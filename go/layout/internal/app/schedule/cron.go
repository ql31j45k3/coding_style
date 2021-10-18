package schedule

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ql31j45k3/coding_style/go/layout/internal/modules/transaction"

	"github.com/ql31j45k3/coding_style/go/layout/internal/modules/order"

	"github.com/ql31j45k3/coding_style/go/layout/internal/utils/tools"

	"github.com/ql31j45k3/coding_style/go/layout/configs"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"go.uber.org/dig"
)

func switchRunning(name string, isStatus func() bool) cron.JobWrapper {
	return func(j cron.Job) cron.Job {
		return cron.FuncJob(func() {
			if isStatus() {
				log.Info(name + " enable")
				j.Run()
			} else {
				log.Info(name + " disable")
			}
		})
	}
}

func runTime(name string) cron.JobWrapper {
	return func(j cron.Job) cron.Job {
		return cron.FuncJob(func() {
			startTime, err := tools.GetNowTime(tools.TimezoneUTC)
			if err != nil {
				log.WithFields(log.Fields{
					"err": err,
				}).Error("runTime - tools.GetNowTime(startTime), name: " + name)
				return
			}

			j.Run()

			endTime, err := tools.GetNowTime(tools.TimezoneUTC)
			if err != nil {
				log.WithFields(log.Fields{
					"err": err,
				}).Error("runTime - tools.GetNowTime(endTime), name: " + name)
				return
			}

			log.WithFields(log.Fields{
				"startTime":  startTime.Format(time.RFC3339),
				"endTime":    endTime.Format(time.RFC3339),
				"subSeconds": endTime.Sub(startTime).Seconds(),
			}).Info("runTime - name: " + name)
		})
	}
}

func notifyJobFinish(ctxStopNotify context.Context, name string, jp *jobPreconditions) cron.JobWrapper {
	return func(j cron.Job) cron.Job {
		done := make(chan jobNotify, 1)
		jp.registerListen(done)

		return cron.FuncJob(func() {
			j.Run()

			nowTimestamp, err := tools.GetNowTimestamp(tools.TimezoneUTC)
			if err != nil {
				log.WithFields(log.Fields{
					"err": err,
				}).Error("notifyJobFinish - tools.GetNowTimestamp, name: " + name)
				return
			}

			select {
			case <-ctxStopNotify.Done():
				return
			case done <- jobNotify{
				name:        name,
				executionAt: nowTimestamp,
			}:
			}
		})
	}
}

// jobPreconditions 未使用鎖的功能，
// 在流程上 run func 是順序執行 不會有同時呼叫 register 相關函數
// run func 執行完，才會呼叫 start 另開 go 做讀資料邏輯，故流程上不會有同時搶佔狀況
type jobPreconditions struct {
	_ struct{}

	listenJob         []<-chan jobNotify
	jobName2NotifyJob map[string][]chan<- jobNotify
}

type jobNotify struct {
	_ struct{}

	name        string
	executionAt int64
}

func newJobPreconditions() *jobPreconditions {
	return &jobPreconditions{
		jobName2NotifyJob: make(map[string][]chan<- jobNotify),
	}
}

func (jp *jobPreconditions) start(ctxStopNotify context.Context) {
	go func(ctxStopNotify context.Context, jp *jobPreconditions) {
		ticker := time.NewTicker(500 * time.Millisecond)
		defer ticker.Stop()

		for {
			<-ticker.C

			for _, listenJob := range jp.getListenJob() {
				select {
				case <-ctxStopNotify.Done():
					return
				case j := <-listenJob:
					for _, notifyJob := range jp.getNotifyJob(j.name) {
						select {
						case <-ctxStopNotify.Done():
							return
						case notifyJob <- j:
						default:
							waitTime := time.NewTimer(500 * time.Millisecond)
							<-waitTime.C

							waitTime.Stop()

							// 注意: 間隔後還是阻塞，直接放棄資料，走不阻塞邏輯並紀錄 log
							select {
							case notifyJob <- j:
							default:
								// 不阻塞
								log.WithFields(log.Fields{
									"err":       errors.New("reTry notifyJob <- j wait 1sec fail, run default data loss"),
									"notifyJob": fmt.Sprintf("%+v", notifyJob),
								}).Error("jobPreconditions.start")
							}
						}
					}
				}
			}
		}
	}(ctxStopNotify, jp)
}

func (jp *jobPreconditions) registerListen(job <-chan jobNotify) {
	jp.listenJob = append(jp.listenJob, job)
}

func (jp *jobPreconditions) getListenJob() []<-chan jobNotify {
	return jp.listenJob
}

func (jp *jobPreconditions) registerNotify(name string, job chan<- jobNotify) {
	v := jp.jobName2NotifyJob[name]
	v = append(v, job)

	jp.jobName2NotifyJob[name] = v
}

func (jp *jobPreconditions) getNotifyJob(name string) []chan<- jobNotify {
	v, ok := jp.jobName2NotifyJob[name]
	if !ok {
		return make([]chan<- jobNotify, 0)
	}

	return v
}

const (
	jobNameOrder = "order"

	jobNameTransaction = "transaction"
)

type jobOrder struct {
	_ struct{}
}

func (j *jobOrder) addJob(ctxStopNotify context.Context, c *cron.Cron, jp *jobPreconditions, container *dig.Container) error {
	cronLog := cron.VerbosePrintfLogger(log.StandardLogger())

	_, err := c.AddJob(configs.Cron.GetOrderSpec(),
		cron.NewChain(switchRunning(jobNameOrder, configs.Cron.GetOrderStatus),
			cron.SkipIfStillRunning(cronLog), runTime(jobNameOrder), notifyJobFinish(ctxStopNotify, jobNameOrder, jp)).
			Then(j.getFunc(ctxStopNotify, container)),
	)
	if err != nil {
		return fmt.Errorf("c.AddJob - jobOrder - %w", err)
	}

	return nil
}

func (j *jobOrder) getFunc(ctxStopNotify context.Context, container *dig.Container) cron.FuncJob {
	return func() {
		err := container.Invoke(func() {
			order.StartOrder(ctxStopNotify)
		})
		if err != nil {
			log.WithFields(log.Fields{
				"err": err,
			}).Error("container.Invoke(order.StartOrder)")
			return
		}
	}
}

type jobTransaction struct {
	_ struct{}
}

func (j *jobTransaction) addJob(ctxStopNotify context.Context, c *cron.Cron, jp *jobPreconditions, container *dig.Container) error {
	cronLog := cron.VerbosePrintfLogger(log.StandardLogger())

	_, err := c.AddJob(configs.Cron.GetTransactionSpec(),
		cron.NewChain(switchRunning(jobNameTransaction, j.isStatus),
			cron.SkipIfStillRunning(cronLog), runTime(jobNameTransaction)).
			Then(j.getFunc(ctxStopNotify, jp, container)))
	if err != nil {
		return fmt.Errorf("c.AddJob - jobTransaction - %w", err)
	}

	return nil
}

func (j *jobTransaction) isStatus() bool {
	var isStatus bool

	// 此排程有前置條件，order 關閉 則一樣無法執行排程
	if configs.Cron.GetOrderStatus() && configs.Cron.GetTransactionStatus() {
		isStatus = true
	} else {
		isStatus = false
	}

	if configs.Cron.GetEnforceTransactionStatus() {
		isStatus = true
	}

	return isStatus
}

func (j *jobTransaction) getFunc(ctxStopNotify context.Context, jp *jobPreconditions, container *dig.Container) cron.FuncJob {
	preconditionsDone := j.notify(ctxStopNotify, jp)

	return func() {
		select {
		case <-ctxStopNotify.Done():
			return
		case <-preconditionsDone:
			err := container.Invoke(func(in containerIn) {
				transaction.StartTransaction()
			})
			if err != nil {
				log.WithFields(log.Fields{
					"err": err,
				}).Error("container.Invoke(transaction.StartTransaction)")
				return
			}
		}
	}
}

func (j *jobTransaction) notify(ctxStopNotify context.Context, jp *jobPreconditions) <-chan struct{} {
	notifyJob := make(chan jobNotify, 1)
	jp.registerNotify(jobNameOrder, notifyJob)

	preconditionsDone := make(chan struct{}, 1)
	go j.checkPreconditions(ctxStopNotify, preconditionsDone, notifyJob)

	return preconditionsDone
}

func (j *jobTransaction) checkPreconditions(ctxStopNotify context.Context, preconditionsDone chan<- struct{}, notifyJob <-chan jobNotify) {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		// 注意: 避免 default 一直執行強佔 CPU
		<-ticker.C

		select {
		case <-ctxStopNotify.Done():
			return
		case v := <-notifyJob:
			executionTime, err := tools.GetTimestampToTime(v.executionAt, tools.TimezoneUTC)
			if err != nil {
				log.WithFields(log.Fields{
					"err": err,
				}).Error("checkPreconditions - tools.GetTimestampToTime")
				return
			}

			executionTime = executionTime.UTC()

			nowTime, err := tools.GetNowTime(tools.TimezoneUTC)
			if err != nil {
				log.WithFields(log.Fields{
					"err": err,
				}).Error("checkPreconditions - tools.GetNowTime")
				return
			}

			// 注意：代表前置條件排程在同時段已執行完畢
			if executionTime.Year() == nowTime.Year() && executionTime.Month() == nowTime.Month() &&
				executionTime.Day() == nowTime.Day() && executionTime.Hour() == nowTime.Hour() {
				select {
				case preconditionsDone <- struct{}{}:
					// preconditionsDone 無法發送成功，只有排程關閉場景造成無人取出資料而阻塞
				default:
					// 不阻塞處理
				}
			} else {
				log.WithFields(log.Fields{
					"executionTime": executionTime.String(),
					"nowTime":       nowTime.String(),
					"message":       "execution conditions are not met",
				}).Info("checkPreconditions - name: " + jobNameTransaction)
			}
		default:
			if configs.Cron.GetEnforceTransactionStatus() {
				select {
				case preconditionsDone <- struct{}{}:
					// preconditionsDone 無法發送成功，只有排程關閉場景造成無人取出資料而阻塞
				default:
					// 不阻塞處理
				}
			}
		}
	}
}

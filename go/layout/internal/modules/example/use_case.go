package example

import (
	"context"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/ql31j45k3/coding_style/go/layout/internal/utils/tools"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

func newUseCaseExample() useCaseExample {
	return &example{}
}

type useCaseExample interface {
	createGoroutine(ctxStopNotify context.Context, taskID string, task *taskMap, goCount int)
	getGoroutineStatus(c *gin.Context, ids []string, task *taskMap)
}

type example struct {
	_ struct{}
}

func (e *example) createGoroutine(ctxStopNotify context.Context, taskID string, task *taskMap, goCount int) {
	taskData := e.producer(ctxStopNotify, taskID, task, goCount)
	e.consumer(taskData)
}

func (e *example) producer(ctxStopNotify context.Context, taskID string, task *taskMap, goCount int) <-chan taskData {
	data := make(chan taskData, 1)

	go func(ctxStopNotify context.Context, data chan<- taskData, taskID string, task *taskMap, goCount int) {
		defer func() {
			if err := recover(); err != nil {
				log.WithFields(log.Fields{
					"err": err,
				}).Panic("producer")

				debug.PrintStack()
			}
		}()

		defer close(data)

		detail := taskMapDetail{
			goCount:     goCount,
			finishCount: 0,
		}

		task.set(taskID, detail)

		for i := 0; i < goCount; i++ {
			select {
			case <-ctxStopNotify.Done():
				v, ok := task.getAndExist(taskID)
				if !ok {
					err := fmt.Errorf("task not find key: %s map", taskID)
					log.WithFields(log.Fields{
						"err": err,
					}).Panic("producer")
					return
				}

				v.goCount = i - 1
				task.set(taskID, v)

				return
			default:
				data <- taskData{
					taskID: taskID,

					task: task,
				}
			}
		}
	}(ctxStopNotify, data, taskID, task, goCount)

	return data
}

func (e *example) consumer(data <-chan taskData) {
	defer func() {
		if err := recover(); err != nil {
			log.WithFields(log.Fields{
				"err": err,
			}).Panic("consumer")

			debug.PrintStack()
		}
	}()

	for v := range data {
		time.Sleep(2 * time.Second)

		detail, ok := v.task.getAndExist(v.taskID)
		if !ok {
			err := fmt.Errorf("task not find key: %s map", v.taskID)
			log.WithFields(log.Fields{
				"err": err,
			}).Panic("consumer")
			return
		}

		detail.finishCount++

		v.task.set(v.taskID, detail)
	}
}

func (e *example) getGoroutineStatus(c *gin.Context, ids []string, task *taskMap) {
	result := make([]responseTaskStatus, 0, len(ids))

	for _, id := range ids {
		v, ok := task.getAndExist(id)
		if !ok {
			err := fmt.Errorf("task not find key: %s map", id)
			tools.NewReturnError(c, http.StatusBadRequest, err)
			return
		}

		temp := responseTaskStatus{
			TaskID: id,

			Detail: responseTaskStatusDetail{
				GoCount:     v.goCount,
				FinishCount: v.finishCount,
			},
		}

		result = append(result, temp)
	}

	c.JSON(http.StatusOK, tools.NewResponseBasicSuccess(result))
}

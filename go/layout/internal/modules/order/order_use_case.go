package order

import (
	"context"
	"fmt"
	"runtime/debug"
	"time"

	log "github.com/sirupsen/logrus"

	"go.mongodb.org/mongo-driver/mongo"
)

func newUseCaseOrder(repositoryOrder repositoryOrder) useCaseOrder {
	return &order{
		repositoryOrder: repositoryOrder,
	}
}

type useCaseOrder interface {
	GetOrderInfo(mongoRS *mongo.Client, account string) (orderData, error)

	DoSyncRecord(ctxStopNotify context.Context)

	DoSyncRecordForSchedule(ctxStopNotify context.Context, goCount int)
}

type order struct {
	_ struct{}

	repositoryOrder
}

func (o *order) GetOrderInfo(mongoRS *mongo.Client, account string) (orderData, error) {
	return o.GetOrderID(mongoRS, account)
}

func (o *order) DoSyncRecord(ctxStopNotify context.Context) {
	for {
		select {
		// 注意: 收到 shutdown 的通知，中止程式
		case <-ctxStopNotify.Done():
			return
		// case: ... 做業務邏輯
		default:
			return
		}
	}
}

func (o *order) DoSyncRecordForSchedule(ctxStopNotify context.Context, goCount int) {
	taskData := o.producer(ctxStopNotify, goCount)
	o.consumer(taskData)
}

func (o *order) producer(ctxStopNotify context.Context, goCount int) <-chan int {
	data := make(chan int, 1)

	go func(ctxStopNotify context.Context, data chan<- int, goCount int) {
		defer func() {
			if err := recover(); err != nil {
				log.WithFields(log.Fields{
					"err": err,
				}).Panic("producer")

				debug.PrintStack()
			}
		}()

		defer close(data)

		for i := 0; i < goCount; i++ {
			select {
			case <-ctxStopNotify.Done():
				return
			default:
				sum := i + i
				data <- sum
			}
		}
	}(ctxStopNotify, data, goCount)

	return data
}

func (o *order) consumer(data <-chan int) {
	defer func() {
		if err := recover(); err != nil {
			log.WithFields(log.Fields{
				"err": err,
			}).Panic("consumer")

			debug.PrintStack()
		}
	}()

	total := 0

	for v := range data {
		time.Sleep(3 * time.Second)

		total += v
	}

	fmt.Println("total: ", total)
}

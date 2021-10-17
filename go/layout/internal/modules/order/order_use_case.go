package order

import (
	"context"

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

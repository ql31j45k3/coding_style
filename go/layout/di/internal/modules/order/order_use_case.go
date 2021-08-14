package order

import "context"

func newUseCaseOrder(repositoryOrder repositoryOrder) useCaseOrder {
	return &order{
		repositoryOrder: repositoryOrder,
	}
}

type useCaseOrder interface {
	GetOrderInfo(account string) string

	DoSyncRecord(ctxStopNotify context.Context)
}

type order struct {
	_ struct{}

	repositoryOrder
}

func (o *order) GetOrderInfo(account string) string {
	return o.GetOrderID(account)
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

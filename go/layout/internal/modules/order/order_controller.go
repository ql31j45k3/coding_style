package order

import (
	"context"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

func RegisterRouterOrder(ctxStopNotify context.Context, condAPI APIOrderCond) {
	order := newUseCaseOrder(newRepositoryOrder())

	registerRouterOrderApp(ctxStopNotify, condAPI, order)
}

func registerRouterOrderApp(ctxStopNotify context.Context, condAPI APIOrderCond, order useCaseOrder) {
	orderControllerFunc := orderControllerApp{
		ctxStopNotify: ctxStopNotify,

		condAPI: condAPI,

		order: order,
	}

	appGroup := condAPI.R.Group("/v1/app")
	appGroup.GET("/order", orderControllerFunc.get)
}

type orderControllerApp struct {
	_ struct{}

	ctxStopNotify context.Context

	condAPI APIOrderCond

	order useCaseOrder
}

func (oca *orderControllerApp) get(c *gin.Context) {
	// 示範有需要被告知服務停止範例
	oca.order.DoSyncRecord(oca.ctxStopNotify)

	// 示範 api get 範例
	condParameter := orderGetCond{}

	condParameter.parse(c)

	member := oca.condAPI.Member.GetMember()
	orderInfo, err := oca.order.GetOrderInfo(oca.condAPI.MongoRS, member.Account)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	log.WithFields(log.Fields{
		"account": member.Account,
	}).Debug("orderControllerApp - get")

	result := responseOrderGet{
		StartTime: condParameter.startTime,
		EndTime:   condParameter.endTime,
		Timezone:  condParameter.timezone,

		OrderInfo: orderInfo.OrderID,
	}
	c.JSON(http.StatusOK, &result)
}

func StartOrder(ctxStopNotify context.Context) {
	// cron job 已處理 recover 功能

	order := newUseCaseOrder(newRepositoryOrder())

	order.DoSyncRecordForSchedule(ctxStopNotify, 10)
}

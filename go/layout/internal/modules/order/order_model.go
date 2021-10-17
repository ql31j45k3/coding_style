package order

import (
	"fmt"
	"strconv"

	"github.com/ql31j45k3/coding_style/go/layout/internal/modules/member"
	"github.com/ql31j45k3/coding_style/go/layout/internal/utils/tools"

	"github.com/shopspring/decimal"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"go.uber.org/dig"
)

const (
	MongoOrder = "order"
)

type APIOrderCond struct {
	dig.In

	R *gin.Engine

	Member member.UseCaseMember

	MongoRS *mongo.Client `name:"mongoRS"`
}

type orderGetCond struct {
	_ struct{}

	startTime string
	endTime   string
	timezone  string
}

func (ogc *orderGetCond) parse(c *gin.Context) {
	ogc.startTime = c.Query("start_time")
	ogc.endTime = c.Query("end_time")
	ogc.timezone = c.Query("timezone")

	log.WithFields(log.Fields{
		"cond": fmt.Sprintf("%+v", ogc),
	}).Debug("orderGetCond - parse")
}

type responseOrderGet struct {
	_ struct{}

	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Timezone  string `json:"timezone"`

	OrderInfo string `json:"order_info"`
}

type orderData struct {
	_ struct{}

	OrderID string `bson:"order_id"`

	Amount float64 `bson:"amount"`
}

func (orderData) GetDatabase() string {
	return tools.MongoDBTestAPI
}

func (orderData) GetCollectionName(_ string) string {
	return MongoOrder
}

func (orderData) Conv(val interface{}) error {
	v, ok := val.(*orderData)
	if !ok {
		return nil
	}

	var err error
	v.Amount, err = strconv.ParseFloat(decimal.NewFromFloat(v.Amount).Truncate(2).String(), 64)
	if err != nil {
		return fmt.Errorf("strconv.ParseFloat - %w", err)
	}

	return nil
}

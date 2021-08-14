package order

import (
	"fmt"

	"github.com/ql31j45k3/coding_style/go/layout/general/internal/modules/member"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type APIOrderCond struct {
	_ struct{}

	R *gin.Engine

	Member member.UseCaseMember
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

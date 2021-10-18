package example

import (
	"context"
	"net/http"
	"time"

	"github.com/ql31j45k3/coding_style/go/layout/internal/utils/tools"

	uuid "github.com/satori/go.uuid"

	"github.com/gin-gonic/gin"
)

func RegisterRouter(ctxStopNotify context.Context, addStopJob func(f func()), condAPI APIExampleCond) {
	example := newUseCaseExample()

	router := exampleRouter{
		ctxStopNotify: ctxStopNotify,

		task: &taskMap{
			taskID2Detail: make(map[string]taskMapDetail),
		},

		example: example,
	}

	addStopJob(router.checkGoroutineStatus)

	routerGroup := condAPI.R.Group("/v1/example")

	// 注意: 示範非同步的 API
	routerGroup.POST("/multi-goroutine", router.createGoroutine)
	routerGroup.GET("/multi-goroutine/status", router.getGoroutineStatus)
}

type exampleRouter struct {
	_ struct{}

	ctxStopNotify context.Context

	task *taskMap

	example useCaseExample
}

func (er *exampleRouter) createGoroutine(c *gin.Context) {
	taskID := uuid.NewV4().String()

	go er.example.createGoroutine(er.ctxStopNotify, taskID, er.task, 10)

	result := responseGoroutine{
		TaskID: taskID,
	}

	c.JSON(http.StatusOK, tools.NewResponseBasicSuccess(result))
}

func (er *exampleRouter) getGoroutineStatus(c *gin.Context) {
	ids := c.QueryArray("ids")

	er.example.getGoroutineStatus(c, ids, er.task)
}

func (er *exampleRouter) checkGoroutineStatus() {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		<-ticker.C

		isAllFinished := false
		er.task.RLock()

		for _, v := range er.task.taskID2Detail {
			if v.goCount == v.finishCount {
				isAllFinished = true
			} else {
				isAllFinished = false
				break
			}
		}

		er.task.RUnlock()

		// 注意: 已經跑完 或 沒有執行過
		if isAllFinished || len(er.task.taskID2Detail) == 0 {
			return
		}
	}
}

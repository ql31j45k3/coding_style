package driver

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ql31j45k3/coding_style/go/layout/configs"

	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/dig"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

func NewGin() *gin.Engine {
	router := gin.Default()

	return router
}

type GinCond struct {
	dig.In

	R *gin.Engine

	MongoRS *mongo.Client `name:"mongoRS"`
}

func StartGin(cancelCtxStopNotify context.CancelFunc, stopFunc func() context.Context, cond GinCond) {
	// 控制調試日誌 log
	gin.SetMode(configs.Gin.GetMode())

	srv := &http.Server{
		Addr:    configs.Host.GetAPIHost(),
		Handler: cond.R,
	}

	go func(srv *http.Server) {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.WithFields(log.Fields{
				"condition": "err != nil and err != http.ErrServerClosed",
				"err":       err,
			}).Error("StartGin - srv.ListenAndServe")
			return
		}
	}(srv)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// 注意: 業務邏輯有做 goroutine 需用 cancel 通知，確保 goroutine 都有正常中止
	cancelCtxStopNotify()

	ctxStop := stopFunc()
	<-ctxStop.Done()

	ctx, cancelCtx := context.WithTimeout(context.Background(), configs.Env.GetShutdownTimeout())
	defer cancelCtx()
	if err := srv.Shutdown(ctx); err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("StartGin - srv.Shutdown")
		return
	}

	ctxMongo, cancelCtxMongo := context.WithTimeout(context.Background(), configs.Mongo.GetTimeout())
	defer cancelCtxMongo()
	if err := Disconnect(ctxMongo, cond.MongoRS); err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("StartGin - Mongo Disconnect")
		// 故意不中斷，後續流程有其他功能需做關閉動作
	}

	log.WithFields(log.Fields{
		"shutdownTimeout": fmt.Sprintf("%ds", int64(configs.Env.GetShutdownTimeout()/time.Second)),
	}).Info("StartGin - Server exiting")
}

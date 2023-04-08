package main

import (
	"context"
	"layout_2/configs"
	deliveryHttp "layout_2/internal/delivery/http"
	"layout_2/internal/libs/container"
	"layout_2/internal/libs/mongodb"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/pyroscope-io/pyroscope/pkg/agent/profiler"

	"fmt"
	"layout_2/internal/libs/logs"
	"layout_2/internal/libs/response"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func main() {
	if err := configs.Start(); err != nil {
		panic(fmt.Sprintf("configs.Start, err: %s", err))
	}

	if err := logs.SetLogEnv(); err != nil {
		panic(fmt.Sprintf("logs.SetLogEnv, err: %s", err))
	}
	configs.SetReloadFunc(logs.ReloadSetLogLevel)

	if configs.App.GetPyroscopeIsRunStart() {
		_, err := profiler.Start(profiler.Config{
			ApplicationName: configs.App.GetServiceName(),
			ServerAddress:   configs.App.GetPyroscopeURL(),
		})

		if err != nil {
			log.WithFields(log.Fields{
				"err": err,
			}).Error("profiler.Start fail")
			return
		}
	}

	log.WithFields(log.Fields{
		"app": fmt.Sprintf("%+v", configs.App),
	}).Debug("check configs app value")

	container.Init()

	if err := container.ProvideInfra(); err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("container.ProvideInfra")
		return
	}

	if err := deliveryHttp.Init(); err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("deliveryHttp.Init")
		return
	}

	// runGin
	if err := container.Get().Invoke(func(r *gin.Engine, mongoClient *mongo.Client) {
		if err := response.RegisterValidator(); err != nil {
			log.WithFields(log.Fields{
				"err": err,
			}).Error("response.RegisterValidator")
			return
		}

		runGin(r, mongoClient)
	}); err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("runGin")
		return
	}
}

func runGin(r *gin.Engine, mongoClient *mongo.Client) {
	srv := newGin(r)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	shutdownTimeout := 10 * time.Second
	ctx, cancelCtx := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancelCtx()

	if err := srv.Shutdown(ctx); err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("srv.Shutdown")
		return
	}

	ctxMongo, cancelCtxMongo := context.WithTimeout(context.Background(), configs.App.GetMongoTimeout())
	defer cancelCtxMongo()
	if err := mongodb.Disconnect(ctxMongo, mongoClient); err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("mongodb.Disconnect")

		// 故意不中斷，後續流程有其他功能需做關閉動作
	}

	timeout := shutdownTimeout / time.Second
	log.WithFields(log.Fields{
		"shutdownTimeout": fmt.Sprintf("%ds", timeout),
	}).Info("Server exiting")
}

func newGin(r *gin.Engine) *http.Server {
	gin.SetMode(configs.App.GetGinMode())

	srv := &http.Server{
		Addr:    configs.App.GetServerPort(),
		Handler: r,
	}

	// 優雅關閉功能，無法攔截 kill -9 信號
	// 當服務做 kill 指令時，會攔截信號，並不接受新的 API 請求，
	// 在執行 shutdownTimeout 秒數，讓正在執行 API 功能，盡量執行結束，
	go func(srv *http.Server) {
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.WithFields(log.Fields{
				"condition": "err != nil and err != http.ErrServerClosed",
				"err":       err,
			}).Error("srv.ListenAndServe")
			return
		}

		log.WithFields(log.Fields{
			"msg": err.Error(),
		}).Info("srv.ListenAndServe")
	}(srv)

	return srv
}

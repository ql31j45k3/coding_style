package main

import (
	"context"
	"layout_2/configs"
	studentRouter "layout_2/internal/example-1/student/delivery/http"
	studentRepo "layout_2/internal/example-1/student/repository"
	studentUseCase "layout_2/internal/example-1/student/usecase"
	deliveryHttp "layout_2/internal/example-2/delivery/http/student"
	"layout_2/internal/example-2/repository/student"
	student2 "layout_2/internal/example-2/usecase/student"

	"fmt"
	"layout_2/internal/libs/logs"
	"layout_2/internal/libs/mysql"
	"layout_2/internal/libs/response"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
	"go.uber.org/dig"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

func main() {
	if err := configs.Start(); err != nil {
		panic(fmt.Sprintf("configs.Start, err: %s", err))
	}

	if err := logs.SetLogEnv(); err != nil {
		panic(fmt.Sprintf("logs.SetLogEnv, err: %s", err))
	}
	configs.SetReloadFunc(logs.ReloadSetLogLevel)

	log.WithFields(log.Fields{
		"app": fmt.Sprintf("%+v", configs.App),
	}).Debug("check configs app value")

	container, err := buildContainer()
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("buildContainer")
		return
	}

	// 調用其他函式，函式參數容器會依照 Provide 提供後自行匹配

	// example-1
	if err := container.Invoke(func(cond studentRouter.StudentHandlerCond) {
		studentRouter.RegisterRouter(cond)
	}); err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("studentRouter")
		return
	}

	// example-2
	if err := container.Invoke(func(cond deliveryHttp.StudentHandlerCond) {
		deliveryHttp.RegisterRouter(cond)
	}); err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("studentRouter")
		return
	}

	// runGin
	if err := container.Invoke(func(r *gin.Engine) {
		if err := response.RegisterValidator(); err != nil {
			log.WithFields(log.Fields{
				"err": err,
			}).Error("response.RegisterValidator")
			return
		}

		runGin(r)
	}); err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("runGin")
		return
	}
}

// buildContainer 建立 DI 容器，提供各個函式的 input 參數
func buildContainer() (*dig.Container, error) {
	container := dig.New()
	provideFunc := containerProvide{}

	if err := container.Provide(provideFunc.gin); err != nil {
		return nil, fmt.Errorf("container.Provide(provideFunc.gin), err: %w", err)
	}

	if err := container.Provide(provideFunc.mysqlMaster, dig.Name("dbM")); err != nil {
		return nil, fmt.Errorf("container.Provide(provideFunc.mysqlMaster), err: %w", err)
	}

	// example-1
	if err := container.Provide(studentRepo.NewStudentRepository); err != nil {
		return nil, fmt.Errorf("container.Provide(studentRepo.NewStudentRepository), err: %w", err)
	}

	if err := container.Provide(studentUseCase.NewStudentUseCase); err != nil {
		return nil, fmt.Errorf("container.Provide(studentUseCase.NewStudentUseCase), err: %w", err)
	}

	// example-2
	if err := container.Provide(student.NewStudentRepository, dig.Name("NewStudentRepository2")); err != nil {
		return nil, fmt.Errorf("container.Provide(repository.NewStudentRepository), err: %w", err)
	}

	if err := container.Provide(student2.NewStudentUseCase, dig.Name("NewStudentUseCase2")); err != nil {
		return nil, fmt.Errorf("container.Provide(usecase.NewStudentUseCase), err: %w", err)
	}

	return container, nil
}

type containerProvide struct {
}

// gin 建立 gin Engine，設定 middleware
func (cp *containerProvide) gin() *gin.Engine {
	return gin.Default()
}

// gorm 建立 gorm.DB 設定，初始化 session 並無實際連線
func (cp *containerProvide) mysqlMaster() (*gorm.DB, error) {
	return mysql.NewMysql(configs.App.GetDBUsername(), configs.App.GetDBPassword(),
		configs.App.GetDBHost(), configs.App.GetDBPort(), configs.App.GetDBName(),
		configs.App.GetGormLogMode())
}

func runGin(r *gin.Engine) {
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

	timeout := shutdownTimeout / time.Second
	log.WithFields(log.Fields{
		"shutdownTimeout": fmt.Sprintf("%ds", timeout),
	}).Info("Server exiting")
}

package api

import (
	"context"
	"fmt"
	"sync"

	"github.com/ql31j45k3/coding_style/go/layout/di/internal/modules/member"

	"github.com/ql31j45k3/coding_style/go/layout/di/internal/modules/order"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"go.uber.org/dig"

	"github.com/ql31j45k3/coding_style/go/layout/di/configs"
	utilsDriver "github.com/ql31j45k3/coding_style/go/layout/di/internal/utils/driver"
)

func Start() {
	// 開始讀取設定檔，順序上必須為容器之前，執行容器內有需要設定檔 struct 取得參數
	if err := configs.Start(); err != nil {
		panic(fmt.Errorf("start - configs.Start - %w", err))
	}

	utilsDriver.SetLogEnv()
	configs.SetReloadFunc(utilsDriver.ReloadSetLogLevel)

	container, err := buildContainer()
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("Start - buildContainer")
		return
	}

	ctxStopNotify, cancelCtxStopNotify := context.WithCancel(context.Background())
	// 注意: cancelCtx 底層保證多個調用，只會執行一次
	defer cancelCtxStopNotify()

	stopJobFunc := stopJob{}

	if err := container.Invoke(func(condAPI order.APIOrderCond) {
		order.RegisterRouterOrder(ctxStopNotify, condAPI)
	}); err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("Start - container.Invoke(order.RegisterRouterOrder)")
		return
	}

	err = container.Invoke(func(in containerIn) {
		utilsDriver.StartGin(cancelCtxStopNotify, stopJobFunc.stop, in.R)
	})
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("Start - utilsDriver.StartGin")
		return
	}
}

// stopJob 為避免其它 package 需 import 此包 package，故用傳遞 func 方式提供功能給其它模組使用，
// 依賴關係都是 start.go 單向 import 其它 package 包功能
type stopJob struct {
	_ struct{}

	sync.Mutex
	stopFunctions []func()
}

func (s *stopJob) stop() context.Context {
	ctx, cancelCtx := context.WithCancel(context.Background())

	go func(s *stopJob, cancelCtx context.CancelFunc) {
		s.Lock()
		defer s.Unlock()

		defer cancelCtx()

		for _, f := range s.stopFunctions {
			f()
		}
	}(s, cancelCtx)

	return ctx
}

func (s *stopJob) add(f func()) {
	s.Lock()
	defer s.Unlock()

	s.stopFunctions = append(s.stopFunctions, f)
}

type containerIn struct {
	dig.In

	R *gin.Engine
}

type containerProvide struct {
	_ struct{}
}

// buildContainer 建立 DI 容器，提供各個函式的 input 參數
func buildContainer() (*dig.Container, error) {
	container := dig.New()
	provideFunc := containerProvide{}

	if err := container.Provide(provideFunc.gin); err != nil {
		return nil, fmt.Errorf("container.Provide(provideFunc.gin) - %w", err)
	}

	if err := member.RegisterContainer(container); err != nil {
		return nil, fmt.Errorf("member.RegisterContainer - %w", err)
	}

	return container, nil
}

// gin 建立 gin Engine，設定 middleware
func (cp *containerProvide) gin() *gin.Engine {
	return utilsDriver.NewGin()
}

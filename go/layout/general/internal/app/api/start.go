package api

import (
	"context"
	"fmt"
	"sync"

	"github.com/ql31j45k3/coding_style/go/layout/general/internal/modules/member"

	"github.com/ql31j45k3/coding_style/go/layout/general/internal/modules/order"

	"github.com/gin-gonic/gin"

	"github.com/ql31j45k3/coding_style/go/layout/general/configs"
	utilsDriver "github.com/ql31j45k3/coding_style/go/layout/general/internal/utils/driver"
)

func Start() {
	// 開始讀取設定檔，順序上必須為容器之前，執行容器內有需要設定檔 struct 取得參數
	if err := configs.Start(); err != nil {
		panic(fmt.Errorf("start - configs.Start - %w", err))
	}

	utilsDriver.SetLogEnv()
	configs.SetReloadFunc(utilsDriver.ReloadSetLogLevel)

	ctxStopNotify, cancelCtxStopNotify := context.WithCancel(context.Background())
	// 注意: cancelCtx 底層保證多個調用，只會執行一次
	defer cancelCtxStopNotify()

	stopJobFunc := stopJob{}
	r := utilsDriver.NewGin()

	registerRouterOrder(ctxStopNotify, r)

	utilsDriver.StartGin(cancelCtxStopNotify, stopJobFunc.stop, r)
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

func registerRouterOrder(ctxStopNotify context.Context, r *gin.Engine) {
	condAPI := order.APIOrderCond{
		R:      r,
		Member: member.NewUseCaseMember(),
	}

	order.RegisterRouterOrder(ctxStopNotify, condAPI)
}

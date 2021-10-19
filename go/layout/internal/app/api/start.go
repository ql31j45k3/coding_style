package api

import (
	"context"
	"fmt"
	"net/http"
	"runtime"
	"sync"

	"github.com/ql31j45k3/coding_style/go/layout/internal/modules/example"

	"github.com/pyroscope-io/pyroscope/pkg/agent/profiler"
	"github.com/ql31j45k3/coding_style/go/layout/configs"
	"github.com/ql31j45k3/coding_style/go/layout/internal/modules/member"
	"github.com/ql31j45k3/coding_style/go/layout/internal/modules/order"
	"github.com/ql31j45k3/coding_style/go/layout/internal/modules/system"
	utilsDriver "github.com/ql31j45k3/coding_style/go/layout/internal/utils/driver"

	transactionDep "github.com/ql31j45k3/coding_style/go/layout/internal/modules/transaction/dependency"

	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"go.uber.org/dig"

	_ "net/http/pprof"
)

func Start() {
	// 開始讀取設定檔，順序上必須為容器之前，執行容器內有需要設定檔 struct 取得參數
	if err := configs.Start(); err != nil {
		panic(fmt.Errorf("start - configs.Start - %w", err))
	}

	// 順序必須在 configs 之後，需取得 設定參數
	if configs.IsPrintVersion() {
		return
	}

	utilsDriver.SetLogEnv()
	configs.SetReloadFunc(utilsDriver.ReloadSetLogLevel)

	go func() {
		if configs.Env.GetPPROFBlockStatus() {
			runtime.SetBlockProfileRate(configs.Env.GetPPROFBlockRate())
		}

		if configs.Env.GetPPROFMutexStatus() {
			runtime.SetMutexProfileFraction(configs.Env.GetPPROFMutexRate())
		}

		if configs.Env.GetPPROFStatus() {
			_ = http.ListenAndServe(configs.Host.GetPPROFAPIHost(), nil)
		}
	}()

	if configs.Env.GetProfilerStatus() {
		_, err := profiler.Start(profiler.Config{
			ApplicationName: configs.Env.GetApplicationName(),
			ServerAddress:   configs.Host.GetProfilerAPIDomain(),
		})

		if err != nil {
			log.WithFields(log.Fields{
				"err": err,
			}).Error("Start - profiler.Start")
			return
		}
	}

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

	if err := container.Invoke(func(condAPI example.APIExampleCond) {
		example.RegisterRouter(ctxStopNotify, stopJobFunc.add, condAPI)
	}); err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("Start - container.Invoke(example.RegisterRouter)")
		return
	}

	if err := container.Invoke(func(condAPI system.APIHealthCond) {
		system.RegisterRouterHealth(condAPI)
	}); err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("Start - container.Invoke(system.RegisterRouterHealth)")
		return
	}

	if err := container.Invoke(func(condAPI order.APIOrderCond) {
		order.RegisterRouterOrder(ctxStopNotify, condAPI)
	}); err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("Start - container.Invoke(order.RegisterRouterOrder)")
		return
	}

	err = container.Invoke(func(cond utilsDriver.GinCond) {
		utilsDriver.StartGin(cancelCtxStopNotify, stopJobFunc.stop, cond)
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

	if err := container.Provide(provideFunc.gormM, dig.Name("dbM")); err != nil {
		return nil, fmt.Errorf("container.Provide(provideFunc.gorm) - %w", err)
	}

	if err := container.Provide(provideFunc.mongoRS, dig.Name("mongoRS")); err != nil {
		return nil, fmt.Errorf("container.Provide(provideFunc.mongoRS) - %w", err)
	}

	if err := member.RegisterContainer(container); err != nil {
		return nil, fmt.Errorf("member.RegisterContainer - %w", err)
	}

	if err := transactionDep.RegisterContainerTransaction(container); err != nil {
		return nil, fmt.Errorf("tranaactionDep.RegisterContainerTransaction - %w", err)
	}

	return container, nil
}

// gin 建立 gin Engine，設定 middleware
func (cp *containerProvide) gin() *gin.Engine {
	return utilsDriver.NewGin()
}

func (cp *containerProvide) gormM() (*gorm.DB, error) {
	return utilsDriver.NewPostgresM(configs.Gorm.GetHost(), configs.Gorm.GetUser(), configs.Gorm.GetPassword(),
		configs.Gorm.GetDBName(), configs.Gorm.GetPort(),
		configs.Gorm.GetMaxIdle(), configs.Gorm.GetMaxOpen(), configs.Gorm.GetMaxLifetime(),
		configs.Gorm.GetLogMode())
}

func (cp *containerProvide) mongoRS() (*mongo.Client, error) {
	ctx, cancelCtx := context.WithTimeout(context.Background(), configs.Mongo.GetTimeout())
	defer cancelCtx()

	return utilsDriver.NewMongoDBConnect(ctx, "",
		utilsDriver.WithMongoHosts(configs.Mongo.GetHosts()),
		utilsDriver.WithMongoAuth(configs.Mongo.GetAuthMechanism(), configs.Mongo.GetUsername(), configs.Mongo.GetPassword()),
		utilsDriver.WithMongoReplicaSet(configs.Mongo.GetReplicaSet()),
		utilsDriver.WithMongoPool(configs.Mongo.GetMinPoolSize(), configs.Mongo.GetMaxPoolSize(), configs.Mongo.GetMaxConnIdleTime()),
		utilsDriver.WithMongoPoolMonitor(),
	)
}

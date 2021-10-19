package schedule

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/ql31j45k3/coding_style/go/layout/internal/modules/member"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/pyroscope-io/pyroscope/pkg/agent/profiler"
	"github.com/ql31j45k3/coding_style/go/layout/internal/utils/tools"

	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"

	utilsDriver "github.com/ql31j45k3/coding_style/go/layout/internal/utils/driver"

	transactionDep "github.com/ql31j45k3/coding_style/go/layout/internal/modules/transaction/dependency"

	"github.com/ql31j45k3/coding_style/go/layout/configs"
	"github.com/ql31j45k3/coding_style/go/layout/internal/modules/index"
	"go.uber.org/dig"
	"gorm.io/gorm"

	_ "net/http/pprof"
)

// Start 控制服務流程、呼叫的依賴性
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

	c := newCron()
	jp := newJobPreconditions()

	ctxStopNotify, cancelCtxStopNotify := context.WithCancel(context.Background())
	// 注意: cancelCtx 底層保證多個調用，只會執行一次
	defer cancelCtxStopNotify()

	// 調用其他函式，函式參數容器會依照 Provide 提供後自行匹配
	if err := run(ctxStopNotify, c, jp, container); err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("Start - startRunJob")
		//nolint:govet
		return
	}

	c.Start()
	jp.start(ctxStopNotify)

	err = container.Invoke(func(in containerIn) {
		shutdown(cancelCtxStopNotify, c, in.MongoRS)
	})
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("Start - shutdown")
		return
	}
}

func newCron() *cron.Cron {
	location, err := time.LoadLocation(tools.TimezoneTaipei)
	if err != nil {
		panic(err)
	}

	cronLog := cron.VerbosePrintfLogger(log.StandardLogger())

	return cron.New(cron.WithLocation(location), cron.WithChain(cron.Recover(cronLog)), cron.WithLogger(cronLog))
}

func shutdown(cancelCtxStopNotify context.CancelFunc, c *cron.Cron, mongoRS *mongo.Client) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// 注意: 業務邏輯有做 goroutine 需用 cancel 通知，確保 goroutine 都有正常中止
	cancelCtxStopNotify()

	ctxStop := c.Stop()
	<-ctxStop.Done()

	ctxMongo, cancelCtxMongo := context.WithTimeout(context.Background(), configs.Mongo.GetTimeout())
	defer cancelCtxMongo()
	if err := utilsDriver.Disconnect(ctxMongo, mongoRS); err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("shutdown - Mongo Disconnect")
		// 故意不中斷，後續流程有其他功能需做關閉動作
	}

	log.WithFields(log.Fields{
		"shutdownTimeout": fmt.Sprintf("%ds", int64(configs.Env.GetShutdownTimeout()/time.Second)),
	}).Info("shutdown - Server exiting")
}

type containerIn struct {
	dig.In

	DBM     *gorm.DB      `name:"dbM"`
	MongoRS *mongo.Client `name:"mongoRS"`
}

// buildContainer 建立 DI 容器，提供各個函式的 input 參數
func buildContainer() (*dig.Container, error) {
	container := dig.New()
	provideFunc := containerProvide{}

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

type containerProvide struct {
	_ struct{}
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

func run(ctxStopNotify context.Context, c *cron.Cron, jp *jobPreconditions, container *dig.Container) error {
	err := container.Invoke(func(in containerIn) {
		index.StartCreateIndex(in.MongoRS)
	})
	if err != nil {
		return fmt.Errorf("container.Invoke(index.StartCreateIndex) - %w", err)
	}

	jOrder := jobOrder{}
	if err := jOrder.addJob(ctxStopNotify, c, jp, container); err != nil {
		return fmt.Errorf("jOrder.addJob - %w", err)
	}

	jTransaction := jobTransaction{}
	if err := jTransaction.addJob(ctxStopNotify, c, jp, container); err != nil {
		return fmt.Errorf("jTransaction.addJob - %w", err)
	}

	return nil
}

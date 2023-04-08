package container

import (
	"context"
	"fmt"
	"layout_2/configs"
	"layout_2/internal/libs/mongodb"
	"layout_2/internal/libs/mysql"
	redisLib "layout_2/internal/libs/redis"
	"layout_2/internal/utils"

	"go.mongodb.org/mongo-driver/mongo"

	"sync"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
	"go.uber.org/dig"
	"gorm.io/gorm"
)

var (
	c *dig.Container

	once sync.Once
)

func Init() {
	once.Do(func() {
		c = dig.New()
	})
}

func Get() *dig.Container {
	return c
}

func ProvideInfra() error {
	provideFunc := containerProvide{}

	if err := c.Provide(provideFunc.gin); err != nil {
		return fmt.Errorf("c.Provide(provideFunc.gin), err: %w", err)
	}

	if err := c.Provide(provideFunc.mysqlMaster, dig.Name("dbM")); err != nil {
		return fmt.Errorf("c.Provide(provideFunc.mysqlMaster), err: %w", err)
	}

	if err := c.Provide(provideFunc.mongo); err != nil {
		return fmt.Errorf("c.Provide(provideFunc.mongo), err: %w", err)
	}

	if err := c.Provide(provideFunc.redis); err != nil {
		return fmt.Errorf("c.Provide(provideFunc.redis), err: %w", err)
	}

	return nil
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
		configs.App.GetGormLogMode(),
		configs.App.GetMaxIdleConns(), configs.App.GetMaxOpenConns(), configs.App.GetConnMaxLifetime())
}

func (cp *containerProvide) mongo() (*mongo.Client, error) {
	ctx, cancelCtx := context.WithTimeout(context.Background(), configs.App.GetMongoTimeout())
	defer cancelCtx()

	return mongodb.NewMongoDBConnect(ctx, "",
		mongodb.WithMongoHosts(configs.App.GetMongoHosts()),
		mongodb.WithMongoAuth(configs.App.GetMongoAuthMechanism(), configs.App.GetMongoUsername(), configs.App.GetMongoPassword()),
		mongodb.WithMongoReplicaSet(configs.App.GetMongoReplicaSet()),
		mongodb.WithMongoPool(configs.App.GetMongoMinPoolSize(), configs.App.GetMongoMaxPoolSize(), configs.App.GetMongoMaxConnIdleTime()),
		mongodb.WithMongoPoolMonitor(configs.App.GetMongoDebug()),
	)
}

func (cp *containerProvide) redis() (*redis.ClusterClient, error) {
	ctx, cancelCtx := context.WithTimeout(context.Background(), utils.Time30S)
	defer cancelCtx()

	return redisLib.NewRedisConnect(ctx, configs.App.GetRedisHosts(),
		configs.App.GetRedisPassword(), configs.App.GetRedisPoolSize())
}

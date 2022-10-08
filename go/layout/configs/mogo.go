package configs

import (
	"time"

	"github.com/spf13/viper"
)

func newConfigMongo() *configMongo {
	viper.SetDefault("database.mongo.timeout", time.Duration(300))

	viper.SetDefault("database.mongo.pool.minSize", 10)
	viper.SetDefault("database.mongo.pool.maxSize", 100)
	viper.SetDefault("database.mongo.pool.maxConnIdleTime", time.Duration(300))

	return &configMongo{
		timeout: viper.GetDuration("database.mongo.timeout") * time.Second,

		authMechanism: viper.GetString("database.mongo.authMechanism"),
		replicaSet:    viper.GetString("database.mongo.replicaSet"),

		hosts:    viper.GetStringSlice("database.mongo.hosts"),
		host:     viper.GetStringSlice("database.mongo.host"),
		port:     viper.GetStringSlice("database.mongo.port"),
		username: viper.GetString("database.mongo.username"),
		password: viper.GetString("database.mongo.password"),

		minPoolSize:     viper.GetUint64("database.mongo.pool.minSize"),
		maxPoolSize:     viper.GetUint64("database.mongo.pool.maxSize"),
		maxConnIdleTime: viper.GetDuration("database.mongo.pool.maxConnIdleTime") * time.Second,
	}
}

type configMongo struct {
	timeout time.Duration

	authMechanism string
	replicaSet    string

	hosts    []string
	host     []string
	port     []string
	username string
	password string

	minPoolSize     uint64
	maxPoolSize     uint64
	maxConnIdleTime time.Duration
}

func (c *configMongo) GetTimeout() time.Duration {
	return c.timeout
}

func (c *configMongo) GetAuthMechanism() string {
	if c.authMechanism == "Direct" || c.authMechanism == "PLAIN" || c.authMechanism == "SCRAM" {
		return c.authMechanism
	}

	c.authMechanism = "Direct"
	return c.authMechanism
}

func (c *configMongo) GetReplicaSet() string {
	return c.replicaSet
}

func (c *configMongo) GetHosts() []string {
	return c.hosts
}

func (c *configMongo) GetHost() []string {
	return c.host
}

func (c *configMongo) GetPort() []string {
	return c.port
}

func (c *configMongo) GetUsername() string {
	return c.username
}

func (c *configMongo) GetPassword() string {
	return c.password
}

func (c *configMongo) GetMinPoolSize() uint64 {
	return c.minPoolSize
}

func (c *configMongo) GetMaxPoolSize() uint64 {
	return c.maxPoolSize
}

func (c *configMongo) GetMaxConnIdleTime() time.Duration {
	return c.maxConnIdleTime
}

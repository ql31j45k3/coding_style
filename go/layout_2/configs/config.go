package configs

import (
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"gorm.io/gorm/logger"

	"github.com/spf13/viper"
)

const (
	envServiceName = "env.serviceName"

	envLogLevel = "env.log.level"
	envLogPath  = "env.log.path"

	envDebug      = "env.debug"
	envServerPort = "env.server.port"

	pyroscopeIsRunStart = "pyroscope.isRunStart"
	pyroscopeURL        = "pyroscope.url"

	dbGormLogMode = "db.gorm.logMode"

	dbMasterHost     = "db.master.host"
	dbMasterPort     = "db.master.port"
	dbMasterUsername = "db.master.username"
	dbMasterPassword = "db.master.password"
	dbMasterName     = "db.master.name"

	mongoTimeout = "mongo.timeout"

	mongoAuthMechanism = "mongo.authMechanism"
	mongoReplicaSet    = "mongo.replicaSet"

	mongoHosts    = "mongo.hosts"
	mongoUsername = "mongo.username"
	mongoPassword = "mongo.password"

	mongoPoolDebug           = "mongo.pool.debug"
	mongoPoolMinSize         = "mongo.pool.minSize"
	mongoPoolMaxSize         = "mongo.pool.maxSize"
	mongoPoolMaxConnIdleTime = "mongo.pool.maxConnIdleTime"

	redisHosts    = "redis.hosts"
	redisPassword = "redis.password"
	redisPoolSize = "redis.poolSize"
)

func newConfigApp() *configApp {
	viper.SetDefault(envLogLevel, "warn")
	viper.SetDefault(envLogPath, "/var/log/layout_demo_log")

	viper.SetDefault(envServerPort, "8080")

	viper.SetDefault(dbGormLogMode, "warn")

	viper.SetDefault(mongoTimeout, time.Duration(300))

	viper.SetDefault(mongoPoolMinSize, 10)
	viper.SetDefault(mongoPoolMaxSize, 100)
	viper.SetDefault(mongoPoolMaxConnIdleTime, time.Duration(300))

	viper.SetDefault(redisPoolSize, 100)

	return &configApp{
		serviceName: viper.GetString(envServiceName),

		logLevel: viper.GetString(envLogLevel),
		logPath:  viper.GetString(envLogPath),

		debug:      viper.GetBool(envDebug),
		serverPort: ":" + viper.GetString(envServerPort),

		pyroscopeIsRunStart: viper.GetBool(pyroscopeIsRunStart),
		pyroscopeURL:        viper.GetString(pyroscopeURL),

		gormLogMode: viper.GetString(dbGormLogMode),

		dbHost:     viper.GetString(dbMasterHost),
		dbPort:     viper.GetString(dbMasterPort),
		dbUsername: viper.GetString(dbMasterUsername),
		dbPassword: viper.GetString(dbMasterPassword),
		dbName:     viper.GetString(dbMasterName),

		mongoTimeout: viper.GetDuration(mongoTimeout) * time.Second,

		mongoAuthMechanism: viper.GetString(mongoAuthMechanism),
		mongoReplicaSet:    viper.GetString(mongoReplicaSet),

		mongoHosts:    viper.GetStringSlice(mongoHosts),
		mongoUsername: viper.GetString(mongoUsername),
		mongoPassword: viper.GetString(mongoPassword),

		mongoDebug:           viper.GetBool(mongoPoolDebug),
		mongoMinPoolSize:     viper.GetUint64(mongoPoolMinSize),
		mongoMaxPoolSize:     viper.GetUint64(mongoPoolMaxSize),
		mongoMaxConnIdleTime: viper.GetDuration(mongoPoolMaxConnIdleTime) * time.Second,

		redisHosts:    viper.GetStringSlice(redisHosts),
		redisPassword: viper.GetString(redisPassword),
		redisPoolSize: viper.GetInt(redisPoolSize),
	}
}

type configApp struct {
	sync.RWMutex

	serviceName string

	logLevel string
	logPath  string

	debug      bool
	serverPort string

	pyroscopeIsRunStart bool
	pyroscopeURL        string

	gormLogMode string

	dbHost     string
	dbPort     string
	dbUsername string
	dbPassword string
	dbName     string

	mongoTimeout time.Duration

	mongoAuthMechanism string
	mongoReplicaSet    string

	mongoHosts    []string
	mongoUsername string
	mongoPassword string

	mongoDebug           bool
	mongoMinPoolSize     uint64
	mongoMaxPoolSize     uint64
	mongoMaxConnIdleTime time.Duration

	redisHosts    []string
	redisPassword string
	redisPoolSize int
}

func (c *configApp) reload() {
	c.Lock()
	defer c.Unlock()

	c.logLevel = viper.GetString(envLogLevel)
}

func (c *configApp) GetServiceName() string {
	return c.serviceName
}

func (c *configApp) GetLogLevel() string {
	c.RLock()
	defer c.RUnlock()

	return c.logLevel
}

func (c *configApp) GetLogPath() string {
	return c.logPath
}

func (c *configApp) GetDebug() bool {
	return c.debug
}

func (c *configApp) GetGinMode() string {
	if c.GetDebug() {
		return gin.DebugMode
	}

	return gin.ReleaseMode
}

func (c *configApp) GetServerPort() string {
	return c.serverPort
}

func (c *configApp) GetPyroscopeIsRunStart() bool {
	return c.pyroscopeIsRunStart
}

func (c *configApp) GetPyroscopeURL() string {
	return c.pyroscopeURL
}

func (c *configApp) GetGormLogMode() logger.LogLevel {
	logMode := strings.ToLower(c.gormLogMode)

	switch logMode {
	case "silent":
		return logger.Silent
	case "error":
		return logger.Error
	case "warn":
		return logger.Warn
	case "info":
		return logger.Info
	}

	return logger.Warn
}

func (c *configApp) GetDBHost() string {
	return c.dbHost
}

func (c *configApp) GetDBPort() string {
	return c.dbPort
}

func (c *configApp) GetDBUsername() string {
	return c.dbUsername
}

func (c *configApp) GetDBPassword() string {
	return c.dbPassword
}

func (c *configApp) GetDBName() string {
	return c.dbName
}

func (c *configApp) GetMongoTimeout() time.Duration {
	return c.mongoTimeout
}

func (c *configApp) GetMongoAuthMechanism() string {
	if c.mongoAuthMechanism == "Direct" || c.mongoAuthMechanism == "PLAIN" || c.mongoAuthMechanism == "SCRAM" {
		return c.mongoAuthMechanism
	}

	c.mongoAuthMechanism = "Direct"
	return c.mongoAuthMechanism
}

func (c *configApp) GetMongoReplicaSet() string {
	return c.mongoReplicaSet
}

func (c *configApp) GetMongoHosts() []string {
	return c.mongoHosts
}

func (c *configApp) GetMongoUsername() string {
	return c.mongoUsername
}

func (c *configApp) GetMongoPassword() string {
	return c.mongoPassword
}

func (c *configApp) GetMongoDebug() bool {
	return c.mongoDebug
}

func (c *configApp) GetMongoMinPoolSize() uint64 {
	return c.mongoMinPoolSize
}

func (c *configApp) GetMongoMaxPoolSize() uint64 {
	return c.mongoMaxPoolSize
}

func (c *configApp) GetMongoMaxConnIdleTime() time.Duration {
	return c.mongoMaxConnIdleTime
}

func (c *configApp) GetRedisHosts() []string {
	return c.redisHosts
}

func (c *configApp) GetRedisPassword() string {
	return c.redisPassword
}

func (c *configApp) GetRedisPoolSize() int {
	return c.redisPoolSize
}

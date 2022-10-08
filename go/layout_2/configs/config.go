package configs

import (
	"strings"
	"sync"

	"github.com/gin-gonic/gin"

	"gorm.io/gorm/logger"

	"github.com/spf13/viper"
)

const (
	envLogLevel = "env.log.level"
	envLogPath  = "env.log.path"

	envDebug      = "env.debug"
	envServerPort = "env.server.port"

	dbGormLogMode = "db.gorm.logMode"

	dbMasterHost     = "db.master.host"
	dbMasterPort     = "db.master.port"
	dbMasterUsername = "db.master.username"
	dbMasterPassword = "db.master.password"
	dbMasterName     = "db.master.name"
)

func newConfigApp() *configApp {

	viper.SetDefault(envLogLevel, "warn")
	viper.SetDefault(envLogPath, "/var/log/layout_demo_log")

	viper.SetDefault(envServerPort, "8080")

	viper.SetDefault(dbGormLogMode, "warn")

	return &configApp{
		logLevel: viper.GetString(envLogLevel),
		logPath:  viper.GetString(envLogPath),

		debug:      viper.GetBool(envDebug),
		serverPort: ":" + viper.GetString(envServerPort),

		gormLogMode: viper.GetString(dbGormLogMode),

		dbHost:     viper.GetString(dbMasterHost),
		dbPort:     viper.GetString(dbMasterPort),
		dbUsername: viper.GetString(dbMasterUsername),
		dbPassword: viper.GetString(dbMasterPassword),
		dbName:     viper.GetString(dbMasterName),
	}
}

type configApp struct {
	sync.RWMutex

	logLevel string
	logPath  string

	debug      bool
	serverPort string

	gormLogMode string

	dbHost     string
	dbPort     string
	dbUsername string
	dbPassword string
	dbName     string
}

func (c *configApp) reload() {
	c.Lock()
	defer c.Unlock()

	c.logLevel = viper.GetString(envLogLevel)
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

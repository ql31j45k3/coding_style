package configs

import (
	"strings"
	"time"

	"github.com/spf13/viper"
	"gorm.io/gorm/logger"
)

func newConfigGorm() *configGorm {
	viper.SetDefault("gorm.logmode", "silent")

	viper.SetDefault("database.postgres.master.timeout", time.Duration(300))

	viper.SetDefault("database.postgres.master.conn.maxIdle", 10)
	viper.SetDefault("database.postgres.master.conn.maxOpen", 100)
	viper.SetDefault("database.postgres.master.conn.maxLifetime", time.Duration(600))

	config := &configGorm{
		host:     viper.GetString("database.postgres.master.host"),
		port:     viper.GetString("database.postgres.master.port"),
		user:     viper.GetString("database.postgres.master.username"),
		password: viper.GetString("database.postgres.master.password"),

		dbName: viper.GetString("database.postgres.master.dbName"),

		maxIdle:     viper.GetInt("database.postgres.master.conn.maxIdle"),
		maxOpen:     viper.GetInt("database.postgres.master.conn.maxOpen"),
		maxLifetime: viper.GetDuration("database.postgres.master.conn.maxLifetime") * time.Second,

		mode:    viper.GetString("gorm.logmode"),
		timeout: viper.GetDuration("database.postgres.master.timeout") * time.Second,
	}

	return config
}

type configGorm struct {
	_ struct{}

	host     string
	port     string
	user     string
	password string

	dbName string

	maxIdle     int
	maxOpen     int
	maxLifetime time.Duration

	mode    string
	timeout time.Duration
}

func (c *configGorm) GetHost() string {
	return c.host
}

func (c *configGorm) GetPort() string {
	return c.port
}

func (c *configGorm) GetUser() string {
	return c.user
}

func (c *configGorm) GetPassword() string {
	return c.password
}

func (c *configGorm) GetDBName() string {
	return c.dbName
}

func (c *configGorm) GetMaxIdle() int {
	return c.maxIdle
}

func (c *configGorm) GetMaxOpen() int {
	return c.maxOpen
}

func (c *configGorm) GetMaxLifetime() time.Duration {
	return c.maxLifetime
}

func (c *configGorm) GetLogMode() logger.LogLevel {
	if strings.ToLower(c.mode) == "silent" {
		return logger.Silent
	}
	if strings.ToLower(c.mode) == "error" {
		return logger.Error
	}
	if strings.ToLower(c.mode) == "warn" {
		return logger.Warn
	}
	if strings.ToLower(c.mode) == "info" {
		return logger.Info
	}

	return logger.Silent
}

func (c *configGorm) GetTimeout() time.Duration {
	return c.timeout
}

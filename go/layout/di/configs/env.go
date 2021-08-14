package configs

import (
	"sync"
	"time"

	"github.com/spf13/viper"
)

func newConfigEnv() *configEnv {
	viper.SetDefault("api.timeout", time.Duration(30))

	viper.SetDefault("system.log.level", "warn")
	viper.SetDefault("system.log.path", "/var/log/ostrich/def_log")

	viper.SetDefault("system.shutdown.timeout", time.Duration(10))

	return &configEnv{
		apiTimeout: viper.GetDuration("api.timeout") * time.Second,

		logLevel: viper.GetString("system.log.level"),
		logPath:  viper.GetString("system.log.path"),

		shutdownTimeout: viper.GetDuration("system.shutdown.timeout") * time.Second,
	}
}

type configEnv struct {
	sync.RWMutex

	apiTimeout time.Duration

	logLevel string
	logPath  string

	shutdownTimeout time.Duration
}

func (c *configEnv) reload() {
	c.Lock()
	defer c.Unlock()

	c.logLevel = viper.GetString("system.log.level")
}

func (c *configEnv) GetAPITimeout() time.Duration {
	return c.apiTimeout
}

func (c *configEnv) GetLogLevel() string {
	c.RLock()
	defer c.RUnlock()

	return c.logLevel
}

func (c *configEnv) GetLogPath() string {
	return c.logPath
}

func (c *configEnv) GetShutdownTimeout() time.Duration {
	return c.shutdownTimeout
}

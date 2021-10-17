package configs

import (
	"sync"
	"time"

	"github.com/spf13/viper"
)

func newConfigEnv() *configEnv {
	viper.SetDefault("api.timeout", time.Duration(30))

	viper.SetDefault("system.log.level", "warn")
	viper.SetDefault("system.log.path", "/var/log/layout/def_log")

	viper.SetDefault("system.shutdown.timeout", time.Duration(10))

	viper.SetDefault("system.pprof.status", false)

	viper.SetDefault("system.pprof.block.status", false)
	viper.SetDefault("system.pprof.block.rate", 1000000000)

	viper.SetDefault("system.pprof.mutex.status", false)
	viper.SetDefault("system.pprof.mutex.rate", 1000000000)

	return &configEnv{
		apiTimeout: viper.GetDuration("api.timeout") * time.Second,

		logLevel: viper.GetString("system.log.level"),
		logPath:  viper.GetString("system.log.path"),

		shutdownTimeout: viper.GetDuration("system.shutdown.timeout") * time.Second,

		pprofStatus: viper.GetBool("system.pprof.status"),

		pprofBlockStatus: viper.GetBool("system.pprof.block.status"),
		pprofBlockRate:   viper.GetInt("system.pprof.block.rate"),

		pprofMutexStatus: viper.GetBool("system.pprof.mutex.status"),
		pprofMutexRate:   viper.GetInt("system.pprof.mutex.rate"),
	}
}

type configEnv struct {
	sync.RWMutex

	apiTimeout time.Duration

	logLevel string
	logPath  string

	shutdownTimeout time.Duration

	pprofStatus bool

	pprofBlockStatus bool
	pprofBlockRate   int

	pprofMutexStatus bool
	pprofMutexRate   int
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

func (c *configEnv) GetPPROFStatus() bool {
	return c.pprofStatus
}

func (c *configEnv) GetPPROFBlockStatus() bool {
	return c.pprofBlockStatus
}

func (c *configEnv) GetPPROFBlockRate() int {
	return c.pprofBlockRate
}

func (c *configEnv) GetPPROFMutexStatus() bool {
	return c.pprofMutexStatus
}

func (c *configEnv) GetPPROFMutexRate() int {
	return c.pprofMutexRate
}

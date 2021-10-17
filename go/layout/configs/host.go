package configs

import "github.com/spf13/viper"

func newConfigHost() *configHost {
	config := &configHost{
		apiHost:      ":" + viper.GetString("api.port"),
		pprofAPIHost: ":" + viper.GetString("api.pprof.port"),

		profilerAPIDomain: viper.GetString("api.profiler.domain"),
	}

	return config
}

type configHost struct {
	apiHost string

	pprofAPIHost string

	profilerAPIDomain string
}

func (c *configHost) GetAPIHost() string {
	return c.apiHost
}

func (c *configHost) GetPPROFAPIHost() string {
	return c.pprofAPIHost
}

func (c *configHost) GetProfilerAPIDomain() string {
	return c.profilerAPIDomain
}

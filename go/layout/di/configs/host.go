package configs

import "github.com/spf13/viper"

func newConfigHost() *configHost {
	config := &configHost{
		apiHost: ":" + viper.GetString("api.port"),
		pprofAPIHost:  ":" + viper.GetString("api.pprof.port"),
	}

	return config
}

type configHost struct {
	apiHost string

	pprofAPIHost  string
}

func (c *configHost) GetAPIHost() string {
	return c.apiHost
}

func (c *configHost) GetPPROFAPIHost() string {
	return c.pprofAPIHost
}

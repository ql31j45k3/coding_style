package configs

import "github.com/spf13/viper"

func newConfigHost() *configHost {
	config := &configHost{
		apiHost: ":" + viper.GetString("api.port"),
	}

	return config
}

type configHost struct {
	apiHost string
}

func (c *configHost) GetAPIHost() string {
	return c.apiHost
}

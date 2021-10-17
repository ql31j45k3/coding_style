package configs

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func newConfigGin() *configGin {
	viper.SetDefault("gin.mode", gin.DebugMode)

	config := &configGin{
		mode: viper.GetString("gin.mode"),
	}

	return config
}

type configGin struct {
	mode string
}

func (c *configGin) GetMode() string {
	return strings.ToLower(c.mode)
}

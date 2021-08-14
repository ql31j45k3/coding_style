package configs

import (
	"errors"
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	Host *configHost
	Gin  *configGin
	Env  *configEnv

	reloadFunc []func()
)

func SetReloadFunc(f func()) {
	reloadFunc = append(reloadFunc, f)
}

// Start 開始 Config 設定參數與讀取檔案並轉成 struct
func Start() error {
	// 設定自定義 flag to viper
	if err := parseFlag(); err != nil {
		return fmt.Errorf("parseFlag - %w", err)
	}

	viper.AddConfigPath(viper.GetString("configFile"))
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("viper.ReadInConfig - %w", err)
	}

	Host = newConfigHost()
	Gin = newConfigGin()
	Env = newConfigEnv()

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		Env.reload()

		for _, f := range reloadFunc {
			f()
		}
	})

	if err := checkBasicValue(); err != nil {
		return fmt.Errorf("checkBasicValue - %w", err)
	}

	return nil
}

func parseFlag() error {
	pflag.String("configFile", "", "configFile path")

	pflag.Parse()

	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		return fmt.Errorf("viper.BindPFlags - %w", err)
	}

	return nil
}

func checkBasicValue() error {
	if Env.shutdownTimeout < 1 {
		return errors.New("Env.shutdownTimeout < 1, need >= 1")
	}

	if Env.apiTimeout <= 0 {
		return errors.New("Env.apiTimeout <= 0, need > 0")
	}

	return nil
}

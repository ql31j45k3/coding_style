package configs

import (
	"sync"

	"github.com/spf13/viper"
)

func newConfigCron() *configCron {
	viper.SetDefault("cron.enforce.transaction.status", false)

	viper.SetDefault("cron.order.status", false)
	viper.SetDefault("cron.transaction.status", false)

	return &configCron{
		enforceTransactionStatus: viper.GetBool("cron.enforce.transaction.status"),

		orderSpec:   viper.GetString("cron.order.spec"),
		orderStatus: viper.GetBool("cron.order.status"),

		transactionSpec:   viper.GetString("cron.transaction.spec"),
		transactionStatus: viper.GetBool("cron.transaction.status"),
	}
}

type configCron struct {
	sync.RWMutex

	enforceTransactionStatus bool

	orderSpec   string
	orderStatus bool

	transactionSpec   string
	transactionStatus bool
}

func (c *configCron) reload() {
	c.Lock()
	defer c.Unlock()

	c.enforceTransactionStatus = viper.GetBool("cron.enforce.transaction.status")

	c.orderStatus = viper.GetBool("cron.order.status")
	c.transactionStatus = viper.GetBool("cron.transaction.status")
}

func (c *configCron) GetEnforceTransactionStatus() bool {
	c.RLock()
	defer c.RUnlock()

	return c.enforceTransactionStatus
}

func (c *configCron) GetOrderSpec() string {
	return c.orderSpec
}

func (c *configCron) GetOrderStatus() bool {
	c.RLock()
	defer c.RUnlock()

	return c.orderStatus
}

func (c *configCron) GetTransactionSpec() string {
	return c.transactionSpec
}

func (c *configCron) GetTransactionStatus() bool {
	c.RLock()
	defer c.RUnlock()

	return c.transactionStatus
}

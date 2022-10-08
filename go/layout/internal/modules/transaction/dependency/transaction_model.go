package dependency

import "gorm.io/gorm"

const (
	pgTransactions = "transactions"
)

type TransactionLog struct {
	TransAt int64  `bson:"trans_at"`
	Account string `bson:"account"`
	OrderID string `bson:"order_id"`
}

func (TransactionLog) getTableName() func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		return tx.Table(pgTransactions)
	}
}

package dependency

import (
	"gorm.io/gorm"
)

func newUseCaseTransaction(repositoryTransaction repositoryTransaction) UseCaseTransaction {
	return &transaction{
		repositoryTransaction: repositoryTransaction,
	}
}

type UseCaseTransaction interface {
	InsertTransactions(dbM *gorm.DB, transactionLog TransactionLog) error
}

type transaction struct {
	repositoryTransaction
}

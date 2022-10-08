package dependency

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"gorm.io/gorm"
)

func newRepositoryTransaction() repositoryTransaction {
	return &transactionDB{}
}

type repositoryTransaction interface {
	InsertTransactions(dbM *gorm.DB, transactionLog TransactionLog) error
}

type transactionDB struct {
	transactionPostgres
}

type transactionPostgres struct {
}

func (tp *transactionPostgres) InsertTransactions(dbM *gorm.DB, transactionLog TransactionLog) error {
	result := dbM.Scopes(transactionLog.getTableName()).Create(&transactionLog)
	if result.Error != nil {
		log.WithFields(log.Fields{
			"transactionLog": fmt.Sprintf("%+v", transactionLog),
			"result.Error":   result.Error,
		}).Error("InsertTransactions - dbM.Create")

		return fmt.Errorf("dbM.Create(&transaction) - %w", result.Error)
	}

	return nil
}

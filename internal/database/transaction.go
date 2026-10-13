package database

import (
	"context"
	"gorm.io/gorm"
)

type transactionKey struct{}

type TransactionFunc func(ctx context.Context) error

type TxManager struct {
	db *gorm.DB
}

func NewTxManager(db *gorm.DB) *TxManager {
	return &TxManager{
		db: db,
	}
}

func (m *TxManager) WithTransaction(ctx context.Context, fn TransactionFunc) error {
	return m.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txCtx := context.WithValue(ctx, transactionKey{}, tx)

		return fn(txCtx)
	})
}

func getTx(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value(transactionKey{}).(*gorm.DB)
	if !ok {
		return nil
	}

	return tx
}

func DB(ctx context.Context, db *gorm.DB) *gorm.DB {
	if tx := getTx(ctx); tx != nil {
		return tx.WithContext(ctx)
	}

	return db.WithContext(ctx)
}

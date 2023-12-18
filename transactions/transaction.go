package transactions

import (
	"context"
	"database/sql"
)

type Transaction struct {
	db *sql.DB
}

func NewTransaction(db *sql.DB) *Transaction {
	return &Transaction{
		db: db,
	}
}

func (r *Transaction) ExecTx(ctx context.Context, h func(tx *sql.Tx) error) error {

	tx, err := r.db.BeginTx(
		ctx,
		&sql.TxOptions{Isolation: sql.LevelSerializable},
	)
	if err != nil {
		return err
	}

	err = h(tx)
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return err
		}
		return err
	}

	return tx.Commit()

}

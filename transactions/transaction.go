package transactions

import (
	"context"
	"database/sql"

	"github.com/felipeflores/utils/persistence"
)

type Transaction struct {
	db *persistence.Service
}

func NewTransaction(db *persistence.Service) *Transaction {
	return &Transaction{
		db: db,
	}
}

func (r *Transaction) ExecTx(ctx context.Context, h func(tx *sql.Tx) error) error {

	tx, err := r.db.DB.BeginTx(
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

package transactions

import (
	"context"
	"database/sql"
)

type ITransaction interface {
	ExecTx(ctx context.Context, h func(tx *sql.Tx) error) error
}

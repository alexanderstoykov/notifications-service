package storage

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/alexanderstoykov/notifications-service/config"
)

const retryAttempts = 10

type Querier interface {
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	Rebind(query string) string
}

type Connection struct {
	db *sqlx.DB
}

type contextKey string

func (c contextKey) String() string {
	return fmt.Sprintf("context key: %s", string(c))
}

var txKey = contextKey("txKey")

type TxFunc func(ctx context.Context) error

func NewConnection(cfg config.DatabaseConfig) (*Connection, error) {
	db, err := sql.Open("postgres", cfg.DSN())
	if err != nil {
		return nil, err
	}

	sqlxDB := sqlx.NewDb(db, "postgres")

	tries := retryAttempts
	for tries > 0 {
		err := sqlxDB.Ping()
		if err == nil {
			break
		}

		time.Sleep(time.Second * 1)
		tries--

		if tries == 0 {
			return nil, fmt.Errorf("database did not become available within %d connection attempts", retryAttempts)
		}
	}

	return &Connection{db: sqlxDB}, nil
}

// DB returns db instance, if there is a transaction going it returns the transaction.
func (c *Connection) DB(ctx context.Context) Querier {
	if dtp, ok := GetTransactionFromContext(ctx); ok {
		return dtp.tx
	}

	return c.db
}

func GetTransactionFromContext(ctx context.Context) (TxPair, bool) {
	dtp, ok := ctx.Value(txKey).(TxPair)

	return dtp, ok
}

type TxPair struct {
	tx *sqlx.Tx
	db *sqlx.DB
}

func (c *Connection) Tx(ctx context.Context, txFunc TxFunc) error {
	return c.execTransaction(ctx, txFunc)
}

func (c *Connection) execTransaction(ctx context.Context, txFunc TxFunc) error {
	if _, ok := GetTransactionFromContext(ctx); ok {
		return txFunc(ctx)
	}

	errRun := c.runInTx(ctx, nil, func(ctx context.Context, tx *sqlx.Tx) error {
		return txFunc(setTransactionToContext(ctx, c.db, tx))
	})

	return errors.WithStack(errRun)
}
func setTransactionToContext(ctx context.Context, db *sqlx.DB, tx *sqlx.Tx) context.Context {
	return context.WithValue(ctx, txKey, TxPair{
		db: db,
		tx: tx,
	})
}

func (c *Connection) runInTx(
	ctx context.Context,
	opts *sql.TxOptions,
	fn func(ctx context.Context, tx *sqlx.Tx) error,
) error {
	tx, err := c.db.BeginTxx(ctx, opts)
	if err != nil {
		return err
	}

	var done bool

	defer func() {
		if !done {
			if err = tx.Rollback(); err != nil {
				log.Printf("error in tx rollback: %s", err)
			}
		}
	}()

	if err := fn(ctx, tx); err != nil {
		return err
	}

	done = true

	return tx.Commit()
}

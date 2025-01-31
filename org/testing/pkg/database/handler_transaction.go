package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/doug-martin/goqu/v9"
)

// SQLTxOpt defines function format to customize transaction creation
type SQLTxOpt func(sqlOpt *sql.TxOptions)

// GetDB returns underlying DB
func (dbh *HandlerDB) GetDB() *goqu.Database {
	return dbh.db
}

// GetTx returns underlying tx if any, else errors
func (dbh *HandlerDB) GetTx() (*goqu.TxDatabase, error) {
	if dbh.tx != nil {
		return dbh.tx, nil
	}
	return nil, ErrNoActiveTransaction
}

// WithReadonly option utility function used when starting a new transaction
func WithReadonly() SQLTxOpt {
	return func(sqlOpt *sql.TxOptions) {
		sqlOpt.ReadOnly = true
	}
}

// WithIsolationReadCommitted option utility function used when starting a new transaction
func WithIsolationReadCommitted() SQLTxOpt {
	return func(sqlOpt *sql.TxOptions) {
		sqlOpt.Isolation = sql.LevelReadCommitted
	}
}

// WithIsolationRepeatableRead option utility function used when starting a new transaction
func WithIsolationRepeatableRead() SQLTxOpt {
	return func(sqlOpt *sql.TxOptions) {
		sqlOpt.Isolation = sql.LevelRepeatableRead
	}
}

// WithIsolationSerializable option utility function used when starting a new transaction
func WithIsolationSerializable() SQLTxOpt {
	return func(sqlOpt *sql.TxOptions) {
		sqlOpt.Isolation = sql.LevelSerializable
	}
}

// WrapInTxReturning executes the given function inside a transaction and returns its result
func WrapInTxReturning[T any](ctx context.Context, dbh *HandlerDB, f func() (T, error), sqlOptions ...SQLTxOpt) (res T, err error) {
	_, err = dbh.BeginTx(ctx, sqlOptions...)
	if err != nil {
		return res, fmt.Errorf(">database._.WrapInTxReturning failed to start tx %w", err)
	}
	res, err = f()
	if err != nil {
		err = dbh.RollbackIfTxOrErr(err)
		return res, fmt.Errorf(">database._.WrapInTxReturning failed to rollback tx %w", err)
	}
	err = dbh.CommitIfTxOrErr()
	if err != nil {
		return res, fmt.Errorf(">database._.WrapInTxReturning failed to commit tx %w", err)
	}
	return res, nil
}

// WrapInTx executes the given function inside a transaction
func WrapInTx(ctx context.Context, dbh *HandlerDB, f func() error, sqlOptions ...SQLTxOpt) (err error) {
	_, err = dbh.BeginTx(ctx, sqlOptions...)
	if err != nil {
		return fmt.Errorf(">database._.WrapInTx %w", err)
	}
	err = f()
	if err != nil {
		err = dbh.RollbackIfTxOrErr(err)
		return fmt.Errorf(">database._.WrapInTx %w", err)
	}
	err = dbh.CommitIfTxOrErr()
	if err != nil {
		return fmt.Errorf(">database._.WrapInTx %w", err)
	}
	return nil
}

// useExistingTxOrWrap reuse existing transaction to execute given function or creates a transaction locally
func useExistingTxOrWrap(ctx context.Context, dbh *HandlerDB, f func() error, sqlOptions ...SQLTxOpt) (err error) {
	if dbh.tx == nil {
		return WrapInTx(ctx, dbh, f, sqlOptions...)
	}
	return f()
}

// useExistingTxOrWrap reuse existing transaction to execute given function or creates a transaction locally and returns the result
func useExistingTxOrWrapReturning[T any](ctx context.Context, dbh *HandlerDB, f func() (T, error), sqlOptions ...SQLTxOpt) (res T, err error) {
	if dbh.tx == nil {
		return WrapInTxReturning(ctx, dbh, f, sqlOptions...)
	}
	return f()
}

// BeginTx starts a new transaction on current DB handler with given options if provided
func (dbh *HandlerDB) BeginTx(ctx context.Context, sqlOptions ...SQLTxOpt) (*goqu.TxDatabase, error) {
	txOpt := &sql.TxOptions{}
	for _, optFunc := range sqlOptions {
		optFunc(txOpt)
	}
	if dbh.tx != nil {
		return nil, NewDBTransactionError("there is already an active transaction, will not start new tx", nil, txOpt)
	}
	tx, err := dbh.db.BeginTx(ctx, txOpt)
	if err != nil {
		return nil, NewDBTransactionError("failed to begin transaction", err, txOpt)
	}
	dbh.tx = tx
	return tx, nil
}

// CommitIfTxOrErr attempts to commit transaction if any, else errors
func (dbh *HandlerDB) CommitIfTxOrErr() error {
	if dbh.tx != nil {
		err := dbh.tx.Commit()
		dbh.tx = nil
		if err != nil {
			return NewDBTransactionError("failed to commit database transaction", err, nil)
		}
		return nil
	}
	return ErrNoActiveTransaction
}

// RollbackIfTxOrErr attempts to rollback transaction if any, else errors
func (dbh *HandlerDB) RollbackIfTxOrErr(initialError error) error {
	if dbh.tx != nil {
		err := dbh.tx.Rollback()
		dbh.tx = nil
		if err != nil {
			return NewDBTransactionError(fmt.Sprintf("failed to rollback transaction [%s]", err.Error()), initialError, nil)
		}
		return initialError
	}
	return ErrNoActiveTransaction
}

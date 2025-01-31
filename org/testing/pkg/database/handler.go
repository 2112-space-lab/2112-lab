package database

import (
	"context"
	"database/sql"

	"github.com/doug-martin/goqu/v9"
)

// HandlerDB wraps DB and TX and provide convenience proxy methods to query from DB or from TX if any
type HandlerDB struct {
	db *goqu.Database
	tx *goqu.TxDatabase
}

// NewHandlerDB create a new HandlerDB from goqu database connection
func NewHandlerDB(db *goqu.Database) *HandlerDB {
	handler := HandlerDB{
		db: db,
		tx: nil,
	}
	return &handler
}

// Update helper to query against transaction if any, else on DB
func (dbh *HandlerDB) Update(table interface{}) *goqu.UpdateDataset {
	tableStr, ok := table.(string)
	if ok {
		table = goqu.I(tableStr)
	}
	if dbh.tx != nil {
		return dbh.tx.Update(table)
	}
	return dbh.db.Update(table)
}

// From helper to query against transaction if any, else on DB
func (dbh *HandlerDB) From(table interface{}) *goqu.SelectDataset {
	tableStr, ok := table.(string)
	if ok {
		table = goqu.I(tableStr)
	}
	if dbh.tx != nil {
		return dbh.tx.From(table)
	}
	return dbh.db.From(table)
}

// Insert helper to query against transaction if any, else on DB
func (dbh *HandlerDB) Insert(table interface{}) *goqu.InsertDataset {
	tableStr, ok := table.(string)
	if ok {
		table = goqu.I(tableStr)
	}
	if dbh.tx != nil {
		return dbh.tx.Insert(table)
	}
	return dbh.db.Insert(table)
}

// Delete helper to query against transaction if any, else on DB
func (dbh *HandlerDB) Delete(table interface{}) *goqu.DeleteDataset {
	tableStr, ok := table.(string)
	if ok {
		table = goqu.I(tableStr)
	}
	if dbh.tx != nil {
		return dbh.tx.Delete(table)
	}
	return dbh.db.Delete(table)
}

// ScanStructContext helper to query against transaction if any, else on DB
func (dbh *HandlerDB) ScanStructContext(ctx context.Context, i interface{}, query string, args ...interface{}) (bool, error) {
	if dbh.tx != nil {
		return dbh.tx.ScanStructContext(ctx, i, query, args...)
	}
	return dbh.db.ScanStructContext(ctx, i, query, args...)
}

// ScanStructsContext helper to query against transaction if any, else on DB
func (dbh *HandlerDB) ScanStructsContext(ctx context.Context, i interface{}, query string, args ...interface{}) error {
	if dbh.tx != nil {
		return dbh.tx.ScanStructsContext(ctx, i, query, args...)
	}
	return dbh.db.ScanStructsContext(ctx, i, query, args...)
}

// QueryContext helper to query against transaction if any, else on DB
func (dbh *HandlerDB) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	if dbh.tx != nil {
		return dbh.tx.QueryContext(ctx, query, args...)
	}
	return dbh.db.QueryContext(ctx, query, args...)
}

// ExecContext helper to query against transaction if any, else on DB
func (dbh *HandlerDB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	if dbh.tx != nil {
		return dbh.tx.ExecContext(ctx, query, args...)
	}
	return dbh.db.ExecContext(ctx, query, args...)
}

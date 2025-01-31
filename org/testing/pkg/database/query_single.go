package database

import (
	"context"
	"database/sql"

	"github.com/doug-martin/goqu/v9"
	"github.com/org/2112-space-lab/org/testing/pkg/fx"
)

// QuerySingleFromDataset type safe helper to execute selecting single row
func QuerySingleFromDataset[T any](ctx context.Context, dbh *HandlerDB, scanner func(*sql.Rows) (T, error), ds *goqu.SelectDataset) (res T, err error) {
	var empty T
	if ds.GetClauses().IsDefaultSelect() {
		ds = ds.Select(res)
	}
	sqlQuery, sqlArgs, sqlErr := ds.Prepared(true).ToSQL()
	if sqlErr != nil {
		return empty, NewDBQueryError(">database._.QuerySingleFromDataset", sqlErr, "", nil)
	}
	return QuerySingle(ctx, dbh, scanner, sqlQuery, sqlArgs...)
}

// QuerySingle type safe helper to execute selecting single row
func QuerySingle[T any](ctx context.Context, dbh *HandlerDB, scanner func(*sql.Rows) (T, error), sqlQuery string, sqlArgs ...interface{}) (res T, err error) {
	var sqlRows *sql.Rows

	sqlRows, err = dbh.QueryContext(ctx, sqlQuery, sqlArgs...)
	if err != nil {
		return res, NewDBQueryError(">database._.QuerySingle", err, sqlQuery, sqlArgs...)
	}
	rows, err := mapRows(sqlRows, err, scanner)
	if err != nil {
		return res, NewDBQueryError(">database._.QuerySingle", err, sqlQuery, sqlArgs...)
	}
	count := len(rows)
	if count == 0 {
		return res, NewDBQueryError(">database._.QuerySingle", ErrNotFoundOnQuerySingle, sqlQuery, sqlArgs...)
	}
	if count > 1 {
		return res, NewDBQueryError(">database._.QuerySingle", ErrMatchManyOnQuerySingle, sqlQuery, sqlArgs...)
	}
	return rows[0], nil
}

// ScanSingleFromDataset type safe helper to execute selecting single row
func ScanSingleFromDataset[T any](ctx context.Context, dbh *HandlerDB, ds *goqu.SelectDataset) (res T, err error) {
	if ds.GetClauses().IsDefaultSelect() {
		ds = ds.Select(res)
	}
	sqlQuery, sqlArgs, sqlErr := ds.Prepared(true).ToSQL()
	if sqlErr != nil {
		return res, NewDBQueryError(">database._.ScanSingleFromDataset", sqlErr, "", err)
	}
	return ScanSingle[T](ctx, dbh, sqlQuery, sqlArgs...)
}

// ScanSingleOptionFromDataset type safe helper to execute selecting single row
func ScanSingleOptionFromDataset[T any](ctx context.Context, dbh *HandlerDB, ds *goqu.SelectDataset) (res fx.Option[T], err error) {
	if ds.GetClauses().IsDefaultSelect() {
		ds = ds.Select(res.Value)
	}
	sqlQuery, sqlArgs, sqlErr := ds.Prepared(true).ToSQL()
	if sqlErr != nil {
		return res, NewDBQueryError(">database._.ScanSingleFromDataset", sqlErr, "", err)
	}
	return ScanSingleOption[T](ctx, dbh, sqlQuery, sqlArgs...)
}

// ScanSingle type safe helper to execute selecting single row
func ScanSingle[T any](ctx context.Context, dbh *HandlerDB, sqlQuery string, sqlArgs ...interface{}) (res T, err error) {
	var entities []T
	err = dbh.ScanStructsContext(ctx, &entities, sqlQuery, sqlArgs...)
	if err != nil {
		return res, NewDBQueryError(">database._.ScanSingle", err, sqlQuery, sqlArgs...)
	}
	count := len(entities)
	if count == 0 {
		return res, NewDBQueryError(">database._.ScanSingle", ErrNotFoundOnQuerySingle, sqlQuery, sqlArgs...)
	}
	if count > 1 {
		return res, NewDBQueryError(">database._.ScanSingle", ErrMatchManyOnQuerySingle, sqlQuery, sqlArgs...)
	}
	return entities[0], nil
}

// ScanSingleOption type safe helper to execute selecting single row
func ScanSingleOption[T any](ctx context.Context, dbh *HandlerDB, sqlQuery string, sqlArgs ...interface{}) (res fx.Option[T], err error) {
	var entities []T
	err = dbh.ScanStructsContext(ctx, &entities, sqlQuery, sqlArgs...)
	if err != nil {
		return res, NewDBQueryError(">database._.ScanSingle", err, sqlQuery, sqlArgs...)
	}
	count := len(entities)
	if count == 0 {
		return fx.NewEmptyOption[T](), nil
	}
	if count > 1 {
		return res, NewDBQueryError(">database._.ScanSingle", ErrMatchManyOnQuerySingle, sqlQuery, sqlArgs...)
	}
	return fx.NewValueOption(entities[0]), nil
}

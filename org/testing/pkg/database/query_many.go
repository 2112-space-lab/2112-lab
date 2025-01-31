package database

import (
	"context"
	"database/sql"

	"github.com/doug-martin/goqu/v9"
	"github.com/org/2112-space-lab/org/testing/pkg/fx"
)

// QueryManyFromDataset type safe helper to execute selecting many rows
func QueryManyFromDataset[T any](ctx context.Context, dbh *HandlerDB, scanner func(*sql.Rows) (T, error), ds *goqu.SelectDataset) (res []T, err error) {
	if ds.GetClauses().IsDefaultSelect() {
		ds = ds.Select(res)
	}
	sqlQuery, sqlArgs, sqlErr := ds.Prepared(true).ToSQL()
	if sqlErr != nil {
		return nil, NewDBQueryError(">database._.QueryManyFromDataset", sqlErr, "", nil)
	}
	return QueryMany(ctx, dbh, scanner, sqlQuery, sqlArgs...)
}

// QueryMany type safe helper to execute selecting many rows
func QueryMany[T any](ctx context.Context, dbh *HandlerDB, scanner func(*sql.Rows) (T, error), sqlQuery string, sqlArgs ...interface{}) (res []T, err error) {
	var sqlRows *sql.Rows
	sqlRows, err = dbh.QueryContext(ctx, sqlQuery, sqlArgs...)
	if err != nil {
		return res, NewDBQueryError(">database._.QueryMany", err, sqlQuery, sqlArgs...)
	}
	rows, err := mapRows(sqlRows, err, scanner)
	if err != nil {
		return res, NewDBQueryError(">database._.QueryMany", err, sqlQuery, sqlArgs...)
	}
	return rows, nil
}

// ScanManyFromDataset type safe helper to execute selecting many rows
func ScanManyFromDataset[T any](ctx context.Context, dbh *HandlerDB, ds *goqu.SelectDataset) (res []T, err error) {
	if ds.GetClauses().IsDefaultSelect() {
		ds = ds.Select(res)
	}
	sqlQuery, sqlArgs, sqlErr := ds.Prepared(true).ToSQL()
	if sqlErr != nil {
		return res, NewDBQueryError(">database._.ScanManyFromDataset", sqlErr, "", nil)
	}
	return ScanMany[T](ctx, dbh, sqlQuery, sqlArgs...)
}

// ScanManyOptionFromDataset type safe helper to execute selecting many rows
func ScanManyOptionFromDataset[T any](ctx context.Context, dbh *HandlerDB, ds *goqu.SelectDataset) (res fx.Option[[]T], err error) {
	if ds.GetClauses().IsDefaultSelect() {
		ds = ds.Select(res.Value)
	}
	sqlQuery, sqlArgs, sqlErr := ds.Prepared(true).ToSQL()
	if sqlErr != nil {
		return res, NewDBQueryError(">database._.ScanManyOptionFromDataset", sqlErr, "", nil)
	}
	return ScanManyOption[T](ctx, dbh, sqlQuery, sqlArgs...)
}

// ScanMany type safe helper to execute selecting many rows
func ScanMany[T any](ctx context.Context, dbh *HandlerDB, sqlQuery string, sqlArgs ...interface{}) (res []T, err error) {
	var entities []T
	err = dbh.ScanStructsContext(ctx, &entities, sqlQuery, sqlArgs...)
	if err != nil {
		return res, NewDBQueryError(">database._.ScanMany", err, sqlQuery, sqlArgs...)
	}
	return entities, nil
}

// ScanManyOption type safe helper to execute selecting many rows
func ScanManyOption[T any](ctx context.Context, dbh *HandlerDB, sqlQuery string, sqlArgs ...interface{}) (res fx.Option[[]T], err error) {
	var entities []T
	err = dbh.ScanStructsContext(ctx, &entities, sqlQuery, sqlArgs...)
	if err != nil {
		return res, NewDBQueryError(">database._.ScanSingle", err, sqlQuery, sqlArgs...)
	}
	count := len(entities)
	if count == 0 {
		return fx.NewEmptyOption[[]T](), nil
	}
	return fx.NewValueOption(entities), nil
}

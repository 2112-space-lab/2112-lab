package database

import (
	"context"
	"database/sql"

	"github.com/doug-martin/goqu/v9"
)

// InsertManyFromDataset type safe helper to execute updating multiple rows
func InsertManyFromDataset(ctx context.Context, dbh *HandlerDB, ds *goqu.InsertDataset) error {
	sqlQuery, sqlArgs, sqlErr := ds.Prepared(true).ToSQL()
	if sqlErr != nil {
		return NewDBQueryError(">database._.InsertManyFromDataset", sqlErr, sqlQuery, sqlArgs...)
	}
	return InsertMany(ctx, dbh, sqlQuery, sqlArgs...)
}

// InsertManyOptionFromDataset type safe helper to execute updating multiple rows
func InsertManyOptionFromDataset(ctx context.Context, dbh *HandlerDB, ds goqu.InsertDataset) error {
	sqlQuery, sqlArgs, sqlErr := ds.Prepared(true).ToSQL()
	if sqlErr != nil {
		return NewDBQueryError(">database._.InsertManyOptionFromDataset", sqlErr, sqlQuery, sqlArgs...)
	}
	return InsertManyOption(ctx, dbh, sqlQuery, sqlArgs...)
}

// InsertMany type safe helper to execute updating multiple rows
func InsertMany(ctx context.Context, dbh *HandlerDB, sqlQuery string, sqlArgs ...interface{}) (err error) {
	err = useExistingTxOrWrap(ctx, dbh, func() error {
		sqlRes, err := dbh.ExecContext(ctx, sqlQuery, sqlArgs...)
		if err != nil {
			return NewDBQueryError(">database._.InsertMany", err, sqlQuery, sqlArgs...)
		}
		rowsAffected, err := sqlRes.RowsAffected()
		if err != nil {
			return NewDBQueryError(">database._.InsertMany", err, sqlQuery, sqlArgs...)
		}
		if rowsAffected == 0 {
			return NewDBQueryError(">database._.InsertMany", ErrNoAffectedRows, sqlQuery, sqlArgs...)
		}
		return nil
	})
	return err
}

// InsertManyOption type safe helper to execute updating multiple rows
func InsertManyOption(ctx context.Context, dbh *HandlerDB, sqlQuery string, sqlArgs ...interface{}) (err error) {
	err = useExistingTxOrWrap(ctx, dbh, func() error {
		_, err := dbh.ExecContext(ctx, sqlQuery, sqlArgs...)
		if err != nil {
			return NewDBQueryError(">database._.InsertMany", err, sqlQuery, sqlArgs...)
		}
		return nil
	})
	return err
}

// InsertManyReturningQueryFromDataset type safe helper to execute updating multiple rows
func InsertManyReturningQueryFromDataset[T any](ctx context.Context, dbh *HandlerDB, scanner func(*sql.Rows) (T, error), ds *goqu.InsertDataset) (res []T, err error) {
	if ds.GetClauses().HasReturning() {
		ds = ds.Returning(res)
	}
	sqlQuery, sqlArgs, sqlErr := ds.Prepared(true).ToSQL()
	if sqlErr != nil {
		return res, NewDBQueryError(">database._.InsertManyReturningQueryFromDataset", sqlErr, "", nil)
	}
	res, err = useExistingTxOrWrapReturning(ctx, dbh, func() ([]T, error) {
		innerRes, innerErr := QueryMany(ctx, dbh, scanner, sqlQuery, sqlArgs...)
		if innerErr != nil {
			return res, NewDBQueryError(">database._.InsertManyReturningScanFromDataset", innerErr, "", nil)
		}
		return innerRes, nil
	})
	return res, err
}

// InsertManyReturningFromDataset type safe helper to execute updating multiple rows
func InsertManyReturningFromDataset[T any](ctx context.Context, dbh *HandlerDB, ds *goqu.InsertDataset) (res []T, err error) {
	if ds.GetClauses().HasReturning() {
		ds = ds.Returning(res)
	}
	sqlQuery, sqlArgs, sqlErr := ds.Prepared(true).ToSQL()
	if sqlErr != nil {
		return res, NewDBQueryError(">database._.InsertManyReturningScanFromDataset", sqlErr, "", nil)
	}
	res, err = useExistingTxOrWrapReturning(ctx, dbh, func() ([]T, error) {
		innerRes, innerErr := ScanMany[T](ctx, dbh, sqlQuery, sqlArgs...)
		if err != nil {
			return innerRes, NewDBQueryError(">database._.InsertManyReturningScanFromDataset", innerErr, "", nil)
		}
		return innerRes, nil
	})
	return res, err
}

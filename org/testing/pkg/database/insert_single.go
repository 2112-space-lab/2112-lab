package database

import (
	"context"
	"database/sql"

	"github.com/doug-martin/goqu/v9"
)

// InsertSingleFromDataset type safe helper to execute updating single row
func InsertSingleFromDataset(ctx context.Context, dbh *HandlerDB, ds *goqu.InsertDataset) error {
	sqlQuery, sqlArgs, sqlErr := ds.Prepared(true).ToSQL()
	if sqlErr != nil {
		return NewDBQueryError(">database._.InsertSingleFromDataset", sqlErr, "", nil)
	}
	return InsertSingle(ctx, dbh, sqlQuery, sqlArgs...)
}

// InsertSingle type safe helper to execute updating single row
func InsertSingle(ctx context.Context, dbh *HandlerDB, sqlQuery string, sqlArgs ...interface{}) (err error) {
	err = useExistingTxOrWrap(ctx, dbh, func() error {
		sqlRes, err := dbh.ExecContext(ctx, sqlQuery, sqlArgs...)
		if err != nil {
			return NewDBQueryError(">database._.InsertSingle", err, sqlQuery, sqlArgs...)
		}
		rowsAffected, err := sqlRes.RowsAffected()
		if err != nil {
			return NewDBQueryError(">database._.InsertSingle", err, sqlQuery, sqlArgs...)
		}
		if rowsAffected == 0 {
			return NewDBQueryError(">database._.InsertSingle", ErrNoAffectedRows, sqlQuery, sqlArgs...)
		}
		if rowsAffected > 1 {
			return NewDBQueryError(">database._.InsertSingle", ErrMoreThanOneAffectedRows, sqlQuery, sqlArgs...)
		}
		return nil
	})
	return err
}

// InsertSingleReturningQueryFromDataset type safe helper to execute updating single row
func InsertSingleReturningQueryFromDataset[T any](ctx context.Context, dbh *HandlerDB, scanner func(*sql.Rows) (T, error), ds *goqu.InsertDataset) (res T, err error) {
	if !ds.GetClauses().HasReturning() {
		ds = ds.Returning(res)
	}
	sqlQuery, sqlArgs, sqlErr := ds.Prepared(true).ToSQL()
	if sqlErr != nil {
		return res, NewDBQueryError(">database._.InsertSingleReturningQueryFromDataset", sqlErr, "", nil)
	}
	res, err = useExistingTxOrWrapReturning(ctx, dbh, func() (T, error) {
		innerRes, innerErr := QuerySingle(ctx, dbh, scanner, sqlQuery, sqlArgs...)
		if innerErr != nil {
			return innerRes, NewDBQueryError(">database._.InsertSingleReturningQueryFromDataset", innerErr, sqlQuery, sqlArgs...)
		}
		return innerRes, nil
	})
	return res, err
}

// InsertSingleReturningScanFromDataset type safe helper to execute updating single row
func InsertSingleReturningScanFromDataset[T any](ctx context.Context, dbh *HandlerDB, ds *goqu.InsertDataset) (res T, err error) {
	if !ds.GetClauses().HasReturning() {
		ds = ds.Returning(res)
	}
	sqlQuery, sqlArgs, sqlErr := ds.Prepared(true).ToSQL()
	if sqlErr != nil {
		return res, NewDBQueryError(">database._.InsertSingleReturningScanFromDataset", sqlErr, "", nil)
	}
	res, err = useExistingTxOrWrapReturning(ctx, dbh, func() (T, error) {
		innerRes, innerErr := ScanSingle[T](ctx, dbh, sqlQuery, sqlArgs...)
		if innerErr != nil {
			return innerRes, NewDBQueryError(">database._.InsertSingleReturningScanFromDataset", innerErr, sqlQuery, sqlArgs...)
		}
		return innerRes, nil
	})
	return res, err
}

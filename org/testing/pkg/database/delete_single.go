package database

import (
	"context"
	"database/sql"

	"github.com/doug-martin/goqu/v9"
)

// DeleteSingleFromDataset type safe helper to execute updating single row
func DeleteSingleFromDataset(ctx context.Context, dbh *HandlerDB, ds *goqu.DeleteDataset) error {
	sqlQuery, sqlArgs, sqlErr := ds.Prepared(true).ToSQL()
	if sqlErr != nil {
		return NewDBQueryError(">database._.DeleteSingleFromDataset", sqlErr, "", nil)
	}
	return DeleteSingle(ctx, dbh, sqlQuery, sqlArgs...)
}

// DeleteSingle type safe helper to execute updating single row
func DeleteSingle(ctx context.Context, dbh *HandlerDB, sqlQuery string, sqlArgs ...interface{}) (err error) {
	err = useExistingTxOrWrap(ctx, dbh, func() error {
		sqlRes, err := dbh.ExecContext(ctx, sqlQuery, sqlArgs...)
		if err != nil {
			return NewDBQueryError(">database._.DeleteSingle", err, sqlQuery, sqlArgs...)
		}
		rowsAffected, err := sqlRes.RowsAffected()
		if err != nil {
			return NewDBQueryError(">database._.DeleteSingle", err, sqlQuery, sqlArgs...)
		}
		if rowsAffected == 0 {
			return NewDBQueryError(">database._.DeleteSingle", ErrNoAffectedRows, sqlQuery, sqlArgs...)
		}
		if rowsAffected > 1 {
			return NewDBQueryError(">database._.DeleteSingle", ErrMoreThanOneAffectedRows, sqlQuery, sqlArgs...)
		}
		return nil
	})
	return err
}

// DeleteSingleReturningQueryFromDataset type safe helper to execute updating single row
func DeleteSingleReturningQueryFromDataset[T any](ctx context.Context, dbh *HandlerDB, scanner func(*sql.Rows) (T, error), ds *goqu.DeleteDataset) (res T, err error) {
	if !ds.GetClauses().HasReturning() {
		ds = ds.Returning(res)
	}
	sqlQuery, sqlArgs, sqlErr := ds.Prepared(true).ToSQL()
	if sqlErr != nil {
		return res, NewDBQueryError(">database._.DeleteSingleReturningQueryFromDataset", sqlErr, "", nil)
	}
	res, err = useExistingTxOrWrapReturning(ctx, dbh, func() (T, error) {
		innerRes, innerErr := QuerySingle(ctx, dbh, scanner, sqlQuery, sqlArgs...)
		if innerErr != nil {
			return innerRes, NewDBQueryError(">database._.DeleteSingleReturningQueryFromDataset", innerErr, sqlQuery, sqlArgs...)
		}
		return innerRes, nil
	})
	return res, err
}

// DeleteSingleReturningScanFromDataset type safe helper to execute updating single row
func DeleteSingleReturningScanFromDataset[T any](ctx context.Context, dbh *HandlerDB, ds *goqu.DeleteDataset) (res T, err error) {
	if !ds.GetClauses().HasReturning() {
		ds = ds.Returning(res)
	}
	sqlQuery, sqlArgs, sqlErr := ds.Prepared(true).ToSQL()
	if sqlErr != nil {
		return res, NewDBQueryError(">database._.DeleteSingleReturningScanFromDataset", sqlErr, "", nil)
	}
	res, err = useExistingTxOrWrapReturning(ctx, dbh, func() (T, error) {
		innerRes, innerErr := ScanSingle[T](ctx, dbh, sqlQuery, sqlArgs...)
		if innerErr != nil {
			return innerRes, NewDBQueryError(">database._.DeleteSingleReturningScanFromDataset", innerErr, sqlQuery, sqlArgs...)
		}
		return innerRes, nil
	})
	return res, err
}

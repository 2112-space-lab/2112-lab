package database

import (
	"context"
	"database/sql"

	"github.com/doug-martin/goqu/v9"
)

// UpdateSingleFromDataset type safe helper to execute updating single row
func UpdateSingleFromDataset(ctx context.Context, dbh *HandlerDB, ds *goqu.UpdateDataset) error {
	sqlQuery, sqlArgs, sqlErr := ds.Prepared(true).ToSQL()
	if sqlErr != nil {
		return NewDBQueryError(">database._.UpdateSingleFromDataset", sqlErr, "", nil)
	}
	return UpdateSingle(ctx, dbh, sqlQuery, sqlArgs...)
}

// UpdateSingle type safe helper to execute updating single row
func UpdateSingle(ctx context.Context, dbh *HandlerDB, sqlQuery string, sqlArgs ...interface{}) (err error) {
	err = useExistingTxOrWrap(ctx, dbh, func() error {
		sqlRes, err := dbh.ExecContext(ctx, sqlQuery, sqlArgs...)
		if err != nil {
			return NewDBQueryError(">database._.UpdateSingle", err, sqlQuery, sqlArgs...)
		}
		rowsAffected, err := sqlRes.RowsAffected()
		if err != nil {
			return NewDBQueryError(">database._.UpdateSingle", err, sqlQuery, sqlArgs...)
		}
		if rowsAffected == 0 {
			return NewDBQueryError(">database._.UpdateSingle", ErrNoAffectedRows, sqlQuery, sqlArgs...)
		}
		if rowsAffected > 1 {
			return NewDBQueryError(">database._.UpdateSingle", ErrMoreThanOneAffectedRows, sqlQuery, sqlArgs...)
		}
		return nil
	})
	return err
}

// UpdateSingleOptionFromDataset type safe helper to execute updating single row
func UpdateSingleOptionFromDataset(ctx context.Context, dbh *HandlerDB, ds *goqu.UpdateDataset) error {
	sqlQuery, sqlArgs, sqlErr := ds.Prepared(true).ToSQL()
	if sqlErr != nil {
		return NewDBQueryError(">database._.UpdateSingleOptionFromDataset", sqlErr, "", nil)
	}
	return UpdateSingleOption(ctx, dbh, sqlQuery, sqlArgs...)
}

// UpdateSingleOption type safe helper to execute updating single row
func UpdateSingleOption(ctx context.Context, dbh *HandlerDB, sqlQuery string, sqlArgs ...interface{}) (err error) {
	err = useExistingTxOrWrap(ctx, dbh, func() error {
		sqlRes, err := dbh.ExecContext(ctx, sqlQuery, sqlArgs...)
		if err != nil {
			return NewDBQueryError(">database._.UpdateSingleOption", err, sqlQuery, sqlArgs...)
		}
		rowsAffected, err := sqlRes.RowsAffected()
		if err != nil {
			return NewDBQueryError(">database._.UpdateSingleOption", err, sqlQuery, sqlArgs...)
		}
		if rowsAffected > 1 {
			return NewDBQueryError(">database._.UpdateSingleOption", ErrMoreThanOneAffectedRows, sqlQuery, sqlArgs...)
		}
		return nil
	})
	return err
}

// UpdateSingleReturningQueryFromDataset type safe helper to execute updating single row
func UpdateSingleReturningQueryFromDataset[T any](ctx context.Context, dbh *HandlerDB, scanner func(*sql.Rows) (T, error), ds *goqu.UpdateDataset) (res T, err error) {
	if !ds.GetClauses().HasReturning() {
		ds = ds.Returning(res)
	}
	sqlQuery, sqlArgs, sqlErr := ds.Prepared(true).ToSQL()
	if sqlErr != nil {
		return res, NewDBQueryError(">database._.UpdateSingleReturningQueryFromDataset", sqlErr, "", nil)
	}
	res, err = useExistingTxOrWrapReturning(ctx, dbh, func() (T, error) {
		innerRes, innerErr := QuerySingle(ctx, dbh, scanner, sqlQuery, sqlArgs...)
		if innerErr != nil {
			return innerRes, NewDBQueryError(">database._.UpdateSingleReturningQueryFromDataset", innerErr, sqlQuery, sqlArgs...)
		}
		return innerRes, nil
	})
	return res, err
}

// UpdateSingleReturningScanFromDataset type safe helper to execute updating single row
func UpdateSingleReturningScanFromDataset[T any](ctx context.Context, dbh *HandlerDB, ds *goqu.UpdateDataset) (res T, err error) {
	if !ds.GetClauses().HasReturning() {
		ds = ds.Returning(res)
	}
	sqlQuery, sqlArgs, sqlErr := ds.Prepared(true).ToSQL()
	if sqlErr != nil {
		return res, NewDBQueryError(">database._.UpdateSingleReturningScanFromDataset", sqlErr, "", nil)
	}
	res, err = useExistingTxOrWrapReturning(ctx, dbh, func() (T, error) {
		innerRes, innerErr := ScanSingle[T](ctx, dbh, sqlQuery, sqlArgs...)
		if innerErr != nil {
			return innerRes, NewDBQueryError(">database._.UpdateSingleReturningScanFromDataset", innerErr, sqlQuery, sqlArgs...)
		}
		return innerRes, nil
	})
	return res, err
}

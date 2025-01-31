package database

import (
	"context"
	"database/sql"

	"github.com/doug-martin/goqu/v9"
)

// UpdateManyFromDataset type safe helper to execute updating multiple rows
func UpdateManyFromDataset(ctx context.Context, dbh *HandlerDB, ds *goqu.UpdateDataset) error {
	sqlQuery, sqlArgs, sqlErr := ds.Prepared(true).ToSQL()
	if sqlErr != nil {
		return NewDBQueryError(">database._.UpdateManyFromDataset", sqlErr, "", nil)
	}
	return UpdateMany(ctx, dbh, sqlQuery, sqlArgs...)
}

// UpdateManyIfAnyFromDataset type safe helper to execute updating multiple rows
func UpdateManyIfAnyFromDataset(ctx context.Context, dbh *HandlerDB, ds *goqu.UpdateDataset) error {
	sqlQuery, sqlArgs, sqlErr := ds.Prepared(true).ToSQL()
	if sqlErr != nil {
		return NewDBQueryError(">database._.UpdateManyIfAnyFromDataset", sqlErr, "", nil)
	}
	return UpdateManyIfAny(ctx, dbh, sqlQuery, sqlArgs...)
}

// UpdateMany type safe helper to execute updating multiple rows
func UpdateMany(ctx context.Context, dbh *HandlerDB, sqlQuery string, sqlArgs ...interface{}) (err error) {
	err = useExistingTxOrWrap(ctx, dbh, func() error {
		sqlRes, err := dbh.ExecContext(ctx, sqlQuery, sqlArgs...)
		if err != nil {
			return NewDBQueryError(">database._.UpdateMany", err, sqlQuery, sqlArgs...)
		}
		rowsAffected, err := sqlRes.RowsAffected()
		if err != nil {
			return NewDBQueryError(">database._.UpdateMany", err, sqlQuery, sqlArgs...)
		}
		if rowsAffected == 0 {
			return NewDBQueryError(">database._.UpdateMany", ErrNoAffectedRows, sqlQuery, sqlArgs...)
		}
		return nil
	})
	return err
}

// UpdateManyIfAny type safe helper to execute updating multiple rows
func UpdateManyIfAny(ctx context.Context, dbh *HandlerDB, sqlQuery string, sqlArgs ...interface{}) (err error) {
	err = useExistingTxOrWrap(ctx, dbh, func() error {
		_, err := dbh.ExecContext(ctx, sqlQuery, sqlArgs...)
		if err != nil {
			return NewDBQueryError(">database._.UpdateMany", err, sqlQuery, sqlArgs...)
		}
		return nil
	})
	return err
}

// UpdateManyReturningQueryFromDataset type safe helper to execute updating multiple rows
func UpdateManyReturningQueryFromDataset[T any](ctx context.Context, dbh *HandlerDB, scanner func(*sql.Rows) (T, error), ds *goqu.UpdateDataset) (res []T, err error) {
	if !ds.GetClauses().HasReturning() {
		ds = ds.Returning(res)
	}
	sqlQuery, sqlArgs, sqlErr := ds.Prepared(true).ToSQL()
	if sqlErr != nil {
		return res, NewDBQueryError(">database._.UpdateManyReturningQueryFromDataset", sqlErr, "", nil)
	}
	res, err = useExistingTxOrWrapReturning(ctx, dbh, func() ([]T, error) {
		innerRes, innerErr := QueryMany(ctx, dbh, scanner, sqlQuery, sqlArgs...)
		if innerErr != nil {
			return res, NewDBQueryError(">database._.UpdateManyReturningScanFromDataset", innerErr, "", nil)
		}
		return innerRes, nil
	})
	return res, err
}

// UpdateManyReturningFromDataset type safe helper to execute updating multiple rows
func UpdateManyReturningFromDataset[T any](ctx context.Context, dbh *HandlerDB, ds *goqu.UpdateDataset) (res []T, err error) {
	if ds.GetClauses().HasReturning() {
		ds = ds.Returning(res)
	}
	sqlQuery, sqlArgs, sqlErr := ds.Prepared(true).ToSQL()
	if sqlErr != nil {
		return res, NewDBQueryError(">database._.UpdateManyReturningScanFromDataset", sqlErr, "", nil)
	}
	res, err = useExistingTxOrWrapReturning(ctx, dbh, func() ([]T, error) {
		innerRes, innerErr := ScanMany[T](ctx, dbh, sqlQuery, sqlArgs...)
		if err != nil {
			return innerRes, NewDBQueryError(">database._.UpdateManyReturningScanFromDataset", innerErr, "", nil)
		}
		return innerRes, nil
	})
	return res, err
}

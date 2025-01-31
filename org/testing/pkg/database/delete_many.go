package database

import (
	"context"
	"database/sql"

	"github.com/doug-martin/goqu/v9"
)

// DeleteManyFromDataset type safe helper to execute updating multiple rows
func DeleteManyFromDataset(ctx context.Context, dbh *HandlerDB, ds *goqu.DeleteDataset) error {
	sqlQuery, sqlArgs, sqlErr := ds.Prepared(true).ToSQL()
	if sqlErr != nil {
		return NewDBQueryError(">database._.DeleteManyFromDataset", sqlErr, "", nil)
	}
	return DeleteMany(ctx, dbh, sqlQuery, sqlArgs...)
}

// DeleteMany type safe helper to execute updating multiple rows
func DeleteMany(ctx context.Context, dbh *HandlerDB, sqlQuery string, sqlArgs ...interface{}) (err error) {
	err = useExistingTxOrWrap(ctx, dbh, func() error {
		sqlRes, err := dbh.ExecContext(ctx, sqlQuery, sqlArgs...)
		if err != nil {
			return NewDBQueryError(">database._.DeleteMany", err, sqlQuery, sqlArgs...)
		}
		rowsAffected, err := sqlRes.RowsAffected()
		if err != nil {
			return NewDBQueryError(">database._.DeleteMany", err, sqlQuery, sqlArgs...)
		}
		if rowsAffected == 0 {
			return NewDBQueryError(">database._.DeleteMany", ErrNoAffectedRows, sqlQuery, sqlArgs...)
		}
		return nil
	})
	return err
}

// DeleteManyReturningQueryFromDataset type safe helper to execute updating multiple rows
func DeleteManyReturningQueryFromDataset[T any](ctx context.Context, dbh *HandlerDB, scanner func(*sql.Rows) (T, error), ds *goqu.DeleteDataset) (res []T, err error) {
	if !ds.GetClauses().HasReturning() {
		ds = ds.Returning(res)
	}
	sqlQuery, sqlArgs, sqlErr := ds.Prepared(true).ToSQL()
	if sqlErr != nil {
		return res, NewDBQueryError(">database._.DeleteManyReturningQueryFromDataset", sqlErr, "", nil)
	}
	res, err = useExistingTxOrWrapReturning(ctx, dbh, func() ([]T, error) {
		innerRes, innerErr := QueryMany(ctx, dbh, scanner, sqlQuery, sqlArgs...)
		if innerErr != nil {
			return res, NewDBQueryError(">database._.DeleteManyReturningScanFromDataset", innerErr, "", nil)
		}
		return innerRes, nil
	})
	return res, err
}

// DeleteManyReturningFromDataset type safe helper to execute updating multiple rows
func DeleteManyReturningFromDataset[T any](ctx context.Context, dbh *HandlerDB, ds *goqu.DeleteDataset) (res []T, err error) {
	if ds.GetClauses().HasReturning() {
		ds = ds.Returning(res)
	}
	sqlQuery, sqlArgs, sqlErr := ds.Prepared(true).ToSQL()
	if sqlErr != nil {
		return res, NewDBQueryError(">database._.DeleteManyReturningScanFromDataset", sqlErr, "", nil)
	}
	res, err = useExistingTxOrWrapReturning(ctx, dbh, func() ([]T, error) {
		innerRes, innerErr := ScanMany[T](ctx, dbh, sqlQuery, sqlArgs...)
		if err != nil {
			return innerRes, NewDBQueryError(">database._.DeleteManyReturningScanFromDataset", innerErr, "", nil)
		}
		return innerRes, nil
	})
	return res, err
}

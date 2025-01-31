package database

import (
	"database/sql"
	"fmt"

	"github.com/org/2112-space-lab/org/testing/pkg/fx"
)

func mapRows[E any](sqlRows *sql.Rows, sqlErr error, scanner func(*sql.Rows) (E, error)) ([]E, error) {
	if sqlErr != nil {
		return nil, sqlErr
	}
	var entities []E
	var errs []error
	for sqlRows.Next() {
		e, err := scanner(sqlRows)
		if err != nil {
			errs = append(errs, err)
		}
		entities = append(entities, e)
	}
	err := sqlRows.Close()
	if err != nil {
		errs = append(errs, err)
	}
	err = sqlRows.Err()
	if err != nil {
		errs = append(errs, err)
	}
	if len(errs) > 0 {
		err := fx.FlattenErrorsIfAny(errs...)
		return nil, fmt.Errorf(">database._.mapRows %w", err)
	}
	return entities, nil
}

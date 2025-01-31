package database

import (
	"database/sql"
	"fmt"
)

// ScannerSingleValue is an utility function for general purpose usage with rows.Scan when we have a single column selected from query
func ScannerSingleValue[T any](r *sql.Rows) (T, error) {
	var id T
	err := r.Scan(&id)
	if err != nil {
		return id, fmt.Errorf("failed to scan id from row %w", err)
	}
	return id, nil
}
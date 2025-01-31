package database

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
)

var (
	// ErrNoActiveTransaction sentinel error value
	ErrNoActiveTransaction = NewDBTransactionError("attempted to get transaction but no active transaction on database handler", nil, nil)
	// ErrNotFoundOnQuerySingle sentinel error value
	ErrNotFoundOnQuerySingle = errors.New("no record found on matching query single")
	// ErrMatchManyOnQuerySingle sentinel error value
	ErrMatchManyOnQuerySingle = errors.New("more than 1 record matching query single")
	// ErrNoAffectedRows sentinel error value
	ErrNoAffectedRows = errors.New("no rows affected by query")
	// ErrMoreThanOneAffectedRows sentinel error value
	ErrMoreThanOneAffectedRows = errors.New("more than 1 row affected by query")
)

// DBErr is base type for all database errors
type DBErr struct {
	Inner   error
	Message string
}

// Error method is implemented for error interface, it converts error to string representation
func (e *DBErr) Error() string {
	if e.Inner != nil {
		return fmt.Sprintf("database error: %s: %v", e.Message, e.Inner)
	}
	return fmt.Sprintf("database error: %s", e.Message)
}

// Unwrap returns next error in chain
func (e *DBErr) Unwrap() error { return e.Inner }

// NewDBError creates new database error containing message and able to wrap inner error, message is given message for error,
// and inner is wrapped error.
func NewDBError(message string, inner error) *DBErr {
	return &DBErr{Message: message, Inner: inner}
}

// DBTransactionErr hold transaction error details
type DBTransactionErr struct {
	DBErr
	TxOpt *sql.TxOptions
}

// Error method is implemented for error interface, it converts error to string representation
func (e *DBTransactionErr) Error() string {
	txMsg := ""
	if e.TxOpt != nil {
		txMsg = fmt.Sprintf("with txOptions [%v] ", e.TxOpt)
	}
	if e.Inner != nil {
		return fmt.Sprintf("transaction error: %s %s[%v]", e.Message, txMsg, e.Inner)
	}
	return fmt.Sprintf("transaction error: %s %s", e.Message, txMsg)
}

// Unwrap returns next error in chain
func (e *DBTransactionErr) Unwrap() error { return e.Inner }

// NewDBTransactionError initializes a new DBTransactionErr
func NewDBTransactionError(msg string, inner error, txOpt *sql.TxOptions) *DBTransactionErr {
	return &DBTransactionErr{DBErr: DBErr{Message: msg, Inner: inner}, TxOpt: txOpt}
}

// DBQueryErr holds details of database query error
type DBQueryErr struct {
	DBErr
	SQLQuery string
	SQLArgs  []interface{}
}

// Error method is implemented for error interface, it converts error to string representation
func (e *DBQueryErr) Error() string {
	if e.Inner != nil {
		argsJSON, err := json.Marshal(e.SQLArgs)
		if err != nil {
			argsJSON = []byte("")
		}
		return fmt.Sprintf("query err: %s for SQL query [%s] with args [%s] - [%v]", e.Message, e.SQLQuery, argsJSON, e.Inner)
	}
	return fmt.Sprintf("query err: %s for SQL query [%s] with args [%v]", e.Message, e.SQLQuery, e.SQLArgs)
}

// Unwrap returns next error in chain
func (e *DBQueryErr) Unwrap() error { return e.Inner }

// NewDBQueryError initializes a new DBQueryErr
func NewDBQueryError(msg string, inner error, sqlQuery string, sqlArgs ...interface{}) *DBQueryErr {
	return &DBQueryErr{DBErr: DBErr{Message: msg, Inner: inner}, SQLQuery: sqlQuery, SQLArgs: sqlArgs}
}

package apperror

import (
	"errors"
	"strings"

	"github.com/afandimsr/go-gin-api/internal/config"
	mssql "github.com/denisenkom/go-mssqldb"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

// isDevelopment checks if the application is running in development mode
func isDevelopment() bool {
	return config.Load().AppEnv == "development" || gin.Mode() == gin.DebugMode
}

// HandleDatabaseError converts database errors to AppError (works with MySQL, PostgreSQL, SQLite, MS SQL Server)
func HandleDatabaseError(err error) *AppError {
	if err == nil {
		return nil
	}

	// GORM specific errors (database agnostic)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return NewNotFoundError("data not found")
	}

	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return NewConflictError("data already exists")
	}

	if errors.Is(err, gorm.ErrInvalidTransaction) {
		if isDevelopment() {
			return NewAppError(500, DatabaseError, err.Error(), err)
		}
		return NewAppError(500, DatabaseError, "invalid database transaction", err)
	}

	if errors.Is(err, gorm.ErrInvalidField) {
		if isDevelopment() {
			return NewAppError(500, DatabaseError, err.Error(), err)
		}
		return NewAppError(500, DatabaseError, "invalid database field", err)
	}

	// MySQL specific errors
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		return handleMySQLError(mysqlErr)
	}

	// PostgreSQL specific errors
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return handlePostgreSQLError(pgErr)
	}

	// MS SQL Server specific errors
	var mssqlErr mssql.Error
	if errors.As(err, &mssqlErr) {
		return handleMSSQLError(mssqlErr)
	}

	// Check for common error patterns in error message (fallback)
	return handleGenericDatabaseError(err)
}

// handleMySQLError handles MySQL specific error codes
func handleMySQLError(err *mysql.MySQLError) *AppError {
	// In development mode, show native error
	if isDevelopment() {
		return NewAppError(500, DatabaseError, err.Error(), err)
	}

	switch err.Number {
	case 1062: // Duplicate entry
		return NewConflictError("data already exists")

	case 1452: // Foreign key constraint fails (insert/update)
		return NewBadRequestError("cannot perform this operation due to related data")

	case 1451: // Cannot delete or update a parent row (foreign key)
		return NewBadRequestError("cannot delete data that is being used")

	case 1064: // SQL syntax error
		return NewAppError(500, DatabaseQueryError, "invalid database query", err)

	case 1146: // Table doesn't exist
		return NewAppError(500, DatabaseError, "database table not found", err)

	case 1054: // Unknown column
		return NewAppError(500, DatabaseError, "invalid database column", err)

	case 1048: // Column cannot be null
		return NewBadRequestError("required field is missing")

	case 1406: // Data too long for column
		return NewBadRequestError("input data exceeds maximum length")

	case 1366: // Incorrect value
		return NewBadRequestError("invalid data format")

	case 2002, 2003, 2006, 2013: // Connection errors
		return NewAppError(503, DatabaseConnectionError, "database connection failed", err)

	case 1205: // Lock wait timeout
		return NewAppError(503, DatabaseError, "database is busy, please try again", err)

	case 1213: // Deadlock
		return NewAppError(503, DatabaseError, "database deadlock detected, please try again", err)

	case 1040: // Too many connections
		return NewAppError(503, DatabaseConnectionError, "database connection limit reached", err)

	default:
		return NewAppError(500, DatabaseError, "database operation failed", err)
	}
}

// handlePostgreSQLError handles PostgreSQL specific error codes
func handlePostgreSQLError(err *pgconn.PgError) *AppError {
	// In development mode, show native error
	if isDevelopment() {
		return NewAppError(500, DatabaseError, err.Message, err)
	}

	switch err.Code {
	case "23505": // unique_violation
		return NewConflictError("data already exists")

	case "23503": // foreign_key_violation
		return NewBadRequestError("cannot perform this operation due to related data")

	case "23502": // not_null_violation
		return NewBadRequestError("required field is missing")

	case "23514": // check_violation
		return NewBadRequestError("data violates check constraint")

	case "22001": // string_data_right_truncation
		return NewBadRequestError("input data exceeds maximum length")

	case "22P02": // invalid_text_representation
		return NewBadRequestError("invalid data format")

	case "42601": // syntax_error
		return NewAppError(500, DatabaseQueryError, "invalid database query", err)

	case "42P01": // undefined_table
		return NewAppError(500, DatabaseError, "database table not found", err)

	case "42703": // undefined_column
		return NewAppError(500, DatabaseError, "invalid database column", err)

	case "08000", "08003", "08006": // connection_exception
		return NewAppError(503, DatabaseConnectionError, "database connection failed", err)

	case "40001": // serialization_failure
		return NewAppError(503, DatabaseError, "database transaction conflict, please try again", err)

	case "40P01": // deadlock_detected
		return NewAppError(503, DatabaseError, "database deadlock detected, please try again", err)

	case "53300": // too_many_connections
		return NewAppError(503, DatabaseConnectionError, "database connection limit reached", err)

	case "57014": // query_canceled
		return NewAppError(500, DatabaseError, "database query was canceled", err)

	default:
		return NewAppError(500, DatabaseError, "database operation failed", err)
	}
}

// handleMSSQLError handles MS SQL Server specific error codes
func handleMSSQLError(err mssql.Error) *AppError {
	// In development mode, show native error
	if isDevelopment() {
		return NewAppError(500, DatabaseError, err.Message, err)
	}

	switch err.Number {
	case 2627, 2601: // Unique constraint violation / Duplicate key
		return NewConflictError("data already exists")

	case 547: // Foreign key constraint violation
		return NewBadRequestError("cannot perform this operation due to related data")

	case 515: // Cannot insert NULL value
		return NewBadRequestError("required field is missing")

	case 8152: // String or binary data would be truncated
		return NewBadRequestError("input data exceeds maximum length")

	case 245, 241, 242, 257: // Conversion/Type errors
		return NewBadRequestError("invalid data format")

	case 102, 105, 156, 170: // Syntax errors
		return NewAppError(500, DatabaseQueryError, "invalid database query", err)

	case 208: // Invalid object name (table doesn't exist)
		return NewAppError(500, DatabaseError, "database table not found", err)

	case 207: // Invalid column name
		return NewAppError(500, DatabaseError, "invalid database column", err)

	case 4060, 18456: // Login/authentication failed
		return NewAppError(503, DatabaseConnectionError, "database authentication failed", err)

	case 64, 233, 10054, 10060, 10061: // Connection errors
		return NewAppError(503, DatabaseConnectionError, "database connection failed", err)

	case 1205: // Deadlock victim
		return NewAppError(503, DatabaseError, "database deadlock detected, please try again", err)

	case 1222: // Lock request timeout
		return NewAppError(503, DatabaseError, "database is busy, please try again", err)

	case 40197, 40501, 40613, 49918, 49919: // Resource/connection limit errors
		return NewAppError(503, DatabaseConnectionError, "database connection limit reached", err)

	case 3621: // Transaction ended in trigger
		return NewAppError(500, DatabaseError, "database transaction was rolled back", err)

	case 2812: // Stored procedure not found
		return NewAppError(500, DatabaseError, "database procedure not found", err)

	case 50000: // Custom error raised by RAISERROR (depends on severity)
		if err.Class >= 16 { // Error severity
			return NewAppError(500, DatabaseError, err.Message, err)
		}
		return NewAppError(400, DatabaseError, err.Message, err)

	default:
		// Handle by severity class
		switch {
		case err.Class >= 20: // Fatal errors (connection will be closed)
			return NewAppError(503, DatabaseConnectionError, "fatal database error", err)
		case err.Class >= 16: // Errors that can be corrected by user
			return NewAppError(500, DatabaseError, "database operation failed", err)
		case err.Class >= 11: // Informational errors
			return NewBadRequestError("database operation failed")
		default:
			return NewAppError(500, DatabaseError, "database operation failed", err)
		}
	}
}

// handleGenericDatabaseError handles errors by examining error message (fallback for unknown databases)
func handleGenericDatabaseError(err error) *AppError {
	errMsg := strings.ToLower(err.Error())

	// In development mode, show native error for unrecognized patterns
	if isDevelopment() {
		return NewAppError(500, DatabaseError, err.Error(), err)
	}

	// Duplicate/Unique constraint
	if strings.Contains(errMsg, "duplicate") ||
		strings.Contains(errMsg, "unique constraint") ||
		strings.Contains(errMsg, "unique violation") ||
		strings.Contains(errMsg, "unique index") {
		return NewConflictError("data already exists")
	}

	// Foreign key constraint
	if strings.Contains(errMsg, "foreign key") ||
		strings.Contains(errMsg, "constraint fails") ||
		strings.Contains(errMsg, "violates foreign key") ||
		strings.Contains(errMsg, "reference constraint") {
		return NewBadRequestError("cannot perform this operation due to related data")
	}

	// Not null constraint
	if strings.Contains(errMsg, "not null") ||
		strings.Contains(errMsg, "cannot be null") ||
		strings.Contains(errMsg, "null value") {
		return NewBadRequestError("required field is missing")
	}

	// Syntax error
	if strings.Contains(errMsg, "syntax") ||
		strings.Contains(errMsg, "parse error") {
		return NewAppError(500, DatabaseQueryError, "invalid database query", err)
	}

	// Table not found
	if strings.Contains(errMsg, "table") &&
		(strings.Contains(errMsg, "not found") ||
			strings.Contains(errMsg, "doesn't exist") ||
			strings.Contains(errMsg, "does not exist") ||
			strings.Contains(errMsg, "invalid object")) {
		return NewAppError(500, DatabaseError, "database table not found", err)
	}

	// Column not found
	if strings.Contains(errMsg, "column") &&
		(strings.Contains(errMsg, "not found") ||
			strings.Contains(errMsg, "doesn't exist") ||
			strings.Contains(errMsg, "unknown") ||
			strings.Contains(errMsg, "invalid")) {
		return NewAppError(500, DatabaseError, "invalid database column", err)
	}

	// Connection errors
	if strings.Contains(errMsg, "connection") ||
		strings.Contains(errMsg, "connect") ||
		strings.Contains(errMsg, "timeout") ||
		strings.Contains(errMsg, "network") {
		return NewAppError(503, DatabaseConnectionError, "database connection failed", err)
	}

	// Lock/Busy errors
	if strings.Contains(errMsg, "lock") ||
		strings.Contains(errMsg, "busy") ||
		strings.Contains(errMsg, "deadlock") {
		return NewAppError(503, DatabaseError, "database is busy, please try again", err)
	}

	// Data too long
	if strings.Contains(errMsg, "too long") ||
		strings.Contains(errMsg, "too big") ||
		strings.Contains(errMsg, "exceeds") ||
		strings.Contains(errMsg, "truncated") {
		return NewBadRequestError("input data exceeds maximum length")
	}

	// Generic database error
	return NewAppError(500, DatabaseError, "database operation failed", err)
}

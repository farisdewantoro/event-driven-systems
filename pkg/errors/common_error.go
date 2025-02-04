package errors

import "net/http"

// Predefined errors for common cases
var (
	ErrNotFound        = NewAPIError("ERR1001", http.StatusNotFound, "The requested resource could not be found.")
	ErrInternal        = NewAPIError("ERR1002", http.StatusInternalServerError, "An unexpected error occurred. Please try again later.")
	ErrBadRequest      = NewAPIError("ERR1003", http.StatusBadRequest, "The request could not be processed due to invalid input.")
	ErrUnauthorized    = NewAPIError("ERR1004", http.StatusUnauthorized, "You are not authorized to access this resource.")
	ErrBindRequest     = NewAPIError("ERR1005", http.StatusBadRequest, "Failed to bind the request data. Please ensure the format is correct.")
	ErrSQLCreate       = NewAPIError("ERR1006", http.StatusInternalServerError, "An error occurred while trying to create the record. Please try again.")
	ErrSQLGet          = NewAPIError("ERR1007", http.StatusInternalServerError, "An error occurred while retrieving the requested data. Please try again.")
	ErrSQLTx           = NewAPIError("ERR1008", http.StatusInternalServerError, "An error occurred while processing the transaction. Please try again.")
	ErrParseJsonOutbox = NewAPIError("ERR1009", http.StatusInternalServerError, "An error occurred while parsing the JSON data. Please try again.")
)

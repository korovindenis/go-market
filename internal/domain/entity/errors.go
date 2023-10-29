package entity

import "errors"

var (
	// ErrMetricNotFound            = errors.New("metric not set")
	// ErrEnvVarNotFound            = errors.New("env var not set")
	// ErrInvalidURLFormat          = errors.New("invalid URL format")
	ErrMethodNotAllowed     = errors.New("method not allowed")
	ErrUnsupportedMediaType = errors.New("unsupported media type")
	ErrInternalServerError  = errors.New("internal server error")
	// ErrInputVarIsWrongType       = errors.New("metric value is wrong type")
	ErrStatusBadRequest                = errors.New("bad request")
	ErrUserLoginNotUnique              = errors.New("login not unique")
	ErrUserLoginUnauthorized           = errors.New("user unauthorized")
	ErrUnprocessableEntity             = errors.New("unprocessable entity")
	ErrOrderAlreadyUploadedAnotherUser = errors.New("order already uploaded another user")
	ErrOrderAlreadyUploaded            = errors.New("order already uploaded")
	// ErrInvalidGzipData           = errors.New("invalid gzip data")
	// ErrReadingRequestBody        = errors.New("error reading request body")
	ErrNoContent = errors.New("no content")
	// ErrNotImplementedServerError = errors.New("not implemented server error")
	ErrInsufficientBalance = errors.New("insufficient balance")
)

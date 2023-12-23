package entity

import "errors"

// define errors
var (
	ErrEnvVarNotFound                  = errors.New("env var not set")
	ErrMethodNotAllowed                = errors.New("method not allowed")
	ErrUnsupportedMediaType            = errors.New("unsupported media type")
	ErrInternalServerError             = errors.New("internal server error")
	ErrStatusBadRequest                = errors.New("bad request")
	ErrUserLoginNotUnique              = errors.New("login not unique")
	ErrUserLoginUnauthorized           = errors.New("user unauthorized")
	ErrUnprocessableEntity             = errors.New("unprocessable entity")
	ErrOrderAlreadyUploadedAnotherUser = errors.New("order already uploaded another user")
	ErrOrderAlreadyUploaded            = errors.New("order already uploaded")
	ErrNoContent                       = errors.New("no content")
	ErrInsufficientBalance             = errors.New("insufficient balance")
)

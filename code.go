package apierrors

import (
	"net/http"
)

var (
	ErrUnknown             = New(http.StatusInternalServerError, 0, "unknown")
	ErrInternalServerError = New(http.StatusInternalServerError, 1, "internal server error")
	ErrBadRequest          = New(http.StatusBadRequest, 2, "bad request")
	ErrNotFound            = New(http.StatusNotFound, 3, "not found")
	ErrUnauthorized        = New(http.StatusUnauthorized, 4, "unauthorized")
)

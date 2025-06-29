package routes

import (
	"errors"
	"fmt"
)

type RouteError struct {
	Method string
	Path   string
	Err    error
}

func NewRouteError(method, path string, err error) *RouteError {
	return &RouteError{
		Method: method,
		Path:   path,
		Err:    err,
	}
}

func (e *RouteError) Unwrap() error {
	return e.Err
}

func (e *RouteError) Error() string {
	return fmt.Sprintf("route error: %s %s: %v", e.Method, e.Path, e.Err)
}

var (
	ErrRouteNotFound      = errors.New("route not found")
	ErrRouteAlreadyExists = errors.New("route already exists")
	ErrRouteMalformedPath = errors.New("malformed path")
)

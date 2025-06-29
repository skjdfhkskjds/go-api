package engine

import (
	"github.com/skjdfhkskjds/go-api/engine/internal/routes"
	"github.com/skjdfhkskjds/go-api/engine/internal/types"
)

// GET registers a GET route
func (e *Engine) GET(path string, handler types.HandlerFunc) *Engine {
	e.routes.GET(path, handler)
	return e
}

// POST registers a POST route
func (e *Engine) POST(path string, handler types.HandlerFunc) *Engine {
	e.routes.POST(path, handler)
	return e
}

// PUT registers a PUT route
func (e *Engine) PUT(path string, handler types.HandlerFunc) *Engine {
	e.routes.PUT(path, handler)
	return e
}

// DELETE registers a DELETE route
func (e *Engine) DELETE(path string, handler types.HandlerFunc) *Engine {
	e.routes.DELETE(path, handler)
	return e
}

// PATCH registers a PATCH route
func (e *Engine) PATCH(path string, handler types.HandlerFunc) *Engine {
	e.routes.PATCH(path, handler)
	return e
}

// Group creates a route group with the specified prefix
func (e *Engine) Group(prefix string) *routes.RouteNode {
	group, _ := e.routes.Group(prefix)
	return group
}

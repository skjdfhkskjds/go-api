package routes

import (
	"net/http"
	"strings"

	"github.com/skjdfhkskjds/go-api/engine/internal/types"
)

// Use adds middleware to the route group
//
// @return: the route group that the middleware was added to
func (n *RouteNode) Use(middlewares ...types.MiddlewareFunc) *RouteNode {
	n.middlewares = append(n.middlewares, middlewares...)
	return n
}

// Group creates a new route group with the specified prefix
//
// @return: the newly created route node for the group
// @return: an error if the route already exists
//
// @see: RouteNode.Route
func (n *RouteNode) Group(
	prefix string,
	middlewares ...types.MiddlewareFunc,
) (*RouteNode, error) {
	return n.Route("", prefix, nil, middlewares...)
}

// Route adds a route to the tree
//
// @return: the newly added route node
// @return: an error if the route already exists
func (n *RouteNode) Route(
	method string,
	path string,
	handler types.HandlerFunc,
	middlewares ...types.MiddlewareFunc,
) (*RouteNode, error) {
	child, err := n.addRoute(method, path, handler, middlewares...)
	if err != nil {
		return nil, NewRouteError(method, path, err)
	}
	return child, nil
}

// GET registers a GET route in the route group
//
// @return: WARNING: returns the parent route node for method chaining,
//
//	not the newly created route node
//
// @return: an error if the route already exists
//
// @see: RouteNode.Route
func (n *RouteNode) GET(
	path string,
	handler types.HandlerFunc,
	middlewares ...types.MiddlewareFunc,
) (*RouteNode, error) {
	_, err := n.Route(
		http.MethodGet,
		path,
		handler,
		middlewares...,
	)
	return n, err
}

// POST registers a POST route in the route group
//
// @return: WARNING: returns the parent route node for method chaining,
//
//	not the newly created route node
//
// @return: an error if the route already exists
//
// @see: RouteNode.Route
func (n *RouteNode) POST(
	path string,
	handler types.HandlerFunc,
	middlewares ...types.MiddlewareFunc,
) (*RouteNode, error) {
	_, err := n.Route(
		http.MethodPost,
		path,
		handler,
		middlewares...,
	)
	return n, err
}

// PUT registers a PUT route in the route group
//
// @return: WARNING: returns the parent route node for method chaining,
//
//	not the newly created route node
//
// @return: an error if the route already exists
//
// @see: RouteNode.Route
func (n *RouteNode) PUT(
	path string,
	handler types.HandlerFunc,
	middlewares ...types.MiddlewareFunc,
) (*RouteNode, error) {
	_, err := n.Route(
		http.MethodPut,
		path,
		handler,
		middlewares...,
	)
	return n, err
}

// DELETE registers a DELETE route in the route group
//
// @return: WARNING: returns the parent route node for method chaining,
//
//	not the newly created route node
//
// @return: an error if the route already exists
//
// @see: RouteNode.Route
func (n *RouteNode) DELETE(
	path string,
	handler types.HandlerFunc,
	middlewares ...types.MiddlewareFunc,
) (*RouteNode, error) {
	_, err := n.Route(
		http.MethodDelete,
		path,
		handler,
		middlewares...,
	)
	return n, err
}

// PATCH registers a PATCH route in the route group
//
// @return: WARNING: returns the parent route node for method chaining,
//
//	not the newly created route node
//
// @return: an error if the route already exists
//
// @see: RouteNode.Route
func (n *RouteNode) PATCH(
	path string,
	handler types.HandlerFunc,
	middlewares ...types.MiddlewareFunc,
) (*RouteNode, error) {
	_, err := n.Route(
		http.MethodPatch,
		path,
		handler,
		middlewares...,
	)
	return n, err
}

// OPTIONS registers an OPTIONS route in the route group
//
// @return: WARNING: returns the parent route node for method chaining,
//
//	not the newly created route node
//
// @return: an error if the route already exists
//
// @see: RouteNode.Route
func (n *RouteNode) OPTIONS(
	path string,
	handler types.HandlerFunc,
	middlewares ...types.MiddlewareFunc,
) (*RouteNode, error) {
	_, err := n.Route(
		http.MethodOptions,
		path,
		handler,
		middlewares...,
	)
	return n, err
}

// HEAD registers a HEAD route in the route group
//
// @return: WARNING: returns the parent route node for method chaining,
//
//	not the newly created route node
//
// @return: an error if the route already exists
//
// @see: RouteNode.Route
func (n *RouteNode) HEAD(
	path string,
	handler types.HandlerFunc,
	middlewares ...types.MiddlewareFunc,
) (*RouteNode, error) {
	_, err := n.Route(
		http.MethodHead,
		path,
		handler,
		middlewares...,
	)
	return n, err
}

// addRoute adds a route to the tree from the given node
//
// @return: the newly added route node
// @return: an error if the route already exists
func (n *RouteNode) addRoute(
	method string,
	path string,
	handler types.HandlerFunc,
	middlewares ...types.MiddlewareFunc,
) (*RouteNode, error) {
	// If we've consumed the entire path, this is our destination
	if path == "" || path == "/" {
		// Check if route already exists for this method
		if _, exists := n.handlers[method]; exists {
			return nil, ErrRouteAlreadyExists
		}

		// Store handler and any route-specific middleware on this node
		n.middlewares = append(n.middlewares, middlewares...)
		n.handlers[method] = handler
		return n, nil
	}

	// Remove leading slash for processing
	if path[0] == '/' {
		path = path[1:]
	}

	// Split path into segments
	segment, remaining := getPathSegment(path)

	// Determine node type
	routeType, paramName := getRouteTypeFromSegment(segment)

	// Parameter route
	if routeType == RouteTypeParam {
		if n.param == nil {
			n.param = NewRouteNode(segment, routeType, paramName, n)
		}
		return n.param.addRoute(method, remaining, handler, middlewares...)
	}

	// Wildcard route
	if routeType == RouteTypeWildcard {
		if n.wildcard == nil {
			n.wildcard = NewRouteNode(segment, routeType, paramName, n)
		}
		return n.wildcard.addRoute(method, remaining, handler, middlewares...)
	}

	// Static route - check if the segment already exists in the static children
	for _, child := range n.static {
		if child.path == segment {
			return child.addRoute(method, remaining, handler, middlewares...)
		}
	}

	// Create new child node
	child := NewRouteNode(segment, RouteTypeStatic, "", n)
	n.static = append(n.static, child)
	return child.addRoute(method, remaining, handler, middlewares...)
}

// getPathSegment splits a path into a segment and a remaining path
func getPathSegment(path string) (string, string) {
	segments := strings.SplitN(path, "/", 2)
	if len(segments) == 0 {
		return "", ""
	}

	segment := segments[0]
	remaining := ""

	if len(segments) > 1 {
		remaining = segments[1]
		if remaining != "" {
			remaining = "/" + remaining
		}
	}

	return segment, remaining
}

// getRouteTypeFromSegment determines the node type and parameter name from a
// segment of a path
func getRouteTypeFromSegment(segment string) (RouteType, string) {
	if isPathParam(segment) {
		return RouteTypeParam, segment[1 : len(segment)-1]
	} else if isWildcard(segment) {
		return RouteTypeWildcard, segment[1:]
	}

	return RouteTypeStatic, ""
}

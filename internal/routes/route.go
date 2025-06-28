package routes

import (
	"maps"

	"github.com/skjdfhkskjds/go-api/internal/types"
)

// Route represents a registered route
//
// This structure is used to represent a matched route from the tree, and
// is the representation used externally.
type Route struct {
	Method      string
	Path        string
	Handler     types.HandlerFunc
	Middlewares []types.MiddlewareFunc
	PathParams  map[string]string
}

// RouteType represents the type of a route node
type RouteType int

const (
	RouteTypeNone   RouteType = -1 // Used for nodes along a path
	RouteTypeStatic RouteType = iota
	RouteTypeParam
	RouteTypeWildcard
)

// RouteNode represents a route in a radix tree
//
// This is the internal representation of a route, and is not meant to be
// used directly by the user.
type RouteNode struct {
	// Path segment for this route
	path string

	// Type of this route (static, param, wildcard)
	routeType RouteType

	// Parameter name (only for param and wildcard routes)
	paramName string

	// Handlers for different HTTP methods (method -> handler)
	handlers map[string]types.HandlerFunc

	// Middleware accumulated from parent nodes
	middlewares []types.MiddlewareFunc

	// Parent node
	parent *RouteNode

	// Child nodes
	static   []*RouteNode
	param    *RouteNode
	wildcard *RouteNode
}

func NewRouteNode(
	path string,
	routeType RouteType,
	paramName string,
	parent *RouteNode,
) *RouteNode {
	return &RouteNode{
		path:        path,
		routeType:   routeType,
		paramName:   paramName,
		parent:      parent,
		handlers:    make(map[string]types.HandlerFunc),
		middlewares: make([]types.MiddlewareFunc, 0),
		static:      make([]*RouteNode, 0),
		param:       nil,
		wildcard:    nil,
	}
}

// Path returns the full path of the route node
//
// @return: the full path of the route node
func (n *RouteNode) Path() string {
	if n.parent == nil {
		return "/"
	}

	parentPath := n.parent.Path()
	if parentPath == "/" {
		return "/" + n.path
	}
	return parentPath + "/" + n.path
}

// Find finds a route in the tree
//
// @return: the matched route as a Route struct, pathParams are populated
// @return: an error if the route is not found
//
// @see: RouteNode.find
func (n *RouteNode) Find(method, path string) (*Route, error) {
	var err error
	route := &Route{
		Method:     method,
		Path:       path,
		PathParams: make(map[string]string),
	}
	route, err = n.find(route, method, path)
	if err != nil {
		return nil, NewRouteError(method, path, err)
	}
	return route, nil
}

// find is a recursive helper function to find a route in the tree
//
// @return: the matched route as a Route struct, pathParams are populated
// @return: an error if the route is not found
//
// @see: RouteNode.Find
func (n *RouteNode) find(route *Route, method, path string) (*Route, error) {
	if path == "" || path == "/" {
		handler, exists := n.handlers[method]
		if !exists {
			return nil, ErrRouteNotFound
		}

		route.Method = method
		route.Handler = handler
		n.collectMiddlewares(&route.Middlewares)
		return route, nil
	}

	// Remove leading slash for processing
	if path[0] == '/' {
		path = path[1:]
	}

	segment, remaining := getPathSegment(path)

	// Check static routes first
	for _, child := range n.static {
		if child.path == segment {
			return child.find(route, method, remaining)
		}
	}

	// Check parameter routes
	if n.param != nil {
		// Save the current pathParams state for backtracking
		originalPathParams := make(map[string]string, len(route.PathParams))
		maps.Copy(route.PathParams, originalPathParams)

		// Store parameter name to value mapping
		route.PathParams[n.param.paramName] = segment
		result, err := n.param.find(route, method, remaining)
		if err == nil {
			return result, nil
		}

		// Restore pathParams state for backtracking
		route.PathParams = originalPathParams
	}

	// Check wildcard routes
	if n.wildcard != nil {
		// Wildcard captures everything remaining
		wildcardValue := segment
		if remaining != "" {
			wildcardValue += remaining
		}
		route.PathParams[n.wildcard.paramName] = wildcardValue
		return n.wildcard.find(route, method, "")
	}

	return nil, ErrRouteNotFound
}

// collectMiddlewares collects middleware from root to current node
func (n *RouteNode) collectMiddlewares(middlewares *[]types.MiddlewareFunc) {
	if n.parent != nil {
		n.parent.collectMiddlewares(middlewares)
	}
	*middlewares = append(*middlewares, n.middlewares...)
}

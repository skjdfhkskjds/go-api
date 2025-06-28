package routes

import (
	"net/http"
	"testing"

	"github.com/skjdfhkskjds/go-api/internal/types"
	"github.com/stretchr/testify/require"
)

// Test helper functions
func newTestHandler(name string) types.HandlerFunc {
	return func(c *types.Context) {
		c.String(200, name)
	}
}

func newTestMiddleware(_ string) types.MiddlewareFunc {
	return func(next types.HandlerFunc) types.HandlerFunc {
		return func(c *types.Context) {
			next(c)
		}
	}
}

func TestNewRouteNode(t *testing.T) {
	tests := []struct {
		name      string
		path      string
		routeType RouteType
		paramName string
		parent    *RouteNode
	}{
		{
			name:      "static node",
			path:      "users",
			routeType: RouteTypeStatic,
			paramName: "",
			parent:    nil,
		},
		{
			name:      "param node",
			path:      ":id",
			routeType: RouteTypeParam,
			paramName: "id",
			parent:    nil,
		},
		{
			name:      "wildcard node",
			path:      "*path",
			routeType: RouteTypeWildcard,
			paramName: "path",
			parent:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := NewRouteNode(tt.path, tt.routeType, tt.paramName, tt.parent)

			require.Equal(t, tt.path, node.path)
			require.Equal(t, tt.routeType, node.routeType)
			require.Equal(t, tt.paramName, node.paramName)
			require.Equal(t, tt.parent, node.parent)
			require.NotNil(t, node.handlers)
			require.NotNil(t, node.middlewares)
			require.NotNil(t, node.static)
		})
	}
}

func TestRouteNode_Route_StaticRoutes(t *testing.T) {
	root := NewRouteNode("", RouteTypeNone, "", nil)
	handler := newTestHandler("users")

	// Test adding a simple static route
	_, err := root.Route(http.MethodGet, "/users", handler)
	require.NoError(t, err)

	// Test finding the route
	route, err := root.Find(http.MethodGet, "/users")
	require.NoError(t, err)
	require.Equal(t, http.MethodGet, route.Method)
	require.Equal(t, "/users", route.Path)
	require.NotNil(t, route.Handler)

	// Test route not found
	_, err = root.Find(http.MethodPost, "/users")
	require.Error(t, err)
	require.EqualError(t, err, NewRouteError(http.MethodPost, "/users", ErrRouteNotFound).Error())

	// Test adding duplicate route
	_, err = root.Route(http.MethodGet, "/users", handler)
	require.Error(t, err)
	require.EqualError(t, err, NewRouteError(http.MethodGet, "/users", ErrRouteAlreadyExists).Error())
}

func TestRouteNode_Route_NestedStaticRoutes(t *testing.T) {
	root := NewRouteNode("", RouteTypeNone, "", nil)

	// Add nested routes
	node, err := root.Route(http.MethodGet, "/users/profile", newTestHandler("profile"))
	require.NoError(t, err)
	require.Equal(t, "/users/profile", node.Path())

	node, err = root.Route(http.MethodGet, "/users/settings", newTestHandler("settings"))
	require.NoError(t, err)
	require.Equal(t, "/users/settings", node.Path())

	// Test finding nested routes
	route, err := root.Find(http.MethodGet, "/users/profile")
	require.NoError(t, err)
	require.Equal(t, "/users/profile", route.Path)

	route, err = root.Find(http.MethodGet, "/users/settings")
	require.NoError(t, err)
	require.Equal(t, "/users/settings", route.Path)

	// Test duplicate nested routes
	_, err = root.Route(http.MethodGet, "/users/profile", newTestHandler("profile"))
	require.Error(t, err)
	require.EqualError(t, err, NewRouteError(http.MethodGet, "/users/profile", ErrRouteAlreadyExists).Error())
}

func TestRouteNode_Route_ParameterRoutes(t *testing.T) {
	root := NewRouteNode("", RouteTypeNone, "", nil)

	// Add parameter routes
	node, err := root.Route(http.MethodGet, "/users/:id", newTestHandler("user"))
	require.NoError(t, err)
	require.Equal(t, "/users/:id", node.Path())

	node, err = root.Route(http.MethodGet, "/users/:id/posts/:postId", newTestHandler("userPost"))
	require.NoError(t, err)
	require.Equal(t, "/users/:id/posts/:postId", node.Path())

	// Test finding parameter routes
	route, err := root.Find(http.MethodGet, "/users/123")
	require.NoError(t, err)
	require.Equal(t, 1, len(route.PathParams))
	require.Equal(t, "123", route.PathParams["id"])

	// Test nested parameter routes
	route, err = root.Find(http.MethodGet, "/users/456/posts/789")
	require.NoError(t, err)
	require.Equal(t, 2, len(route.PathParams))
	require.Equal(t, "456", route.PathParams["id"])
	require.Equal(t, "789", route.PathParams["postId"])
}

func TestRouteNode_Route_WildcardRoutes(t *testing.T) {
	root := NewRouteNode("", RouteTypeNone, "", nil)

	// Add wildcard route
	node, err := root.Route(http.MethodGet, "/files/*path", newTestHandler("files"))
	require.NoError(t, err)
	require.Equal(t, "/files/*path", node.Path())

	// Test finding wildcard routes
	tests := []struct {
		path     string
		expected string
	}{
		{"/files/documents", "documents"},
		{"/files/images/photo.jpg", "images/photo.jpg"},
		{"/files/folder/subfolder/file.txt", "folder/subfolder/file.txt"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			route, err := root.Find(http.MethodGet, tt.path)
			require.NoError(t, err)
			require.Equal(t, 1, len(route.PathParams))
			require.Equal(t, tt.expected, route.PathParams["path"])
		})
	}
}

func TestRouteNode_Route_RoutePriority(t *testing.T) {
	root := NewRouteNode("", RouteTypeNone, "", nil)

	// Add routes with different priorities
	node, err := root.Route(http.MethodGet, "/users/:id", newTestHandler("param"))
	require.NoError(t, err)
	require.Equal(t, "/users/:id", node.Path())

	node, err = root.Route(http.MethodGet, "/users/admin", newTestHandler("static"))
	require.NoError(t, err)
	require.Equal(t, "/users/admin", node.Path())

	node, err = root.Route(http.MethodGet, "/users/*path", newTestHandler("wildcard"))
	require.NoError(t, err)
	require.Equal(t, "/users/*path", node.Path())

	// Test that static routes have highest priority
	route, err := root.Find(http.MethodGet, "/users/admin")
	require.NoError(t, err)
	require.Empty(t, route.PathParams)

	// Test that param routes have higher priority than wildcard
	route, err = root.Find(http.MethodGet, "/users/123")
	require.NoError(t, err)
	require.Equal(t, 1, len(route.PathParams))
	require.Equal(t, "123", route.PathParams["id"])

	// Test that wildcard matches anything else
	route, err = root.Find(http.MethodGet, "/users/some/long/path")
	require.NoError(t, err)
	require.Equal(t, 1, len(route.PathParams))
	require.Equal(t, "some/long/path", route.PathParams["path"])
}

func TestRouteNode_Route_MultipleMethodsSamePath(t *testing.T) {
	root := NewRouteNode("", RouteTypeNone, "", nil)

	// Add multiple methods for the same path
	methods := []string{
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodDelete,
		http.MethodPatch,
	}

	for _, method := range methods {
		node, err := root.Route(method, "/users", newTestHandler(method))
		require.NoError(t, err)
		require.Equal(t, "/users", node.Path())
	}

	// Test that all methods can be found
	for _, method := range methods {
		route, err := root.Find(method, "/users")
		require.NoError(t, err)
		require.Equal(t, method, route.Method)
	}

	// Test that non-existent method returns error
	_, err := root.Find(http.MethodOptions, "/users")
	require.Error(t, err)
	require.EqualError(t, err, NewRouteError(http.MethodOptions, "/users", ErrRouteNotFound).Error())
}

func TestRouteNode_Use_Middleware(t *testing.T) {
	root := NewRouteNode("", RouteTypeNone, "", nil)

	// Add middleware to root
	middleware1 := newTestMiddleware("auth")
	middleware2 := newTestMiddleware("logging")
	root.Use(middleware1, middleware2)

	// Add a route
	node, err := root.Route(http.MethodGet, "/users", newTestHandler("users"))
	require.NoError(t, err)
	require.Equal(t, "/users", node.Path())

	// Find the route and check middleware
	route, err := root.Find(http.MethodGet, "/users")
	require.NoError(t, err)
	require.Equal(t, 2, len(route.Middlewares))

	// Add middleware to the route
	node.Use(newTestMiddleware("rateLimit"))

	// Find the route and check middleware again
	route, err = root.Find(http.MethodGet, "/users")
	require.NoError(t, err)
	require.Equal(t, 3, len(route.Middlewares))
}

func TestRouteNode_HTTPMethods(t *testing.T) {
	root := NewRouteNode("", RouteTypeNone, "", nil)
	tests := []struct {
		name   string
		method string
		fn     func(path string, handler types.HandlerFunc, middleware ...types.MiddlewareFunc) (*RouteNode, error)
	}{
		{"GET", http.MethodGet, root.GET},
		{"POST", http.MethodPost, root.POST},
		{"PUT", http.MethodPut, root.PUT},
		{"DELETE", http.MethodDelete, root.DELETE},
		{"PATCH", http.MethodPatch, root.PATCH},
		{"OPTIONS", http.MethodOptions, root.OPTIONS},
		{"HEAD", http.MethodHead, root.HEAD},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := "/test-" + tt.name
			handler := newTestHandler(tt.name)
			middleware := newTestMiddleware(tt.name)

			_, err := tt.fn(path, handler, middleware)
			require.NoError(t, err)

			// Verify the route was added correctly
			route, err := root.Find(tt.method, path)
			require.NoError(t, err)
			require.Equal(t, tt.method, route.Method)
			require.Equal(t, path, route.Path)
		})
	}
}

func TestRouteNode_EdgeCases(t *testing.T) {
	root := NewRouteNode("", RouteTypeNone, "", nil)

	t.Run("root path", func(t *testing.T) {
		node, err := root.Route(http.MethodGet, "/", newTestHandler("root"))
		require.NoError(t, err)
		require.Equal(t, "/", node.Path())

		route, err := root.Find(http.MethodGet, "/")
		require.NoError(t, err)
		require.Equal(t, "/", route.Path)
	})

	t.Run("empty path equivalent to root", func(t *testing.T) {
		// Empty path should be treated the same as root path "/"
		// So we should be able to find the root route with empty path
		route, err := root.Find(http.MethodGet, "")
		require.NoError(t, err)
		require.Equal(t, "", route.Path)

		// And trying to add a route with empty path should fail since root exists
		_, err = root.Route(http.MethodGet, "", newTestHandler("empty"))
		require.Error(t, err)
		require.EqualError(t, err, NewRouteError(http.MethodGet, "", ErrRouteAlreadyExists).Error())
	})

	t.Run("path without leading slash", func(t *testing.T) {
		node, err := root.Route(http.MethodGet, "users", newTestHandler("no-slash"))
		require.NoError(t, err)
		require.Equal(t, "/users", node.Path())

		route, err := root.Find(http.MethodGet, "users")
		require.NoError(t, err)
		require.Equal(t, "users", route.Path)
	})
}

func TestGetPathSegment(t *testing.T) {
	tests := []struct {
		path              string
		expectedSegment   string
		expectedRemaining string
	}{
		{"users", "users", ""},
		{"users/profile", "users", "/profile"},
		{"users/profile/settings", "users", "/profile/settings"},
		{"", "", ""},
		{"/", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			segment, remaining := getPathSegment(tt.path)
			require.Equal(t, tt.expectedSegment, segment)
			require.Equal(t, tt.expectedRemaining, remaining)
		})
	}
}

func TestGetRouteTypeFromSegment(t *testing.T) {
	tests := []struct {
		segment           string
		expectedType      RouteType
		expectedParamName string
	}{
		{"users", RouteTypeStatic, ""},
		{":id", RouteTypeParam, "id"},
		{":userId", RouteTypeParam, "userId"},
		{"*path", RouteTypeWildcard, "path"},
		{"*filepath", RouteTypeWildcard, "filepath"},
	}

	for _, tt := range tests {
		t.Run(tt.segment, func(t *testing.T) {
			routeType, paramName := getRouteTypeFromSegment(tt.segment)
			require.Equal(t, tt.expectedType, routeType)
			require.Equal(t, tt.expectedParamName, paramName)
		})
	}
}

// Benchmark tests
func BenchmarkRouteNode_AddStaticRoute(b *testing.B) {
	NewRouteNode("", RouteTypeNone, "", nil)
	handler := newTestHandler("benchmark")

	b.ResetTimer()
	for b.Loop() {
		// Create a new root for each iteration to avoid duplicate route errors
		testRoot := NewRouteNode("", RouteTypeNone, "", nil)
		testRoot.Route(http.MethodGet, "/users", handler)
	}
}

func BenchmarkRouteNode_FindStaticRoute(b *testing.B) {
	root := NewRouteNode("", RouteTypeNone, "", nil)
	handler := newTestHandler("benchmark")
	root.Route(http.MethodGet, "/users", handler)

	b.ResetTimer()
	for b.Loop() {
		root.Find(http.MethodGet, "/users")
	}
}

func BenchmarkRouteNode_FindParameterRoute(b *testing.B) {
	root := NewRouteNode("", RouteTypeNone, "", nil)
	handler := newTestHandler("benchmark")
	root.Route(http.MethodGet, "/users/:id", handler)

	b.ResetTimer()
	for b.Loop() {
		root.Find(http.MethodGet, "/users/123")
	}
}

func BenchmarkRouteNode_FindWildcardRoute(b *testing.B) {
	root := NewRouteNode("", RouteTypeNone, "", nil)
	handler := newTestHandler("benchmark")
	root.Route(http.MethodGet, "/files/*path", handler)

	b.ResetTimer()
	for b.Loop() {
		root.Find(http.MethodGet, "/files/documents/file.txt")
	}
}

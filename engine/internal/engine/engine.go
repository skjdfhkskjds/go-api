package engine

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/skjdfhkskjds/go-api/engine/internal/routes"
	"github.com/skjdfhkskjds/go-api/engine/internal/types"
)

// Engine is the core framework engine
type Engine struct {
	config *Config
	routes *routes.RouteNode
	server *http.Server
}

// New creates a new Engine instance with the provided configuration
func New(config *Config) *Engine {
	if config == nil {
		config = DefaultConfig()
	}

	engine := &Engine{
		config: config,
		routes: routes.NewRouteNode("", routes.RouteTypeNone, "", nil),
	}

	return engine
}

// ServeHTTP implements http.Handler interface
func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Convert net/http request to our Context type
	ctx := &types.Context{
		Request:    r,
		Writer:     w,
		PathParams: make(map[string]string),
	}

	// Find matching route using RouteNode
	route, err := e.routes.Find(r.Method, r.URL.Path)
	if err != nil {
		ctx.ErrorString(http.StatusNotFound, "Not Found")
		return
	}

	// Set path parameters from route matching
	ctx.PathParams = route.PathParams

	// Execute handler (middleware support can be added later)
	route.Handler(ctx)
}

// Run starts the HTTP server
func (e *Engine) Run(addr ...string) error {
	address := e.resolveAddress(addr)

	e.server = &http.Server{
		Addr:         address,
		Handler:      e,
		ReadTimeout:  time.Duration(e.config.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(e.config.Server.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(e.config.Server.IdleTimeout) * time.Second,
	}

	log.Printf("Server starting on %s", address)
	return e.server.ListenAndServe()
}

// resolveAddress resolves the server address
func (e *Engine) resolveAddress(addr []string) string {
	switch len(addr) {
	case 0:
		if e.config.Server.Port != 0 {
			return fmt.Sprintf(":%d", e.config.Server.Port)
		}
		return ":8080"
	case 1:
		return addr[0]
	default:
		panic("too many parameters")
	}
}

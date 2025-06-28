package types

// MiddlewareFunc defines the signature for middleware functions
type MiddlewareFunc func(HandlerFunc) HandlerFunc

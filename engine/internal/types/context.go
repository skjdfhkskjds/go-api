package types

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// Context provides request context and response utilities
type Context struct {
	context.Context

	Request    *http.Request
	Writer     http.ResponseWriter
	PathParams map[string]string
}

// JSON sends a JSON response
func (c *Context) JSON(status int, data any) {
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(status)

	if err := json.NewEncoder(c.Writer).Encode(data); err != nil {
		c.Error(http.StatusInternalServerError, err)
	}
}

// String sends a string response
func (c *Context) String(status int, data string) {
	c.Writer.Header().Set("Content-Type", "text/plain")
	c.Writer.WriteHeader(status)
	c.Writer.Write([]byte(data))
}

// HTML sends an HTML response
func (c *Context) HTML(status int, html string) {
	c.Writer.Header().Set("Content-Type", "text/html")
	c.Writer.WriteHeader(status)
	c.Writer.Write([]byte(html))
}

// Data sends raw data response
func (c *Context) Data(status int, contentType string, data []byte) {
	c.Writer.Header().Set("Content-Type", contentType)
	c.Writer.WriteHeader(status)
	c.Writer.Write(data)
}

// Status sets the HTTP status code
func (c *Context) Status(status int) {
	c.Writer.WriteHeader(status)
}

// Header sets a response header
func (c *Context) Header(key, value string) {
	c.Writer.Header().Set(key, value)
}

// GetHeader gets a request header
func (c *Context) GetHeader(key string) string {
	return c.Request.Header.Get(key)
}

// GetParam gets a path parameter by name
func (c *Context) GetParam(name string) string {
	return c.PathParams[name]
}

// GetParamInt gets a path parameter as integer
func (c *Context) GetParamInt(name string) (int, error) {
	param := c.PathParams[name]
	if param == "" {
		return 0, fmt.Errorf("parameter %s not found", name)
	}
	return strconv.Atoi(param)
}

// GetParamInt64 gets a path parameter as int64
func (c *Context) GetParamInt64(name string) (int64, error) {
	param := c.PathParams[name]
	if param == "" {
		return 0, fmt.Errorf("parameter %s not found", name)
	}
	return strconv.ParseInt(param, 10, 64)
}

// GetQuery gets a query parameter
func (c *Context) GetQuery(name string) string {
	return c.Request.URL.Query().Get(name)
}

// GetQueryDefault gets a query parameter with default value
func (c *Context) GetQueryDefault(name, defaultValue string) string {
	value := c.Request.URL.Query().Get(name)
	if value == "" {
		return defaultValue
	}
	return value
}

// GetQueryInt gets a query parameter as integer
func (c *Context) GetQueryInt(name string) (int, error) {
	param := c.Request.URL.Query().Get(name)
	if param == "" {
		return 0, fmt.Errorf("query parameter %s not found", name)
	}
	return strconv.Atoi(param)
}

// GetQueryIntDefault gets a query parameter as integer with default value
func (c *Context) GetQueryIntDefault(name string, defaultValue int) int {
	param := c.Request.URL.Query().Get(name)
	if param == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(param)
	if err != nil {
		return defaultValue
	}
	return value
}

// GetCookie gets a cookie value
func (c *Context) GetCookie(name string) (string, error) {
	cookie, err := c.Request.Cookie(name)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

// SetCookie sets a cookie
func (c *Context) SetCookie(name, value string, maxAge int, path, domain string, secure, httpOnly bool) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		MaxAge:   maxAge,
		Path:     path,
		Domain:   domain,
		Secure:   secure,
		HttpOnly: httpOnly,
	}
	http.SetCookie(c.Writer, cookie)
}

// BindJSON binds JSON request body to a struct
func (c *Context) BindJSON(obj any) error {
	return json.NewDecoder(c.Request.Body).Decode(obj)
}

// GetUserAgent gets the User-Agent header
func (c *Context) GetUserAgent() string {
	return c.Request.Header.Get("User-Agent")
}

// GetClientIP gets the client IP address
func (c *Context) GetClientIP() string {
	// Check for X-Forwarded-For header first
	if xff := c.Request.Header.Get("X-Forwarded-For"); xff != "" {
		return xff
	}

	// Check for X-Real-IP header
	if xri := c.Request.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	// Fall back to RemoteAddr
	return c.Request.RemoteAddr
}

// IsAjax checks if the request is an AJAX request
func (c *Context) IsAjax() bool {
	return c.Request.Header.Get("X-Requested-With") == "XMLHttpRequest"
}

// GetContentType gets the Content-Type header
func (c *Context) GetContentType() string {
	return c.Request.Header.Get("Content-Type")
}

// Redirect redirects to the given URL
func (c *Context) Redirect(status int, url string) {
	http.Redirect(c.Writer, c.Request, url, status)
}

// Error sends an error response
func (c *Context) Error(status int, err error) {
	c.JSON(status, map[string]any{
		"error":   http.StatusText(status),
		"message": err.Error(),
	})
}

// ErrorString sends an error response with string message
func (c *Context) ErrorString(status int, message string) {
	c.JSON(status, map[string]any{
		"error":   http.StatusText(status),
		"message": message,
	})
}

// Success sends a success response
func (c *Context) Success(data any) {
	c.JSON(200, map[string]any{
		"success": true,
		"data":    data,
	})
}

// SuccessWithMessage sends a success response with message
func (c *Context) SuccessWithMessage(message string, data any) {
	c.JSON(200, map[string]any{
		"success": true,
		"message": message,
		"data":    data,
	})
}

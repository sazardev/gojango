package gojango

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

// JSON sends a JSON response
func (c *Context) JSON(data interface{}) error {
	c.Response.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(c.Response).Encode(data)
}

// ErrorJSON sends an error JSON response
func (c *Context) ErrorJSON(status int, message string, err error) error {
	c.Response.Header().Set("Content-Type", "application/json")
	c.Response.WriteHeader(status)
	
	errorResponse := map[string]interface{}{
		"error":   message,
		"status":  status,
	}
	
	if err != nil {
		errorResponse["details"] = err.Error()
	}
	
	return json.NewEncoder(c.Response).Encode(errorResponse)
}

// BindJSON binds request body to a struct
func (c *Context) BindJSON(v interface{}) error {
	if c.Request.Header.Get("Content-Type") != "application/json" {
		return fmt.Errorf("content-type must be application/json")
	}
	
	decoder := json.NewDecoder(c.Request.Body)
	defer c.Request.Body.Close()
	
	return decoder.Decode(v)
}

// Param gets a URL parameter by name
func (c *Context) Param(name string) string {
	// First check if it's already parsed
	if val, exists := c.Params[name]; exists {
		return val
	}
	
	// Extract from URL path (simple implementation)
	// This would be set by the router when matching routes
	return c.Request.URL.Query().Get(name)
}

// ParamInt gets a URL parameter as integer
func (c *Context) ParamInt(name string) (int, error) {
	val := c.Param(name)
	if val == "" {
		return 0, fmt.Errorf("parameter %s not found", name)
	}
	
	return strconv.Atoi(val)
}

// Query gets a query parameter
func (c *Context) Query(name string) string {
	return c.Request.URL.Query().Get(name)
}

// QueryInt gets a query parameter as integer
func (c *Context) QueryInt(name string) (int, error) {
	val := c.Query(name)
	if val == "" {
		return 0, fmt.Errorf("query parameter %s not found", name)
	}
	
	return strconv.Atoi(val)
}

// FormValue gets a form value
func (c *Context) FormValue(name string) string {
	return c.Request.FormValue(name)
}

// String sends a plain text response
func (c *Context) String(data string) error {
	c.Response.Header().Set("Content-Type", "text/plain")
	_, err := c.Response.Write([]byte(data))
	return err
}

// HTML sends an HTML response
func (c *Context) HTML(html string) error {
	c.Response.Header().Set("Content-Type", "text/html")
	_, err := c.Response.Write([]byte(html))
	return err
}

// Render renders a template with data
func (c *Context) Render(templateName string, data interface{}) error {
	if c.app.templates == nil {
		return fmt.Errorf("template engine not configured")
	}
	
	return c.app.templates.Render(c.Response, templateName, data)
}

// Status sets the HTTP status code
func (c *Context) Status(code int) {
	c.Response.WriteHeader(code)
}

// Header sets a response header
func (c *Context) Header(key, value string) {
	c.Response.Header().Set(key, value)
}

// GetHeader gets a request header
func (c *Context) GetHeader(key string) string {
	return c.Request.Header.Get(key)
}

// Body gets the request body as bytes
func (c *Context) Body() ([]byte, error) {
	defer c.Request.Body.Close()
	return io.ReadAll(c.Request.Body)
}

// Method gets the HTTP method
func (c *Context) Method() string {
	return c.Request.Method
}

// Path gets the request path
func (c *Context) Path() string {
	return c.Request.URL.Path
}

// IsAjax checks if the request is an AJAX request
func (c *Context) IsAjax() bool {
	return strings.ToLower(c.GetHeader("X-Requested-With")) == "xmlhttprequest"
}

// IsJSON checks if the request contains JSON
func (c *Context) IsJSON() bool {
	return strings.Contains(strings.ToLower(c.GetHeader("Content-Type")), "application/json")
}

// ClientIP gets the client IP address
func (c *Context) ClientIP() string {
	// Check for forwarded headers first
	if ip := c.GetHeader("X-Forwarded-For"); ip != "" {
		return strings.Split(ip, ",")[0]
	}
	
	if ip := c.GetHeader("X-Real-IP"); ip != "" {
		return ip
	}
	
	return c.Request.RemoteAddr
}

// Set stores a value in the context (for middleware communication)
func (c *Context) Set(key string, value interface{}) {
	if c.Params == nil {
		c.Params = make(map[string]string)
	}
	// Using string conversion for simplicity - in production you'd want a proper context store
	c.Params["__context_"+key] = fmt.Sprintf("%v", value)
}

// Get retrieves a value from the context
func (c *Context) Get(key string) (interface{}, bool) {
	val, exists := c.Params["__context_"+key]
	return val, exists
}

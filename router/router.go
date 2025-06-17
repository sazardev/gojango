package router

import (
	"net/http"
	"regexp"
	"strings"
)

// Router handles HTTP routing with parameter extraction
type Router struct {
	routes map[string][]*Route
	mux    *http.ServeMux
}

// Route represents a single route
type Route struct {
	Pattern string
	Handler http.HandlerFunc
	Regex   *regexp.Regexp
	Params  []string
}

// New creates a new router
func New() *Router {
	return &Router{
		routes: make(map[string][]*Route),
		mux:    http.NewServeMux(),
	}
}

// GET registers a GET route
func (r *Router) GET(pattern string, handler http.HandlerFunc) {
	r.addRoute("GET", pattern, handler)
}

// POST registers a POST route
func (r *Router) POST(pattern string, handler http.HandlerFunc) {
	r.addRoute("POST", pattern, handler)
}

// PUT registers a PUT route
func (r *Router) PUT(pattern string, handler http.HandlerFunc) {
	r.addRoute("PUT", pattern, handler)
}

// DELETE registers a DELETE route
func (r *Router) DELETE(pattern string, handler http.HandlerFunc) {
	r.addRoute("DELETE", pattern, handler)
}

// PATCH registers a PATCH route
func (r *Router) PATCH(pattern string, handler http.HandlerFunc) {
	r.addRoute("PATCH", pattern, handler)
}

// addRoute adds a route to the router
func (r *Router) addRoute(method, pattern string, handler http.HandlerFunc) {
	route := &Route{
		Pattern: pattern,
		Handler: handler,
	}
	
	// Convert pattern to regex for parameter extraction
	regexPattern, params := r.patternToRegex(pattern)
	route.Regex = regexp.MustCompile("^" + regexPattern + "$")
	route.Params = params
	
	if r.routes[method] == nil {
		r.routes[method] = make([]*Route, 0)
	}
	
	r.routes[method] = append(r.routes[method], route)
}

// patternToRegex converts a route pattern to regex
func (r *Router) patternToRegex(pattern string) (string, []string) {
	var params []string
	
	// Replace :param with ([^/]+) and collect parameter names
	paramRegex := regexp.MustCompile(`:([a-zA-Z_][a-zA-Z0-9_]*)`)
	
	regexPattern := paramRegex.ReplaceAllStringFunc(pattern, func(match string) string {
		paramName := match[1:] // Remove the :
		params = append(params, paramName)
		return `([^/]+)` // Match any character except /
	})
	
	// Escape other regex special characters
	regexPattern = strings.ReplaceAll(regexPattern, ".", `\.`)
	regexPattern = strings.ReplaceAll(regexPattern, "*", `.*`)
	
	return regexPattern, params
}

// ServeHTTP implements http.Handler
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	method := req.Method
	path := req.URL.Path
	
	// Find matching route
	routes, exists := r.routes[method]
	if !exists {
		http.NotFound(w, req)
		return
	}
	
	for _, route := range routes {
		matches := route.Regex.FindStringSubmatch(path)
		if matches != nil {
			// Extract parameters
			params := make(map[string]string)
			for i, paramName := range route.Params {
				if i+1 < len(matches) {
					params[paramName] = matches[i+1]
				}
			}
			
			// Store parameters in request context (simplified approach)
			// In a real implementation, you'd use context.Context
			req.Header.Set("X-Route-Params", r.encodeParams(params))
			
			route.Handler(w, req)
			return
		}
	}
	
	http.NotFound(w, req)
}

// encodeParams encodes parameters for header storage (simplified)
func (r *Router) encodeParams(params map[string]string) string {
	var parts []string
	for k, v := range params {
		parts = append(parts, k+"="+v)
	}
	return strings.Join(parts, "&")
}

// DecodeParams decodes parameters from header (helper function)
func DecodeParams(encoded string) map[string]string {
	params := make(map[string]string)
	if encoded == "" {
		return params
	}
	
	parts := strings.Split(encoded, "&")
	for _, part := range parts {
		kv := strings.SplitN(part, "=", 2)
		if len(kv) == 2 {
			params[kv[0]] = kv[1]
		}
	}
	
	return params
}

// Package gojango provides a Django-inspired web framework for Go
// with batteries included: ORM, routing, templates, and automatic CRUD.
package gojango

import (
	"fmt"
	"log"
	"net/http"
	"reflect"

	"gojango/config"
	"gojango/database"
	"gojango/models"
	"gojango/router"
	"gojango/templates"
)

// App represents the main application instance
type App struct {
	router     *router.Router
	db         *database.DB
	config     *config.Config
	templates  *templates.Engine
	middleware []Middleware
}

// Context wraps HTTP request/response with useful methods
type Context struct {
	Request  *http.Request
	Response http.ResponseWriter
	Params   map[string]string
	app      *App
}

// Middleware defines the middleware function signature
type Middleware func(*Context) error

// HandlerFunc defines the handler function signature
type HandlerFunc func(*Context) error

// New creates a new GoJango application with sensible defaults
func New(opts ...Option) *App {
	app := &App{
		router:    router.New(),
		config:    config.New(),
		templates: templates.New(),
	}

	// Apply options
	for _, opt := range opts {
		opt(app)
	}

	// Initialize database if configured
	if app.config.DatabaseURL != "" {
		var err error
		app.db, err = database.Connect(app.config.DatabaseURL)
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}
	}

	return app
}

// Option defines configuration options for the app
type Option func(*App)

// WithConfig sets custom configuration
func WithConfig(cfg *config.Config) Option {
	return func(app *App) {
		app.config = cfg
	}
}

// WithDatabase sets custom database connection
func WithDatabase(db *database.DB) Option {
	return func(app *App) {
		app.db = db
	}
}

// GET registers a GET route
func (app *App) GET(path string, handler HandlerFunc) {
	app.router.GET(path, app.wrapHandler(handler))
}

// POST registers a POST route
func (app *App) POST(path string, handler HandlerFunc) {
	app.router.POST(path, app.wrapHandler(handler))
}

// PUT registers a PUT route
func (app *App) PUT(path string, handler HandlerFunc) {
	app.router.PUT(path, app.wrapHandler(handler))
}

// DELETE registers a DELETE route
func (app *App) DELETE(path string, handler HandlerFunc) {
	app.router.DELETE(path, app.wrapHandler(handler))
}

// Use adds middleware to the application
func (app *App) Use(middleware Middleware) {
	app.middleware = append(app.middleware, middleware)
}

// AutoMigrate automatically creates/updates database tables for models
func (app *App) AutoMigrate(models ...interface{}) error {
	if app.db == nil {
		return fmt.Errorf("database not configured")
	}
	
	for _, model := range models {
		if err := app.db.AutoMigrate(model); err != nil {
			return fmt.Errorf("failed to migrate %T: %v", model, err)
		}
	}
	
	return nil
}

// RegisterCRUD automatically creates CRUD endpoints for a model
func (app *App) RegisterCRUD(basePath string, model interface{}) {
	modelType := reflect.TypeOf(model)
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}
	
	// List endpoint
	app.GET(basePath, func(c *Context) error {
		results, err := app.db.FindAll(model)
		if err != nil {
			return c.ErrorJSON(500, "Database error", err)
		}
		return c.JSON(results)
	})
	
	// Create endpoint
	app.POST(basePath, func(c *Context) error {
		newModel := reflect.New(modelType).Interface()
		if err := c.BindJSON(newModel); err != nil {
			return c.ErrorJSON(400, "Invalid JSON", err)
		}
		
		if err := app.db.Create(newModel); err != nil {
			return c.ErrorJSON(500, "Database error", err)
		}
		
		return c.JSON(newModel)
	})
	
	// Get by ID endpoint
	app.GET(basePath+"/:id", func(c *Context) error {
		id := c.Param("id")
		result := reflect.New(modelType).Interface()
		
		if err := app.db.FindByID(result, id); err != nil {
			return c.ErrorJSON(404, "Not found", err)
		}
		
		return c.JSON(result)
	})
	
	// Update endpoint
	app.PUT(basePath+"/:id", func(c *Context) error {
		id := c.Param("id")
		updateModel := reflect.New(modelType).Interface()
		
		if err := c.BindJSON(updateModel); err != nil {
			return c.ErrorJSON(400, "Invalid JSON", err)
		}
		
		if err := app.db.Update(updateModel, id); err != nil {
			return c.ErrorJSON(500, "Database error", err)
		}
		
		return c.JSON(updateModel)
	})
	
	// Delete endpoint
	app.DELETE(basePath+"/:id", func(c *Context) error {
		id := c.Param("id")
		deleteModel := reflect.New(modelType).Interface()
		
		if err := app.db.Delete(deleteModel, id); err != nil {
			return c.ErrorJSON(500, "Database error", err)
		}
		
		return c.JSON(map[string]string{"message": "Deleted successfully"})
	})
}

// Run starts the HTTP server
func (app *App) Run(addr string) error {
	if addr == "" {
		addr = app.config.GetString("server.port", ":8000")
	}
	
	log.Printf("🚀 GoJango server starting on %s", addr)
	return http.ListenAndServe(addr, app.router)
}

// wrapHandler wraps a HandlerFunc to work with the router
func (app *App) wrapHandler(handler HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := &Context{
			Request:  r,
			Response: w,
			Params:   make(map[string]string),
			app:      app,
		}
		
		// Extract route parameters from header (set by router)
		if paramHeader := r.Header.Get("X-Route-Params"); paramHeader != "" {
			r.Header.Del("X-Route-Params") // Clean up
			for k, v := range router.DecodeParams(paramHeader) {
				ctx.Params[k] = v
			}
		}
		
		// Execute middleware chain
		for _, middleware := range app.middleware {
			if err := middleware(ctx); err != nil {
				ctx.ErrorJSON(500, "Middleware error", err)
				return
			}
		}
		
		// Execute handler
		if err := handler(ctx); err != nil {
			ctx.ErrorJSON(500, "Handler error", err)
		}
	}
}

// RouteGroup allows grouping routes with common middleware
type RouteGroup struct {
	app        *App
	prefix     string
	middleware []Middleware
}

// Group creates a new route group with a prefix
func (app *App) Group(prefix string) *RouteGroup {
	return &RouteGroup{
		app:        app,
		prefix:     prefix,
		middleware: make([]Middleware, 0),
	}
}

// Use adds middleware to the route group
func (rg *RouteGroup) Use(middleware Middleware) {
	rg.middleware = append(rg.middleware, middleware)
}

// GET registers a GET route in the group
func (rg *RouteGroup) GET(path string, handler HandlerFunc) {
	fullPath := rg.prefix + path
	wrappedHandler := rg.wrapWithGroupMiddleware(handler)
	rg.app.router.GET(fullPath, rg.app.wrapHandler(wrappedHandler))
}

// POST registers a POST route in the group
func (rg *RouteGroup) POST(path string, handler HandlerFunc) {
	fullPath := rg.prefix + path
	wrappedHandler := rg.wrapWithGroupMiddleware(handler)
	rg.app.router.POST(fullPath, rg.app.wrapHandler(wrappedHandler))
}

// PUT registers a PUT route in the group
func (rg *RouteGroup) PUT(path string, handler HandlerFunc) {
	fullPath := rg.prefix + path
	wrappedHandler := rg.wrapWithGroupMiddleware(handler)
	rg.app.router.PUT(fullPath, rg.app.wrapHandler(wrappedHandler))
}

// DELETE registers a DELETE route in the group
func (rg *RouteGroup) DELETE(path string, handler HandlerFunc) {
	fullPath := rg.prefix + path
	wrappedHandler := rg.wrapWithGroupMiddleware(handler)
	rg.app.router.DELETE(fullPath, rg.app.wrapHandler(wrappedHandler))
}

// wrapWithGroupMiddleware wraps handler with group-specific middleware
func (rg *RouteGroup) wrapWithGroupMiddleware(handler HandlerFunc) HandlerFunc {
	return func(c *Context) error {
		// Execute group middleware first
		for _, middleware := range rg.middleware {
			if err := middleware(c); err != nil {
				return err
			}
		}
		
		// Then execute the handler
		return handler(c)
	}
}

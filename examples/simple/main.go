package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// Simple minimal example to test the concept
type SimpleApp struct {
	mux *http.ServeMux
}

type SimpleContext struct {
	w http.ResponseWriter
	r *http.Request
}

func NewSimpleApp() *SimpleApp {
	return &SimpleApp{
		mux: http.NewServeMux(),
	}
}

func (app *SimpleApp) GET(pattern string, handler func(*SimpleContext)) {
	app.mux.HandleFunc("GET "+pattern, func(w http.ResponseWriter, r *http.Request) {
		ctx := &SimpleContext{w: w, r: r}
		handler(ctx)
	})
}

func (app *SimpleApp) POST(pattern string, handler func(*SimpleContext)) {
	app.mux.HandleFunc("POST "+pattern, func(w http.ResponseWriter, r *http.Request) {
		ctx := &SimpleContext{w: w, r: r}
		handler(ctx)
	})
}

func (c *SimpleContext) JSON(data interface{}) {
	c.w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(c.w).Encode(data)
}

func (c *SimpleContext) String(text string) {
	c.w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(c.w, text)
}

func (app *SimpleApp) Run(addr string) error {
	log.Printf("🚀 Simple GoJango demo running on %s", addr)
	return http.ListenAndServe(addr, app.mux)
}

func main() {
	app := NewSimpleApp()

	// Basic routes
	app.GET("/", func(c *SimpleContext) {
		c.JSON(map[string]interface{}{
			"message":   "Hello GoJango! 🐍🐹",
			"timestamp": time.Now().Format(time.RFC3339),
			"status":    "working",
		})
	})

	app.GET("/health", func(c *SimpleContext) {
		c.JSON(map[string]string{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	app.GET("/hello/{name}", func(c *SimpleContext) {
		name := c.r.PathValue("name")
		c.JSON(map[string]string{
			"message": fmt.Sprintf("Hello %s from GoJango!", name),
		})
	})

	log.Println("📝 Available endpoints:")
	log.Println("   GET /")
	log.Println("   GET /health")  
	log.Println("   GET /hello/{name}")
	log.Println("")
	log.Println("🎯 Try: curl http://localhost:8000/")
	log.Println("🎯 Try: curl http://localhost:8000/hello/John")

	if err := app.Run(":8000"); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

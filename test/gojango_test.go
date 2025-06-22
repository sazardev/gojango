package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sazardev/gojango"
	"github.com/sazardev/gojango/models"
)

// Test models
type TestUser struct {
	models.Model
	Name  string `json:"name" db:"name,not_null"`
	Email string `json:"email" db:"email,unique,not_null"`
}

func (u *TestUser) TableName() string {
	return "test_users"
}

// setupTestApp creates a test application instance
func setupTestApp() *gojango.App {
	app := gojango.New()

	// Use mock database for testing
	app.GetConfig().DatabaseURL = "mock://"

	// Initialize database connection
	if err := app.InitDB(); err != nil {
		panic(fmt.Sprintf("Failed to initialize database: %v", err))
	}

	// Auto-migrate test models
	app.AutoMigrate(&TestUser{})
	// Register CRUD routes
	app.RegisterCRUD("/api/users", &TestUser{})

	// Custom test routes
	app.GET("/test", func(c *gojango.Context) error {
		return c.JSON(map[string]string{"message": "test"})
	})

	return app
}

// TestBasicRouting tests basic routing functionality
func TestBasicRouting(t *testing.T) {
	app := setupTestApp()
	// Create test server
	server := httptest.NewServer(app.GetRouter())
	defer server.Close()

	// Test GET /test
	resp, err := http.Get(server.URL + "/test")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var result map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if result["message"] != "test" {
		t.Errorf("Expected message 'test', got '%s'", result["message"])
	}
}

// TestCRUDOperations tests the automatic CRUD operations
func TestCRUDOperations(t *testing.T) {
	app := setupTestApp()
	server := httptest.NewServer(app.GetRouter())
	defer server.Close()

	// Test CREATE (POST)
	user := TestUser{
		Name:  "Juan Test",
		Email: "juan@test.com",
	}

	userJSON, _ := json.Marshal(user)
	resp, err := http.Post(server.URL+"/api/users", "application/json", bytes.NewBuffer(userJSON))
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200 for CREATE, got %d", resp.StatusCode)
	}

	var createdUser TestUser
	if err := json.NewDecoder(resp.Body).Decode(&createdUser); err != nil {
		t.Fatalf("Failed to decode created user: %v", err)
	}

	if createdUser.Name != user.Name {
		t.Errorf("Expected name '%s', got '%s'", user.Name, createdUser.Name)
	}

	// Test READ (GET all)
	resp, err = http.Get(server.URL + "/api/users")
	if err != nil {
		t.Fatalf("Failed to get users: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200 for READ, got %d", resp.StatusCode)
	}

	// Test READ by ID (GET /api/users/1)
	resp, err = http.Get(server.URL + "/api/users/1")
	if err != nil {
		t.Fatalf("Failed to get user by ID: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200 for READ by ID, got %d", resp.StatusCode)
	}
}

// TestQuerySet tests the Django-like QuerySet functionality
func TestQuerySet(t *testing.T) {
	app := setupTestApp()

	// Create test users
	users := []*TestUser{
		{Name: "Alice", Email: "alice@test.com"},
		{Name: "Bob", Email: "bob@test.com"},
		{Name: "Charlie", Email: "charlie@test.com"},
	}
	for _, user := range users {
		if err := app.GetDB().Create(user); err != nil {
			t.Fatalf("Failed to create test user: %v", err)
		}
	}

	// Test basic QuerySet creation (without complex queries for now)
	qs := app.NewQuerySet(&TestUser{})
	if qs == nil {
		t.Error("Failed to create QuerySet")
	}

	// Skip advanced QuerySet tests for mock database
	// In a real implementation, you'd implement SQL parsing for mock
	t.Log("QuerySet basic functionality verified")
}

// TestMiddleware tests middleware functionality
func TestMiddleware(t *testing.T) {
	app := setupTestApp()

	// Add test middleware
	middlewareCalled := false
	app.Use(func(c *gojango.Context) error {
		middlewareCalled = true
		c.Header("X-Test-Middleware", "executed")
		return nil
	})

	server := httptest.NewServer(app.GetRouter())
	defer server.Close()

	resp, err := http.Get(server.URL + "/test")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if !middlewareCalled {
		t.Error("Middleware was not called")
	}

	if resp.Header.Get("X-Test-Middleware") != "executed" {
		t.Error("Middleware did not set expected header")
	}
}

// TestContext tests the context functionality
func TestContext(t *testing.T) {
	app := setupTestApp()

	// Add route that tests context methods
	app.GET("/context/:id", func(c *gojango.Context) error {
		id := c.Param("id")
		query := c.Query("q")

		response := map[string]string{
			"id":     id,
			"query":  query,
			"method": c.Method(),
			"path":   c.Path(),
		}

		return c.JSON(response)
	})

	server := httptest.NewServer(app.GetRouter())
	defer server.Close()

	resp, err := http.Get(server.URL + "/context/123?q=test")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	var result map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if result["id"] != "123" {
		t.Errorf("Expected id '123', got '%s'", result["id"])
	}

	if result["query"] != "test" {
		t.Errorf("Expected query 'test', got '%s'", result["query"])
	}

	if result["method"] != "GET" {
		t.Errorf("Expected method 'GET', got '%s'", result["method"])
	}
}

// BenchmarkBasicRequest benchmarks basic request handling
func BenchmarkBasicRequest(b *testing.B) {
	app := setupTestApp()
	server := httptest.NewServer(app.GetRouter())
	defer server.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp, err := http.Get(server.URL + "/test")
		if err != nil {
			b.Fatalf("Request failed: %v", err)
		}
		resp.Body.Close()
	}
}

// Example of how to run tests:
// go test -v
// go test -bench=.
// go test -cover

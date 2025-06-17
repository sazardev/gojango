package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"gojango"
	"gojango/models"
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
	
	// Use in-memory SQLite for testing
	app.config.DatabaseURL = "sqlite://:memory:"
	
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
	server := httptest.NewServer(app.router)
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
	server := httptest.NewServer(app.router)
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
		if err := app.db.Create(user); err != nil {
			t.Fatalf("Failed to create test user: %v", err)
		}
	}
	
	// Test Filter
	qs := app.NewQuerySet(&TestUser{})
	results, err := qs.Filter("name__icontains", "a").All() // Should find Alice and Charlie
	if err != nil {
		t.Fatalf("Filter query failed: %v", err)
	}
	
	// Note: results is an interface{}, need to type assert
	if results == nil {
		t.Error("Expected results, got nil")
	}
	
	// Test Count
	count, err := qs.Count()
	if err != nil {
		t.Fatalf("Count query failed: %v", err)
	}
	
	if count != 3 {
		t.Errorf("Expected count 3, got %d", count)
	}
	
	// Test First
	first, err := qs.OrderBy("name").First()
	if err != nil {
		t.Fatalf("First query failed: %v", err)
	}
	
	if first == nil {
		t.Error("Expected first result, got nil")
	}
	
	// Test Exists
	exists, err := qs.Filter("name", "Alice").Exists()
	if err != nil {
		t.Fatalf("Exists query failed: %v", err)
	}
	
	if !exists {
		t.Error("Expected Alice to exist")
	}
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
	
	server := httptest.NewServer(app.router)
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
	
	server := httptest.NewServer(app.router)
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
	server := httptest.NewServer(app.router)
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

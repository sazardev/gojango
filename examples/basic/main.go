package main

import (
	"log"

	"gojango"
	"gojango/middleware"
	"gojango/models"
)

// User model - similar to Django models
type User struct {
	models.Model
	Name     string `json:"name" db:"name,not_null,size:100"`
	Email    string `json:"email" db:"email,unique,not_null,size:255"`
	Password string `json:"-" db:"password,not_null,size:255"`
	Active   bool   `json:"active" db:"active,default:true"`
}

// TableName defines the table name (like in Django)
func (u *User) TableName() string {
	return "users"
}

// Post model
type Post struct {
	models.Model
	Title   string `json:"title" db:"title,not_null,size:200"`
	Content string `json:"content" db:"content,type:TEXT"`
	UserID  uint   `json:"user_id" db:"user_id,not_null"`
}

func (p *Post) TableName() string {
	return "posts"
}

func main() { // Create application with automatic configuration
	app := gojango.New()

	// Configure database (SQLite by default)
	app.GetConfig().DatabaseURL = "sqlite://./app.db"

	// Initialize database connection
	if err := app.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Auto-migration (like Django migrate)
	if err := app.AutoMigrate(&User{}, &Post{}); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	// Global middleware (like Django middleware)
	app.Use(func(c *gojango.Context) error {
		return middleware.Logger()(c)
	})
	app.Use(func(c *gojango.Context) error {
		return middleware.CORS("*")(c)
	})
	app.Use(func(c *gojango.Context) error {
		return middleware.Recovery()(c)
	})

	// Automatic CRUD (like Django admin)
	app.RegisterCRUD("/api/users", &User{})
	app.RegisterCRUD("/api/posts", &Post{})

	// Custom routes (like Django URLs)
	app.GET("/", homeHandler)
	app.GET("/api/health", healthHandler)
	app.POST("/api/login", loginHandler)
	app.GET("/api/users/:id/posts", userPostsHandler)

	// Routes with specific middleware (temporary - without groups for now)
	app.GET("/admin/dashboard", func(c *gojango.Context) error {
		// Here you would apply middleware manually if needed
		return adminDashboardHandler(c)
	})

	log.Println("🚀 GoJango app running on :8000")
	log.Println("📝 API endpoints:")
	log.Println("   GET    /")
	log.Println("   GET    /api/health")
	log.Println("   POST   /api/login")
	log.Println("   GET    /api/users (CRUD)")
	log.Println("   POST   /api/users (CRUD)")
	log.Println("   GET    /api/users/:id (CRUD)")
	log.Println("   PUT    /api/users/:id (CRUD)")
	log.Println("   DELETE /api/users/:id (CRUD)")
	log.Println("   GET    /api/posts (CRUD)")
	log.Println("   POST   /api/posts (CRUD)")
	log.Println("   GET    /api/posts/:id (CRUD)")
	log.Println("   PUT    /api/posts/:id (CRUD)")
	log.Println("   DELETE /api/posts/:id (CRUD)")

	// Start server
	if err := app.Run(":8000"); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

// Handlers - simple and clean like Django views
func homeHandler(c *gojango.Context) error {
	return c.JSON(map[string]interface{}{
		"message": "Welcome to GoJango! 🐍🐹",
		"version": "1.0.0",
		"docs":    "https://github.com/sazardev/gojango",
	})
}

func healthHandler(c *gojango.Context) error {
	return c.JSON(map[string]string{
		"status": "ok",
		"time":   "2025-06-17T10:00:00Z",
	})
}

func loginHandler(c *gojango.Context) error {
	var loginData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BindJSON(&loginData); err != nil {
		return c.ErrorJSON(400, "Invalid JSON", err)
	}

	// Here you would implement authentication logic
	// For simplicity, we accept any email/password
	if loginData.Email == "" || loginData.Password == "" {
		return c.ErrorJSON(400, "Email and password required", nil)
	}

	// Simulate JWT token
	token := "fake-jwt-token-" + loginData.Email

	return c.JSON(map[string]interface{}{
		"token": token,
		"user": map[string]string{
			"email": loginData.Email,
		},
	})
}

func userPostsHandler(c *gojango.Context) error {
	userID := c.Param("id")
	if userID == "" {
		return c.ErrorJSON(400, "User ID is required", nil)
	}

	// Here you would do the actual database query
	// For simplicity, we return mock data
	posts := []map[string]interface{}{
		{
			"id":      1,
			"title":   "My first post",
			"content": "This is the post content",
			"user_id": userID,
		},
		{
			"id":      2,
			"title":   "Second post",
			"content": "More content here",
			"user_id": userID,
		},
	}

	return c.JSON(map[string]interface{}{
		"user_id": userID,
		"posts":   posts,
		"count":   len(posts),
	})
}

func adminDashboardHandler(c *gojango.Context) error {
	return c.JSON(map[string]interface{}{
		"message":    "Admin Dashboard",
		"user_count": 42,
		"post_count": 128,
	})
}

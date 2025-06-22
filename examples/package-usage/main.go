package main

import (
	"log"

	"github.com/sazardev/gojango"
	"github.com/sazardev/gojango/middleware"
	"github.com/sazardev/gojango/models"
)

// User model - like Django models
type User struct {
	models.Model
	Name     string `json:"name" db:"name,not_null,size:100"`
	Email    string `json:"email" db:"email,unique,not_null,size:255"`
	Password string `json:"-" db:"password,not_null,size:255"`
	Active   bool   `json:"active" db:"active,default:true"`
}

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

func main() {
	// Create GoJango application
	app := gojango.New()

	// Configure database
	// Use mock for development/testing (no CGO required)
	app.GetConfig().DatabaseURL = "mock://"

	// Use SQLite for production (requires CGO)
	// app.GetConfig().DatabaseURL = "sqlite://./app.db"

	// Initialize database
	if err := app.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Auto-migrate models (like Django migrate)
	if err := app.AutoMigrate(&User{}, &Post{}); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	// Add middleware (like Django middleware)
	app.Use(func(c *gojango.Context) error {
		return middleware.Logger()(c)
	})
	app.Use(func(c *gojango.Context) error {
		return middleware.CORS("*")(c)
	})

	// Register automatic CRUD endpoints (like Django admin)
	app.RegisterCRUD("/api/users", &User{})
	app.RegisterCRUD("/api/posts", &Post{})

	// Custom routes (like Django URLs)
	app.GET("/", func(c *gojango.Context) error {
		return c.JSON(map[string]interface{}{
			"message": "Welcome to GoJango! 🐍🐹",
			"version": "1.0.1",
			"endpoints": map[string][]string{
				"users": {
					"GET    /api/users",
					"POST   /api/users",
					"GET    /api/users/:id",
					"PUT    /api/users/:id",
					"DELETE /api/users/:id",
				},
				"posts": {
					"GET    /api/posts",
					"POST   /api/posts",
					"GET    /api/posts/:id",
					"PUT    /api/posts/:id",
					"DELETE /api/posts/:id",
				},
			},
		})
	})

	app.GET("/api/health", func(c *gojango.Context) error {
		return c.JSON(map[string]string{
			"status":    "ok",
			"framework": "GoJango",
		})
	})

	// Route groups (like Django URL namespaces)
	api := app.Group("/api/v1")
	api.GET("/status", func(c *gojango.Context) error {
		return c.JSON(map[string]string{"status": "v1 API working"})
	})

	log.Println("🚀 GoJango server starting...")
	log.Println("📖 Visit http://localhost:8000 for API info")
	log.Println("🔧 Available endpoints:")
	log.Println("   GET    /                (API info)")
	log.Println("   GET    /api/health      (health check)")
	log.Println("   GET    /api/users       (list users)")
	log.Println("   POST   /api/users       (create user)")
	log.Println("   GET    /api/users/:id   (get user)")
	log.Println("   PUT    /api/users/:id   (update user)")
	log.Println("   DELETE /api/users/:id   (delete user)")
	log.Println("   GET    /api/posts       (list posts)")
	log.Println("   POST   /api/posts       (create post)")

	// Start server
	if err := app.Run(":8000"); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

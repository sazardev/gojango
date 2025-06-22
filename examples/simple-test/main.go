package main

import (
	"log"

	"gojango"
	"gojango/models"
)

type User struct {
	models.Model
	Name  string `json:"name" db:"name,not_null"`
	Email string `json:"email" db:"email,unique,not_null"`
}

func main() {
	app := gojango.New()
	app.GetConfig().DatabaseURL = "sqlite://./app.db"

	// Initialize database with error handling
	if err := app.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Auto-migrate with error handling
	if err := app.AutoMigrate(&User{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	app.RegisterCRUD("/api/users", &User{})

	log.Println("🚀 Server starting on :8000")
	if err := app.Run(":8000"); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

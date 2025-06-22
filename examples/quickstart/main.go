package main

import (
	"log"

	"github.com/sazardev/gojango"
	"github.com/sazardev/gojango/models"
)

// Define tu modelo (como en Django models.py)
type User struct {
	models.Model
	Name   string `json:"name" db:"name,not_null,size:100"`
	Email  string `json:"email" db:"email,unique,not_null,size:255"`
	Active bool   `json:"active" db:"active,default:true"`
}

func (u *User) TableName() string {
	return "users"
}

func main() {
	// Configuración automática
	app := gojango.New()

	// Configurar base de datos
	app.GetConfig().DatabaseURL = "sqlite://./app.db"

	// Inicializar base de datos
	if err := app.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Auto-migración (como Django migrate)
	app.AutoMigrate(&User{})

	// CRUD automático (como Django admin)
	app.RegisterCRUD("/api/users", &User{})

	// Rutas personalizadas (como Django URLs)
	app.GET("/", func(c *gojango.Context) error {
		return c.JSON(map[string]string{"message": "Hello GoJango!"})
	})

	// Ejecutar servidor
	log.Println("🚀 GoJango server starting on :8000")
	app.Run(":8000")
}

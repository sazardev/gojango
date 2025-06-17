package main

import (
	"log"

	"gojango"
	"gojango/models"
)

// User model extendido
type User struct {
	models.Model
	Name     string `json:"name" db:"name,not_null,size:100"`
	Email    string `json:"email" db:"email,unique,not_null,size:255"`
	Password string `json:"-" db:"password,not_null,size:255"`
	Active   bool   `json:"active" db:"active,default:true"`
	Age      int    `json:"age" db:"age"`
}

func (u *User) TableName() string {
	return "users"
}

func main() {
	// Crear aplicación
	app := gojango.New()
	
	// Migrar modelos
	app.AutoMigrate(&User{})
	
	// Ejemplos de QuerySet (estilo Django ORM)
	demonstrateQuerySet(app)
	
	// Rutas con QuerySet
	app.GET("/api/users/active", func(c *gojango.Context) error {
		// Usuarios activos
		qs := app.NewQuerySet(&User{})
		users, err := qs.Filter("active", true).All()
		if err != nil {
			return c.ErrorJSON(500, "Database error", err)
		}
		return c.JSON(users)
	})
	
	app.GET("/api/users/search", func(c *gojango.Context) error {
		name := c.Query("name")
		if name == "" {
			return c.ErrorJSON(400, "name parameter required", nil)
		}
		
		// Buscar por nombre (case insensitive, contains)
		qs := app.NewQuerySet(&User{})
		users, err := qs.Filter("name__icontains", name).OrderBy("name").All()
		if err != nil {
			return c.ErrorJSON(500, "Database error", err)
		}
		return c.JSON(users)
	})
	
	app.GET("/api/users/adults", func(c *gojango.Context) error {
		// Usuarios mayores de edad
		qs := app.NewQuerySet(&User{})
		users, err := qs.Filter("age__gte", 18).OrderBy("-age").All()
		if err != nil {
			return c.ErrorJSON(500, "Database error", err)
		}
		return c.JSON(users)
	})
	
	app.PUT("/api/users/activate", func(c *gojango.Context) error {
		var request struct {
			UserIDs []int `json:"user_ids"`
		}
		
		if err := c.BindJSON(&request); err != nil {
			return c.ErrorJSON(400, "Invalid JSON", err)
		}
		
		// Activar múltiples usuarios
		qs := app.NewQuerySet(&User{})
		err := qs.Filter("id__in", request.UserIDs).Update(map[string]interface{}{
			"active": true,
		})
		if err != nil {
			return c.ErrorJSON(500, "Database error", err)
		}
		
		return c.JSON(map[string]string{"message": "Users activated"})
	})
	
	log.Println("🚀 Advanced QuerySet demo running on :8000")
	log.Println("📝 Try these endpoints:")
	log.Println("   GET /api/users/active")
	log.Println("   GET /api/users/search?name=john")
	log.Println("   GET /api/users/adults")
	log.Println("   PUT /api/users/activate")
	
	app.Run(":8000")
}

func demonstrateQuerySet(app *gojango.App) {
	log.Println("🔍 Demonstrating Django-like QuerySet operations...")
	
	// Crear algunos usuarios de ejemplo
	users := []*User{
		{Name: "Juan Pérez", Email: "juan@example.com", Age: 25, Active: true},
		{Name: "María García", Email: "maria@example.com", Age: 30, Active: true},
		{Name: "Carlos López", Email: "carlos@example.com", Age: 17, Active: false},
		{Name: "Ana Martín", Email: "ana@example.com", Age: 35, Active: true},
	}
	
	for _, user := range users {
		if err := app.db.Create(user); err != nil {
			log.Printf("Error creating user: %v", err)
		}
	}
	
	qs := app.NewQuerySet(&User{})
	
	// 1. Filtrar usuarios activos
	log.Println("\n1. Usuarios activos:")
	activeUsers, err := qs.Filter("active", true).All()
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		log.Printf("Encontrados: %v", activeUsers)
	}
	
	// 2. Buscar por nombre (contiene)
	log.Println("\n2. Usuarios con 'ar' en el nombre:")
	nameUsers, err := qs.Filter("name__icontains", "ar").All()
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		log.Printf("Encontrados: %v", nameUsers)
	}
	
	// 3. Usuarios mayores de edad ordenados por edad descendente
	log.Println("\n3. Usuarios adultos (age >= 18) ordenados por edad:")
	adultUsers, err := qs.Filter("age__gte", 18).OrderBy("-age").All()
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		log.Printf("Encontrados: %v", adultUsers)
	}
	
	// 4. Contar usuarios activos
	log.Println("\n4. Conteo de usuarios activos:")
	count, err := qs.Filter("active", true).Count()
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		log.Printf("Total usuarios activos: %d", count)
	}
	
	// 5. Primer usuario por orden alfabético
	log.Println("\n5. Primer usuario alfabéticamente:")
	firstUser, err := qs.OrderBy("name").First()
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		log.Printf("Primer usuario: %v", firstUser)
	}
	
	// 6. Verificar si existen usuarios menores de edad
	log.Println("\n6. ¿Existen usuarios menores de edad?")
	exists, err := qs.Filter("age__lt", 18).Exists()
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		log.Printf("Menores de edad: %t", exists)
	}
	
	// 7. Actualizar usuarios inactivos
	log.Println("\n7. Activando usuarios inactivos...")
	err = qs.Filter("active", false).Update(map[string]interface{}{
		"active": true,
	})
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		log.Println("Usuarios actualizados correctamente")
	}
	
	// 8. Consulta compleja: usuarios activos con edad entre 20 y 35
	log.Println("\n8. Usuarios activos entre 20 y 35 años:")
	complexUsers, err := qs.Filter("active", true).Filter("age__gte", 20).Filter("age__lte", 35).OrderBy("age").All()
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		log.Printf("Encontrados: %v", complexUsers)
	}
	
	log.Println("\n✅ QuerySet demonstration completed!")
}

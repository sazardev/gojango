package main

import (
	"log"

	"gojango"
	"gojango/models"
)

// Extended User model
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
	// Create application
	app := gojango.New()
	
	// Migrar modelos
	if err := app.AutoMigrate(&User{}); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
	
	// QuerySet examples (Django ORM style)
	demonstrateQuerySet(app)
	
	// Rutas con QuerySet
	app.GET("/api/users/active", func(c *gojango.Context) error {
		// Active users
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
		
		// Search by name (case insensitive, contains)
		qs := app.NewQuerySet(&User{})
		users, err := qs.Filter("name__icontains", name).OrderBy("name").All()
		if err != nil {
			return c.ErrorJSON(500, "Database error", err)
		}
		return c.JSON(users)
	})
	
	app.GET("/api/users/adults", func(c *gojango.Context) error {
		// Users of legal age
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
		
		// Activate multiple users
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
	
	// Create some example users
	users := []*User{
		{Name: "John Doe", Email: "john@example.com", Age: 25, Active: true},
		{Name: "Jane Smith", Email: "jane@example.com", Age: 30, Active: true},
		{Name: "Bob Wilson", Email: "bob@example.com", Age: 17, Active: false},
		{Name: "Alice Brown", Email: "alice@example.com", Age: 35, Active: true},
	}
	
	for _, user := range users {
		if err := app.GetDB().Create(user); err != nil {
			log.Printf("Error creating user: %v", err)
		}
	}
	
	qs := app.NewQuerySet(&User{})
	
	// 1. Filter active users
	log.Println("\n1. Active users:")
	activeUsers, err := qs.Filter("active", true).All()
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		log.Printf("Found: %v", activeUsers)
	}
	
	// 2. Search by name (contains)
	log.Println("\n2. Users with 'o' in name:")
	nameUsers, err := qs.Filter("name__icontains", "o").All()
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		log.Printf("Found: %v", nameUsers)
	}
	
	// 3. Adult users ordered by age descending
	log.Println("\n3. Adult users (age >= 18) ordered by age:")
	adultUsers, err := qs.Filter("age__gte", 18).OrderBy("-age").All()
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		log.Printf("Found: %v", adultUsers)
	}
	
	// 4. Count active users
	log.Println("\n4. Count of active users:")
	count, err := qs.Filter("active", true).Count()
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		log.Printf("Total active users: %d", count)
	}
	
	// 5. Primer usuario por orden alfabético
	log.Println("\n5. Primer usuario alfabéticamente:")
	firstUser, err := qs.OrderBy("name").First()
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		log.Printf("Primer usuario: %v", firstUser)
	}
	
	// 6. Check if underage users exist
	log.Println("\n6. Do underage users exist?")
	exists, err := qs.Filter("age__lt", 18).Exists()
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		log.Printf("Underage users: %t", exists)
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
	
	// 8. Complex query: active users aged between 20 and 35
	log.Println("\n8. Active users between 20 and 35 years old:")
	complexUsers, err := qs.Filter("active", true).Filter("age__gte", 20).Filter("age__lte", 35).OrderBy("age").All()
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		log.Printf("Found: %v", complexUsers)
	}
	
	log.Println("\n✅ QuerySet demonstration completed!")
}

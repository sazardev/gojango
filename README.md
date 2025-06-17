# GoJango 🐍🐹

A Django-inspired web framework for Go with batteries included. Designed to be **ultra easy to use**, **opinionated**, and with **minimal dependencies**.

## 🌟 Features

- 🚀 **Ultra easy to use** - Django-like syntax
- 🔋 **Batteries included** - ORM, Router, Templates, Auto CRUD, Middleware
- 🏗️ **Clean architecture** - Modular with dependency injection
- 📝 **Opinionated** - Convention over configuration
- 🎯 **Minimal dependencies** - Only SQLite as external dependency
- 🔍 **Django-style QuerySet** - Intuitive ORM with filters and complex queries
- 🛡️ **Built-in middleware** - CORS, Logging, Recovery, Auth, Rate Limiting
- 🌐 **Flexible routing** - Route groups, parameters, specific middleware

## 📦 Installation

```bash
go get github.com/your-username/gojango
```

## 🚀 Quick Start

```go
package main

import (
    "gojango"
    "gojango/models"
)

// Define your model (like Django models.py)
type User struct {
    models.Model
    Name  string `json:"name" db:"name,not_null,size:100"`
    Email string `json:"email" db:"email,unique,not_null,size:255"`
    Active bool  `json:"active" db:"active,default:true"`
}

func (u *User) TableName() string {
    return "users"
}

func main() {
    // Automatic configuration
    app := gojango.New()
    
    // Auto-migration (like Django migrate)
    app.AutoMigrate(&User{})
    
    // Automatic CRUD (like Django admin)
    app.RegisterCRUD("/api/users", &User{})
    
    // Custom routes (like Django URLs)
    app.GET("/", func(c *gojango.Context) error {
        return c.JSON(map[string]string{"message": "Hello GoJango!"})
    })
    
    // Start server
    app.Run(":8000")
}
```

## 🎯 Key Concepts

### 1. Models

Similar to Django models, with tags for database configuration:

```go
type User struct {
    models.Model  // Automatic ID, CreatedAt, UpdatedAt
    Name     string `json:"name" db:"name,not_null,size:100"`
    Email    string `json:"email" db:"email,unique,not_null,size:255"`
    Password string `json:"-" db:"password,not_null,size:255"`
    Active   bool   `json:"active" db:"active,default:true"`
    Age      int    `json:"age" db:"age"`
}

func (u *User) TableName() string {
    return "users"
}
```

**Available database tags:**
- `not_null` - Required field
- `unique` - Unique value
- `primary_key` - Primary key
- `auto_increment` - Auto increment
- `size:N` - Maximum size
- `default:value` - Default value
- `type:TYPE` - Specific DB type

### 2. QuerySet (Django-style ORM)

Intuitive and chainable queries:

```go
qs := app.NewQuerySet(&User{})

// Basic filters
users, _ := qs.Filter("active", true).All()

// Advanced filters (Django-style lookups)
users, _ := qs.Filter("name__icontains", "john").All()       // LIKE %john%
users, _ := qs.Filter("age__gte", 18).All()                  // age >= 18
users, _ := qs.Filter("email__endswith", "@gmail.com").All() // email LIKE %@gmail.com

// Ordering
users, _ := qs.OrderBy("name").All()     // ASC
users, _ := qs.OrderBy("-created_at").All() // DESC

// Pagination
users, _ := qs.Limit(10).Offset(20).All()

// Combinations
adults, _ := qs.Filter("active", true).
               Filter("age__gte", 18).
               OrderBy("-age").
               Limit(5).All()

// Useful operations
count, _ := qs.Filter("active", true).Count()
exists, _ := qs.Filter("email", "john@example.com").Exists()
first, _ := qs.OrderBy("created_at").First()

// Bulk updates
qs.Filter("active", false).Update(map[string]interface{}{
    "active": true,
})

// Bulk deletions
qs.Filter("age__lt", 18).Delete()
```

**Available lookups:**
- `exact` - Exact equality (default)
- `iexact` - Case-insensitive equality
- `contains` - Contains (LIKE %value%)
- `icontains` - Case-insensitive contains
- `startswith` - Starts with
- `endswith` - Ends with
- `gt`, `gte` - Greater than, greater or equal
- `lt`, `lte` - Less than, less or equal
- `in` - In a list of values
- `isnull` - Is NULL or not NULL

### 3. Routes and Controllers

```go
// Basic routes
app.GET("/users", listUsers)
app.POST("/users", createUser)
app.PUT("/users/:id", updateUser)
app.DELETE("/users/:id", deleteUser)

// Route groups with middleware
api := app.Group("/api")
api.Use(middleware.CORS("*"))
api.GET("/users", listUsers)

admin := app.Group("/admin")
admin.Use(middleware.BasicAuth("admin", "secret"))
admin.GET("/dashboard", adminDashboard)

// Automatic CRUD
app.RegisterCRUD("/api/users", &User{})
// Automatically generates:
// GET    /api/users     (list)
// POST   /api/users     (create)
// GET    /api/users/:id (get)
// PUT    /api/users/:id (update) 
// DELETE /api/users/:id (delete)
```

### 4. Context (Request/Response)

Rich API for handling requests and responses:

```go
func handler(c *gojango.Context) error {
    // URL parameters
    id := c.Param("id")
    idInt, _ := c.ParamInt("id")
    
    // Query parameters
    name := c.Query("name")
    page, _ := c.QueryInt("page")
    
    // Headers
    auth := c.GetHeader("Authorization")
    c.Header("X-Custom", "value")
    
    // Body parsing
    var user User
    if err := c.BindJSON(&user); err != nil {
        return c.ErrorJSON(400, "Invalid JSON", err)
    }
    
    // Responses
    return c.JSON(user)                    // JSON response
    return c.String("Hello World")         // Text response
    return c.HTML("<h1>Hello</h1>")       // HTML response
    return c.Render("template.html", data) // Template response
    
    // Helpers
    ip := c.ClientIP()
    isAjax := c.IsAjax()
    isJSON := c.IsJSON()
    
    return nil
}
```

### 5. Middleware

Built-in and easy-to-use middleware:

```go
import "gojango/middleware"

// Global middleware
app.Use(func(c *gojango.Context) error {
    return middleware.Logger()(c)
})
app.Use(func(c *gojango.Context) error {
    return middleware.CORS("*")(c)
})
app.Use(func(c *gojango.Context) error {
    return middleware.Recovery()(c)
})

// Specific middleware
admin := app.Group("/admin")
admin.Use(func(c *gojango.Context) error {
    return middleware.BasicAuth("admin", "secret")(c)
})

// Custom middleware
app.Use(func(c *gojango.Context) error {
    log.Printf("Request: %s %s", c.Method(), c.Path())
    return nil
})
```

**Built-in middleware:**
- `Logger()` - Request logging
- `CORS(origin)` - CORS headers
- `Recovery()` - Panic recovery
- `BasicAuth(user, pass)` - Basic authentication
- `RequestID()` - Unique request ID
- `RateLimit(req, window)` - Request rate limiting
- `Security()` - Security headers

## 📁 Recommended project structure

```
my-project/
├── main.go
├── models/
│   ├── user.go
│   └── post.go
├── handlers/
│   ├── auth.go
│   └── api.go
├── middleware/
│   └── custom.go
├── templates/
│   ├── index.html
│   └── layout.html
├── static/
│   ├── css/
│   └── js/
└── go.mod
```

## 🔧 Configuration

```go
// Default configuration
app := gojango.New()

// Custom configuration
config := config.New()
config.DatabaseURL = "sqlite://./my-app.db"
config.Debug = true
config.Set("app.name", "My App")

app := gojango.New(gojango.WithConfig(config))

// Environment variables
config.LoadFromEnv("MYAPP_") // Load vars starting with MYAPP_

// Using configuration
appName := app.config.GetString("app.name", "Default App")
debug := app.config.GetBool("debug", false)
```

## 🗄️ Database

Uses SQLite by default, perfect for development and small applications:

```go
// SQLite in memory (default)
app := gojango.New()

// SQLite file
app.config.DatabaseURL = "sqlite://./app.db"

// Auto-migration
app.AutoMigrate(&User{}, &Post{}, &Comment{})
```

## 🎨 Templates

Built-in template system with helper functions:

```go
// Configure templates directory
app.templates.SetBaseDir("templates")
app.templates.LoadTemplates()

// Render in handler
func homePage(c *gojango.Context) error {
    data := map[string]interface{}{
        "title": "My App",
        "users": users,
    }
    return c.Render("home.html", data)
}
```

**Available helper functions:**
- `upper`, `lower`, `title` - String manipulation
- `add`, `sub`, `mul`, `div` - Math operations
- `eq`, `ne`, `lt`, `gt` - Comparisons
- `default` - Default values

## 📚 Examples

### Complete REST API

```go
package main

import (
    "gojango"
    "gojango/middleware"
    "gojango/models"
)

type User struct {
    models.Model
    Name  string `json:"name" db:"name,not_null"`
    Email string `json:"email" db:"email,unique,not_null"`
    Active bool  `json:"active" db:"active,default:true"`
}

func (u *User) TableName() string { return "users" }

func main() {
    app := gojango.New()
    
    // Middleware
    app.Use(func(c *gojango.Context) error { return middleware.Logger()(c) })
    app.Use(func(c *gojango.Context) error { return middleware.CORS("*")(c) })
    
    // Migration
    app.AutoMigrate(&User{})
    
    // API routes
    api := app.Group("/api/v1")
    api.RegisterCRUD("/users", &User{})
    
    // Custom routes
    api.GET("/users/search", searchUsers)
    api.POST("/users/bulk", bulkCreateUsers)
    
    app.Run(":8000")
}

func searchUsers(c *gojango.Context) error {
    query := c.Query("q")
    qs := c.app.NewQuerySet(&User{})
    users, err := qs.Filter("name__icontains", query).
                    Filter("active", true).
                    OrderBy("name").All()
    if err != nil {
        return c.ErrorJSON(500, "Search error", err)
    }
    return c.JSON(users)
}
```

### Web application with templates

```go
func main() {
    app := gojango.New()
    app.templates.SetBaseDir("templates")
    app.templates.LoadTemplates()
    
    app.AutoMigrate(&User{})
    
    app.GET("/", homePage)
    app.GET("/users", usersPage)
    app.POST("/users", createUserForm)
    
    app.Run(":8000")
}

func homePage(c *gojango.Context) error {
    return c.Render("home.html", map[string]interface{}{
        "title": "Welcome to GoJango",
    })
}

func usersPage(c *gojango.Context) error {
    qs := c.app.NewQuerySet(&User{})
    users, _ := qs.Filter("active", true).OrderBy("name").All()
    
    return c.Render("users.html", map[string]interface{}{
        "title": "Users",
        "users": users,
    })
}
```

## 🤝 Comparison with Django

| Feature | Django | GoJango |
|---------|--------|---------|
| Models | `models.Model` | `models.Model` |
| ORM | `User.objects.filter()` | `qs.Filter()` |
| URLs | `urlpatterns` | `app.GET()` |
| Views | Functions/Classes | `HandlerFunc` |
| Templates | Jinja-like | Go templates |
| Admin | Automatic | `RegisterCRUD()` |
| Middleware | List in settings | `app.Use()` |
| Migrations | `migrate` | `AutoMigrate()` |

## 📖 Complete documentation

- [Getting started guide](./docs/getting-started.md)
- [Models and ORM](./docs/models.md)
- [Routes and controllers](./docs/routing.md)
- [Middleware](./docs/middleware.md)
- [Templates](./docs/templates.md)
- [Configuration](./docs/config.md)
- [Examples](./examples/)

## 🤔 Why GoJango?

### ✅ Advantages

- **Familiarity**: If you know Django, you'll feel at home
- **Productivity**: Auto CRUD, migrations, built-in middleware
- **Simplicity**: Single external dependency (SQLite)
- **Performance**: Go's speed with Django's comfort
- **Type safety**: Go's type safety
- **Deploy**: Single binary, easy deployment

### 🎯 Ideal for

- Developers coming from Django/Python
- Fast REST APIs
- Small to medium web applications
- Prototypes and MVPs
- Microservices
- Applications that need simple deployment

## 🚀 Get started now

1. **Install Go** (1.22+)
2. **Create a new project**:
   ```bash
   mkdir my-app && cd my-app
   go mod init my-app
   go get github.com/your-username/gojango
   ```
3. **Create `main.go`** with the quick start example
4. **Run**: `go run main.go`
5. **Visit**: `http://localhost:8000`

## 📄 License

MIT License - see [LICENSE](LICENSE) for details.

## 🤝 Contributing

Contributions are welcome! See [CONTRIBUTING.md](CONTRIBUTING.md) for details.

---

**GoJango** - Django-like web framework for Go 🐍🐹

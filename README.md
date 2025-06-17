# GoJango 🐍🐹

Un framework web para Go inspirado en Django con baterías incluidas. Diseñado para ser **ultra fácil de usar**, **opinionado** y con **mínimas dependencias**.

## 🌟 Características

- 🚀 **Ultra fácil de usar** - Sintaxis similar a Django
- 🔋 **Baterías incluidas** - ORM, Router, Templates, CRUD automático, Middleware
- 🏗️ **Arquitectura limpia** - Modular con inyección de dependencias
- 📝 **Opinionado** - Convenciones sobre configuración (Convention over Configuration)
- 🎯 **Mínimas dependencias** - Solo SQLite como dependencia externa
- 🔍 **QuerySet estilo Django** - ORM intuitivo con filtros y consultas complejas
- 🛡️ **Middleware integrado** - CORS, Logging, Recovery, Auth, Rate Limiting
- 🌐 **Rutas flexibles** - Grupos de rutas, parámetros, middleware específico

## 📦 Instalación

```bash
go get github.com/tu-usuario/gojango
```

## 🚀 Inicio rápido

```go
package main

import (
    "gojango"
    "gojango/models"
)

// Define tu modelo (como Django models.py)
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
    // Configuración automática
    app := gojango.New()
    
    // Auto-migración (como Django migrate)
    app.AutoMigrate(&User{})
    
    // CRUD automático (como Django admin)
    app.RegisterCRUD("/api/users", &User{})
    
    // Rutas custom (como Django URLs)
    app.GET("/", func(c *gojango.Context) error {
        return c.JSON(map[string]string{"message": "¡Hola GoJango!"})
    })
    
    // Servidor
    app.Run(":8000")
}
```

## 🎯 Conceptos clave

### 1. Modelos (Models)

Similares a Django models, con tags para configuración de base de datos:

```go
type User struct {
    models.Model  // ID, CreatedAt, UpdatedAt automáticos
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

**Tags de base de datos disponibles:**
- `not_null` - Campo obligatorio
- `unique` - Valor único
- `primary_key` - Clave primaria
- `auto_increment` - Auto incremento
- `size:N` - Tamaño máximo
- `default:valor` - Valor por defecto
- `type:TIPO` - Tipo específico de BD

### 2. QuerySet (ORM estilo Django)

Consultas intuitivas y encadenables:

```go
qs := app.NewQuerySet(&User{})

// Filtros básicos
users, _ := qs.Filter("active", true).All()

// Filtros avanzados (Django-style lookups)
users, _ := qs.Filter("name__icontains", "juan").All()       // LIKE %juan%
users, _ := qs.Filter("age__gte", 18).All()                  // age >= 18
users, _ := qs.Filter("email__endswith", "@gmail.com").All() // email LIKE %@gmail.com

// Ordenamiento
users, _ := qs.OrderBy("name").All()     // ASC
users, _ := qs.OrderBy("-created_at").All() // DESC

// Paginación
users, _ := qs.Limit(10).Offset(20).All()

// Combinaciones
adults, _ := qs.Filter("active", true).
               Filter("age__gte", 18).
               OrderBy("-age").
               Limit(5).All()

// Operaciones útiles
count, _ := qs.Filter("active", true).Count()
exists, _ := qs.Filter("email", "juan@example.com").Exists()
first, _ := qs.OrderBy("created_at").First()

// Actualizaciones masivas
qs.Filter("active", false).Update(map[string]interface{}{
    "active": true,
})

// Eliminaciones masivas
qs.Filter("age__lt", 18).Delete()
```

**Lookups disponibles:**
- `exact` - Igualdad exacta (por defecto)
- `iexact` - Igualdad sin case sensitive
- `contains` - Contiene (LIKE %valor%)
- `icontains` - Contiene sin case sensitive
- `startswith` - Comienza con
- `endswith` - Termina con
- `gt`, `gte` - Mayor que, mayor o igual
- `lt`, `lte` - Menor que, menor o igual
- `in` - En una lista de valores
- `isnull` - Es NULL o no es NULL

### 3. Rutas y Controladores

```go
// Rutas básicas
app.GET("/users", listUsers)
app.POST("/users", createUser)
app.PUT("/users/:id", updateUser)
app.DELETE("/users/:id", deleteUser)

// Grupos de rutas con middleware
api := app.Group("/api")
api.Use(middleware.CORS("*"))
api.GET("/users", listUsers)

admin := app.Group("/admin")
admin.Use(middleware.BasicAuth("admin", "secret"))
admin.GET("/dashboard", adminDashboard)

// CRUD automático
app.RegisterCRUD("/api/users", &User{})
// Genera automáticamente:
// GET    /api/users     (listar)
// POST   /api/users     (crear)
// GET    /api/users/:id (obtener)
// PUT    /api/users/:id (actualizar) 
// DELETE /api/users/:id (eliminar)
```

### 4. Context (Request/Response)

API rica para manejar requests y responses:

```go
func handler(c *gojango.Context) error {
    // Parámetros de URL
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

Middleware integrado y fácil de usar:

```go
import "gojango/middleware"

// Middleware global
app.Use(func(c *gojango.Context) error {
    return middleware.Logger()(c)
})
app.Use(func(c *gojango.Context) error {
    return middleware.CORS("*")(c)
})
app.Use(func(c *gojango.Context) error {
    return middleware.Recovery()(c)
})

// Middleware específico
admin := app.Group("/admin")
admin.Use(func(c *gojango.Context) error {
    return middleware.BasicAuth("admin", "secret")(c)
})

// Middleware custom
app.Use(func(c *gojango.Context) error {
    log.Printf("Request: %s %s", c.Method(), c.Path())
    return nil
})
```

**Middleware incluido:**
- `Logger()` - Log de requests
- `CORS(origin)` - Headers CORS
- `Recovery()` - Recuperación de panics
- `BasicAuth(user, pass)` - Autenticación básica
- `RequestID()` - ID único por request
- `RateLimit(req, window)` - Limitación de requests
- `Security()` - Headers de seguridad

## 📁 Estructura de proyecto recomendada

```
mi-proyecto/
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

## 🔧 Configuración

```go
// Configuración por defecto
app := gojango.New()

// Configuración custom
config := config.New()
config.DatabaseURL = "sqlite://./mi-app.db"
config.Debug = true
config.Set("app.name", "Mi App")

app := gojango.New(gojango.WithConfig(config))

// Variables de entorno
config.LoadFromEnv("MYAPP_") // Carga vars que empiecen con MYAPP_

// Uso de configuración
appName := app.config.GetString("app.name", "Default App")
debug := app.config.GetBool("debug", false)
```

## 🗄️ Base de datos

Por defecto usa SQLite, perfecto para desarrollo y aplicaciones pequeñas:

```go
// SQLite en memoria (por defecto)
app := gojango.New()

// SQLite en archivo
app.config.DatabaseURL = "sqlite://./app.db"

// Auto-migración
app.AutoMigrate(&User{}, &Post{}, &Comment{})
```

## 🎨 Templates

Sistema de templates integrado con funciones helper:

```go
// Configurar directorio de templates
app.templates.SetBaseDir("templates")
app.templates.LoadTemplates()

// Render en handler
func homePage(c *gojango.Context) error {
    data := map[string]interface{}{
        "title": "Mi App",
        "users": users,
    }
    return c.Render("home.html", data)
}
```

**Funciones helper disponibles:**
- `upper`, `lower`, `title` - Manipulación de strings
- `add`, `sub`, `mul`, `div` - Operaciones matemáticas
- `eq`, `ne`, `lt`, `gt` - Comparaciones
- `default` - Valores por defecto

## 📚 Ejemplos

### API REST completa

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
    
    // Migración
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

### Aplicación web con templates

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
        "title": "Bienvenido a GoJango",
    })
}

func usersPage(c *gojango.Context) error {
    qs := c.app.NewQuerySet(&User{})
    users, _ := qs.Filter("active", true).OrderBy("name").All()
    
    return c.Render("users.html", map[string]interface{}{
        "title": "Usuarios",
        "users": users,
    })
}
```

## 🤝 Comparación con Django

| Característica | Django | GoJango |
|---------------|--------|---------|
| Modelos | `models.Model` | `models.Model` |
| ORM | `User.objects.filter()` | `qs.Filter()` |
| URLs | `urlpatterns` | `app.GET()` |
| Views | Functions/Classes | `HandlerFunc` |
| Templates | Jinja-like | Go templates |
| Admin | Automático | `RegisterCRUD()` |
| Middleware | Lista en settings | `app.Use()` |
| Migrations | `migrate` | `AutoMigrate()` |

## 📖 Documentación completa

- [Guía de inicio](./docs/getting-started.md)
- [Modelos y ORM](./docs/models.md)
- [Rutas y controladores](./docs/routing.md)
- [Middleware](./docs/middleware.md)
- [Templates](./docs/templates.md)
- [Configuración](./docs/config.md)
- [Ejemplos](./examples/)

## 🤔 ¿Por qué GoJango?

### ✅ Ventajas

- **Familiaridad**: Si conoces Django, te sentirás como en casa
- **Productividad**: CRUD automático, migraciones, middleware integrado
- **Simplicidad**: Una sola dependencia externa (SQLite)
- **Performance**: La velocidad de Go con la comodidad de Django
- **Tipado**: Seguridad de tipos de Go
- **Deploy**: Binario único, fácil deployment

### 🎯 Ideal para

- Desarrolladores que vienen de Django/Python
- APIs REST rápidas
- Aplicaciones web pequeñas a medianas
- Prototipos y MVPs
- Microservicios
- Aplicaciones que necesitan deployment sencillo

## 🚀 Empezar ahora

1. **Instala Go** (1.22+)
2. **Crea un nuevo proyecto**:
   ```bash
   mkdir mi-app && cd mi-app
   go mod init mi-app
   go get github.com/tu-usuario/gojango
   ```
3. **Crea `main.go`** con el ejemplo de inicio rápido
4. **Ejecuta**: `go run main.go`
5. **Visita**: `http://localhost:8000`

## 📄 Licencia

MIT License - ve [LICENSE](LICENSE) para detalles.

## 🤝 Contribuir

¡Las contribuciones son bienvenidas! Ve [CONTRIBUTING.md](CONTRIBUTING.md) para más detalles.

---

**GoJango** - Django-like web framework for Go 🐍🐹

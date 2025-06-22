# 📦 Cómo Usar GoJango como Paquete

## 🚀 Instalación

### Opción 1: Usando go get
```bash
go get github.com/sazardev/gojango@latest
```

### Opción 2: Usando go mod
```bash
# En tu proyecto
go mod init mi-app-gojango
go get github.com/sazardev/gojango
```

## 📁 Estructura de Proyecto Recomendada

```
mi-app-gojango/
├── main.go
├── go.mod
├── go.sum
├── models/
│   ├── user.go
│   └── post.go
├── handlers/
│   ├── auth.go
│   └── api.go
└── static/
    ├── css/
    └── js/
```

## 🔥 Quick Start

### 1. Crear `main.go`
```go
package main

import (
    "log"

    "github.com/sazardev/gojango"
    "github.com/sazardev/gojango/models"
)

// Define tu modelo
type User struct {
    models.Model
    Name  string `json:"name" db:"name,not_null"`
    Email string `json:"email" db:"email,unique,not_null"`
}

func (u *User) TableName() string {
    return "users"
}

func main() {
    // Crear aplicación
    app := gojango.New()
    
    // Configurar base de datos
    app.GetConfig().DatabaseURL = "sqlite://./app.db"
    app.InitDB()
    
    // Migrar modelos
    app.AutoMigrate(&User{})
    
    // CRUD automático
    app.RegisterCRUD("/api/users", &User{})
    
    // Rutas personalizadas
    app.GET("/", func(c *gojango.Context) error {
        return c.JSON(map[string]string{
            "message": "Hello GoJango!",
            "version": "1.0.0",
        })
    })
    
    // Iniciar servidor
    log.Println("🚀 Server starting on :8000")
    app.Run(":8000")
}
```

### 2. Inicializar el proyecto
```bash
go mod init mi-app-gojango
go get github.com/sazardev/gojango
go run main.go
```

### 3. Probar la API
```bash
# Página principal
curl http://localhost:8000/

# Listar usuarios
curl http://localhost:8000/api/users

# Crear usuario
curl -X POST http://localhost:8000/api/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Juan","email":"juan@example.com"}'

# Obtener usuario por ID
curl http://localhost:8000/api/users/1
```

## 📚 Ejemplos Avanzados

### Con Middleware
```go
import "github.com/sazardev/gojango/middleware"

app.Use(func(c *gojango.Context) error {
    return middleware.Logger()(c)
})
app.Use(func(c *gojango.Context) error {
    return middleware.CORS("*")(c)
})
```

### Con Grupos de Rutas
```go
// API v1
api := app.Group("/api/v1")
api.RegisterCRUD("/users", &User{})
api.GET("/health", healthHandler)

// Admin (con autenticación)
admin := app.Group("/admin")
admin.Use(func(c *gojango.Context) error {
    return middleware.BasicAuth("admin", "secret")(c)
})
admin.GET("/dashboard", adminHandler)
```

### Con QuerySet (ORM estilo Django)
```go
func getUsersHandler(c *gojango.Context) error {
    // Obtener usuarios activos
    qs := c.App().NewQuerySet(&User{})
    users, err := qs.Filter("active", true).
                     OrderBy("name").
                     Limit(10).All()
    
    if err != nil {
        return c.ErrorJSON(500, "Database error", err)
    }
    
    return c.JSON(users)
}
```

### Con Base de Datos Mock para Testing
```go
func TestAPI(t *testing.T) {
    app := gojango.New()
    
    // Usar base de datos mock
    app.GetConfig().DatabaseURL = "mock://"
    app.InitDB()
    app.AutoMigrate(&User{})
    
    // Crear usuario de prueba
    user := &User{Name: "Test", Email: "test@test.com"}
    app.GetDB().Create(user)
    
    // Probar API...
}
```

## 🔧 Configuración Avanzada

### Variables de Entorno
```go
// Usar variables de entorno
config := config.New()
config.LoadFromEnv("MYAPP_")

app := gojango.New(gojango.WithConfig(config))
```

### Base de Datos Personalizada
```go
// Conectar a base de datos externa
db, err := database.Connect("postgres://user:pass@localhost/db")
if err != nil {
    log.Fatal(err)
}

app := gojango.New(gojango.WithDatabase(db))
```

### Templates
```go
app.GET("/home", func(c *gojango.Context) error {
    data := map[string]interface{}{
        "title": "Home Page",
        "users": users,
    }
    return c.Render("home.html", data)
})
```

## 🏗️ Modelos Avanzados

```go
type Post struct {
    models.Model
    Title    string `json:"title" db:"title,not_null,size:200"`
    Content  string `json:"content" db:"content,type:TEXT"`
    UserID   uint   `json:"user_id" db:"user_id,not_null"`
    Tags     string `json:"tags" db:"tags,size:500"`
    Published bool  `json:"published" db:"published,default:false"`
}

func (p *Post) TableName() string {
    return "posts"
}

// Hooks como Django
func (p *Post) BeforeCreate() {
    // Lógica antes de crear
}

func (p *Post) BeforeUpdate() {
    // Lógica antes de actualizar
}
```

## 📋 Tags Disponibles para Campos

- `not_null` - Campo requerido
- `unique` - Valor único
- `primary_key` - Clave primaria
- `auto_increment` - Auto incremento
- `size:N` - Tamaño máximo
- `default:value` - Valor por defecto
- `type:TYPE` - Tipo específico de DB

## 🌐 Endpoints CRUD Automáticos

Cuando usas `app.RegisterCRUD("/api/users", &User{})`, automáticamente obtienes:

- `GET /api/users` - Listar todos
- `POST /api/users` - Crear nuevo
- `GET /api/users/:id` - Obtener por ID
- `PUT /api/users/:id` - Actualizar por ID
- `DELETE /api/users/:id` - Eliminar por ID

## 🚀 Deployment

### Compilar para producción
```bash
CGO_ENABLED=1 go build -o app main.go
./app
```

### Docker
```dockerfile
FROM golang:1.22-alpine AS builder
RUN apk add --no-cache gcc musl-dev sqlite-dev
WORKDIR /app
COPY . .
RUN CGO_ENABLED=1 go build -o main .

FROM alpine:latest
RUN apk add --no-cache sqlite
WORKDIR /root/
COPY --from=builder /app/main .
EXPOSE 8000
CMD ["./main"]
```

## 🤝 Soporte

- **Documentación**: [README.md](../README.md)
- **Ejemplos**: [examples/](../examples/)
- **Issues**: GitHub Issues
- **Contribuir**: [CONTRIBUTING.md](../CONTRIBUTING.md)

---

¡Disfruta construyendo con **GoJango**! 🐍🐹

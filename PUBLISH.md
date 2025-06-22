# 📦 Publicar y Usar GoJango como Paquete

## 🚀 Pasos para Publicar el Paquete

### 1. Configurar Git Repository
```bash
cd gojango
git init
git add .
git commit -m "Initial GoJango framework"
```

### 2. Crear Repository en GitHub
1. Ve a GitHub.com
2. Crea un nuevo repository: `gojango`
3. **Importante**: Usa el nombre de usuario `sazardev` o cambia las importaciones

### 3. Subir a GitHub
```bash
git remote add origin https://github.com/sazardev/gojango.git
git branch -M main
git push -u origin main
```

### 4. Crear Release Tag
```bash
git tag v1.0.0
git push origin v1.0.0
```

### 5. Hacer el Repository Público
- En GitHub, ve a Settings → General → Danger Zone
- Hacer el repository público

## 📥 Cómo Usar el Paquete

### Instalación
```bash
go get github.com/sazardev/gojango@latest
```

### Ejemplo Básico
```go
package main

import (
    "log"
    "github.com/sazardev/gojango"
    "github.com/sazardev/gojango/models"
)

type User struct {
    models.Model
    Name  string `json:"name" db:"name,not_null"`
    Email string `json:"email" db:"email,unique,not_null"`
}

func (u *User) TableName() string {
    return "users"
}

func main() {
    app := gojango.New()
    
    // Configurar base de datos
    app.GetConfig().DatabaseURL = "sqlite://./app.db"
    app.InitDB()
    
    // Migrar y crear CRUD
    app.AutoMigrate(&User{})
    app.RegisterCRUD("/api/users", &User{})
    
    // Ruta personalizada
    app.GET("/", func(c *gojango.Context) error {
        return c.JSON(map[string]string{
            "message": "Hello GoJango!",
        })
    })
    
    app.Run(":8000")
}
```

## 🛠️ Crear Proyecto Nuevo

### Script Automático (Windows)
```powershell
# Ejecutar el script incluido
.\create-example.ps1
```

### Manual
```bash
# 1. Crear proyecto
mkdir mi-app-gojango
cd mi-app-gojango

# 2. Inicializar Go module
go mod init mi-app-gojango

# 3. Instalar GoJango
go get github.com/sazardev/gojango

# 4. Crear main.go (ver ejemplo arriba)

# 5. Ejecutar
go run main.go
```

## 🎯 API Endpoints Automáticos

Con `app.RegisterCRUD("/api/users", &User{})` obtienes:

| Método | Endpoint | Descripción |
|--------|----------|-------------|
| GET | `/api/users` | Listar todos los usuarios |
| POST | `/api/users` | Crear nuevo usuario |
| GET | `/api/users/:id` | Obtener usuario por ID |
| PUT | `/api/users/:id` | Actualizar usuario |
| DELETE | `/api/users/:id` | Eliminar usuario |

## 📝 Ejemplos de Uso

### Crear Usuario
```bash
curl -X POST http://localhost:8000/api/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Juan Pérez","email":"juan@example.com"}'
```

### Listar Usuarios
```bash
curl http://localhost:8000/api/users
```

### Obtener Usuario
```bash
curl http://localhost:8000/api/users/1
```

## 🔧 Configuración Avanzada

### Con Middleware
```go
import "github.com/sazardev/gojango/middleware"

app.Use(func(c *gojango.Context) error {
    return middleware.Logger()(c)
})
```

### Con QuerySet (ORM estilo Django)
```go
func getActiveUsers(c *gojango.Context) error {
    qs := c.App().NewQuerySet(&User{})
    users, err := qs.Filter("active", true).
                     OrderBy("name").
                     Limit(10).All()
    return c.JSON(users)
}
```

## 🧪 Testing con Mock Database

```go
func TestAPI(t *testing.T) {
    app := gojango.New()
    app.GetConfig().DatabaseURL = "mock://"  // Usar mock
    app.InitDB()
    app.AutoMigrate(&User{})
    
    // Tests...
}
```

## 📁 Estructura Recomendada

```
mi-proyecto/
├── main.go              # Archivo principal
├── go.mod               # Dependencias
├── models/              # Modelos de datos
│   ├── user.go
│   └── post.go
├── handlers/            # Controladores
│   ├── auth.go
│   └── users.go
├── middleware/          # Middleware personalizado
│   └── custom.go
├── static/              # Archivos estáticos
│   ├── css/
│   └── js/
└── templates/           # Plantillas HTML
    ├── layout.html
    └── index.html
```

## 🔀 Versiones

### Usar Versión Específica
```bash
go get github.com/sazardev/gojango@v1.0.0
```

### Usar Última Versión
```bash
go get github.com/sazardev/gojango@latest
```

## 📚 Documentación Completa

- [README Principal](README.md) - Documentación completa del framework
- [USAGE.md](USAGE.md) - Guía detallada de uso
- [Ejemplos](examples/) - Proyectos de ejemplo
- [Tests](test/) - Suite de tests completa

## 🤝 Contribuir

1. Fork el proyecto
2. Crea una branch (`git checkout -b feature/nueva-funcionalidad`)
3. Commit los cambios (`git commit -am 'Agregar nueva funcionalidad'`)
4. Push a la branch (`git push origin feature/nueva-funcionalidad`)
5. Crear Pull Request

---

**GoJango** - Framework web inspirado en Django para Go 🐍🐹

¡Disfruta construyendo APIs rápidas y elegantes!

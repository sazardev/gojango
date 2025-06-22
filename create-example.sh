#!/bin/bash

# Script para demostrar cómo usar GoJango como paquete

echo "🚀 Creando proyecto de ejemplo con GoJango..."

# Crear directorio temporal para el proyecto de ejemplo
mkdir -p /tmp/gojango-example
cd /tmp/gojango-example

# Inicializar módulo Go
echo "📦 Inicializando módulo Go..."
go mod init gojango-example

# Crear main.go de ejemplo
cat > main.go << 'EOF'
package main

import (
    "log"

    "github.com/sazardev/gojango"
    "github.com/sazardev/gojango/models"
)

// Modelo de ejemplo
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
    // Crear aplicación GoJango
    app := gojango.New()
    
    // Configurar base de datos
    app.GetConfig().DatabaseURL = "sqlite://./example.db"
    if err := app.InitDB(); err != nil {
        log.Fatalf("Failed to initialize database: %v", err)
    }
    
    // Auto-migración
    app.AutoMigrate(&User{})
    
    // CRUD automático
    app.RegisterCRUD("/api/users", &User{})
    
    // Rutas personalizadas
    app.GET("/", func(c *gojango.Context) error {
        return c.JSON(map[string]interface{}{
            "message": "¡Hola desde GoJango! 🐍🐹",
            "version": "1.0.0",
            "endpoints": []string{
                "GET /api/users",
                "POST /api/users", 
                "GET /api/users/:id",
                "PUT /api/users/:id",
                "DELETE /api/users/:id",
            },
        })
    })
    
    app.GET("/health", func(c *gojango.Context) error {
        return c.JSON(map[string]string{
            "status": "OK",
            "framework": "GoJango",
        })
    })
    
    // Iniciar servidor
    log.Println("🚀 GoJango server starting on :8000")
    log.Println("📝 Endpoints disponibles:")
    log.Println("   GET    / (home)")
    log.Println("   GET    /health")
    log.Println("   GET    /api/users (listar usuarios)")
    log.Println("   POST   /api/users (crear usuario)")
    log.Println("   GET    /api/users/:id (obtener usuario)")
    log.Println("   PUT    /api/users/:id (actualizar usuario)")
    log.Println("   DELETE /api/users/:id (eliminar usuario)")
    
    if err := app.Run(":8000"); err != nil {
        log.Fatalf("Server failed to start: %v", err)
    }
}
EOF

echo "📋 Archivo main.go creado:"
echo "--------------------------------"
cat main.go
echo "--------------------------------"

echo ""
echo "🔧 Para usar este ejemplo:"
echo "1. Copia el proyecto GoJango a GitHub"
echo "2. Ejecuta: go get github.com/sazardev/gojango"
echo "3. Ejecuta: go run main.go"
echo "4. Visita: http://localhost:8000"

echo ""
echo "📝 Comandos de prueba:"
echo "curl http://localhost:8000/"
echo "curl http://localhost:8000/health"
echo "curl http://localhost:8000/api/users"
echo "curl -X POST http://localhost:8000/api/users -H 'Content-Type: application/json' -d '{\"name\":\"Juan\",\"email\":\"juan@example.com\"}'"

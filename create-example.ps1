# Script para demostrar cómo usar GoJango como paquete en Windows

Write-Host "🚀 Creando proyecto de ejemplo con GoJango..." -ForegroundColor Green

# Crear directorio temporal para el proyecto de ejemplo
$exampleDir = "$env:TEMP\gojango-example"
New-Item -ItemType Directory -Force -Path $exampleDir | Out-Null
Set-Location $exampleDir

# Inicializar módulo Go
Write-Host "📦 Inicializando módulo Go..." -ForegroundColor Yellow
go mod init gojango-example

# Crear main.go de ejemplo
$mainGoContent = @'
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
'@

$mainGoContent | Out-File -FilePath "main.go" -Encoding UTF8

Write-Host "📋 Archivo main.go creado en: $exampleDir" -ForegroundColor Cyan
Write-Host "📁 Contenido del directorio:" -ForegroundColor Cyan
Get-ChildItem

Write-Host ""
Write-Host "🔧 Para usar este ejemplo:" -ForegroundColor Green
Write-Host "1. Sube el proyecto GoJango a GitHub (github.com/sazardev/gojango)"
Write-Host "2. Ejecuta: go get github.com/sazardev/gojango"
Write-Host "3. Ejecuta: go run main.go"
Write-Host "4. Visita: http://localhost:8000"

Write-Host ""
Write-Host "📝 Comandos de prueba:" -ForegroundColor Yellow
Write-Host "Invoke-RestMethod -Uri 'http://localhost:8000/'"
Write-Host "Invoke-RestMethod -Uri 'http://localhost:8000/health'"
Write-Host "Invoke-RestMethod -Uri 'http://localhost:8000/api/users'"
Write-Host "Invoke-RestMethod -Uri 'http://localhost:8000/api/users' -Method POST -ContentType 'application/json' -Body '{\"name\":\"Juan\",\"email\":\"juan@example.com\"}'"

Write-Host ""
Write-Host "🌐 O usando curl si lo tienes instalado:" -ForegroundColor Yellow
Write-Host "curl http://localhost:8000/"
Write-Host "curl http://localhost:8000/health"
Write-Host "curl http://localhost:8000/api/users"
Write-Host "curl -X POST http://localhost:8000/api/users -H 'Content-Type: application/json' -d '{\"name\":\"Juan\",\"email\":\"juan@example.com\"}'"

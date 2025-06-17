#!/bin/bash

# GoJango Quick Start Script
# This script helps users get started quickly with GoJango

set -e

echo "🐍🐹 GoJango Quick Start"
echo "======================="
echo ""

# Verificar Go instalado
if ! command -v go &> /dev/null; then
    echo "❌ Go no está instalado. Por favor instala Go 1.22+ desde https://golang.org/"
    exit 1
fi

echo "✅ Go encontrado: $(go version)"
echo ""

# Create new project
read -p "📝 Project name: " PROJECT_NAME

if [ -z "$PROJECT_NAME" ]; then
    echo "❌ Project name required"
    exit 1
fi

if [ -d "$PROJECT_NAME" ]; then
    echo "❌ Directory '$PROJECT_NAME' already exists"
    exit 1
fi

echo ""
echo "🚀 Creating project '$PROJECT_NAME'..."

# Create directory structure
mkdir -p "$PROJECT_NAME"/{models,handlers,middleware,templates,static/{css,js}}

cd "$PROJECT_NAME"

# Create go.mod
cat > go.mod << EOF
module $PROJECT_NAME

go 1.22

require (
    github.com/mattn/go-sqlite3 v1.14.17
)
EOF

# Create basic main.go
cat > main.go << 'EOF'
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// Simple GoJango-inspired app
type App struct {
	mux *http.ServeMux
}

type Context struct {
	w http.ResponseWriter
	r *http.Request
}

func New() *App {
	return &App{mux: http.NewServeMux()}
}

func (app *App) GET(pattern string, handler func(*Context)) {
	app.mux.HandleFunc("GET "+pattern, func(w http.ResponseWriter, r *http.Request) {
		ctx := &Context{w: w, r: r}
		handler(ctx)
	})
}

func (app *App) POST(pattern string, handler func(*Context)) {
	app.mux.HandleFunc("POST "+pattern, func(w http.ResponseWriter, r *http.Request) {
		ctx := &Context{w: w, r: r}
		handler(ctx)
	})
}

func (c *Context) JSON(data interface{}) {
	c.w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(c.w).Encode(data)
}

func (c *Context) Param(name string) string {
	return c.r.PathValue(name)
}

func (app *App) Run(addr string) error {
	log.Printf("🚀 %s server running on %s", "GoJango App", addr)
	return http.ListenAndServe(addr, app.mux)
}

func main() {
	app := New()

	// Rutas principales
	app.GET("/", func(c *Context) {
		c.JSON(map[string]interface{}{
			"message":   "¡Bienvenido a tu app GoJango! 🐍🐹",
			"timestamp": time.Now().Format(time.RFC3339),
			"status":    "running",
			"endpoints": []string{
				"GET /",
				"GET /health",
				"GET /api/hello/{name}",
			},
		})
	})

	app.GET("/health", func(c *Context) {
		c.JSON(map[string]string{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	app.GET("/api/hello/{name}", func(c *Context) {
		name := c.Param("name")
		if name == "" {
			name = "Anonymous"
		}
		c.JSON(map[string]string{
			"message": fmt.Sprintf("Hello %s! 👋", name),
			"time":    time.Now().Format("15:04:05"),
		})
	})

	log.Println("📝 Endpoints disponibles:")
	log.Println("   GET  /              - Página principal")
	log.Println("   GET  /health        - Health check")
	log.Println("   GET  /api/hello/{name} - Saludo personalizado")
	log.Println("")
	log.Println("🎯 Prueba:")
	log.Println("   curl http://localhost:8000/")
	log.Println("   curl http://localhost:8000/api/hello/Juan")
	log.Println("")

	if err := app.Run(":8000"); err != nil {
		log.Fatalf("❌ Error starting server: %v", err)
	}
}
EOF

# Create project README
cat > README.md << EOF
# $PROJECT_NAME

Web application created with GoJango 🐍🐹

## Run

\`\`\`bash
go run main.go
\`\`\`

Luego visita: http://localhost:8000

## Desarrollo

\`\`\`bash
# Instalar dependencias
go mod tidy

# Run
go run main.go

# Compilar
go build

# Testear
go test
\`\`\`

## Endpoints

- \`GET /\` - Página principal
- \`GET /health\` - Health check  
- \`GET /api/hello/{name}\` - Saludo personalizado

## Estructura

\`\`\`
$PROJECT_NAME/
├── main.go           # Main application
├── models/           # Data models
├── handlers/         # Controladores
├── middleware/       # Middleware custom
├── templates/        # Templates HTML
├── static/           # Assets estáticos
└── README.md        # This file
\`\`\`

---
Created with ❤️ using GoJango
EOF

# Create .gitignore
cat > .gitignore << 'EOF'
# Binarios
*.exe
*.exe~
*.dll
*.so
*.dylib
main

# Test binaries
*.test
*.out

# Dependency directories
vendor/

# Go workspace file
go.work

# IDE
.vscode/
.idea/
*.swp
*.swo

# OS
.DS_Store
Thumbs.db

# App specific
*.db
*.sqlite
*.log
EOF

echo ""
echo "✅ Project '$PROJECT_NAME' created successfully!"
echo ""
echo "🎯 Próximos pasos:"
echo "   cd $PROJECT_NAME"
echo "   go mod tidy"
echo "   go run main.go"
echo ""
echo "🌐 Luego visita: http://localhost:8000"
echo ""
echo "🎉 ¡Happy coding con GoJango! 🐍🐹"

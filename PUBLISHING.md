# 🚀 Publicar GoJango como Paquete Go

## 1. Preparar el Repositorio

### Crear repositorio en GitHub:
```bash
# Ir al directorio raíz del proyecto
cd c:\Users\cerbe\Documents\go\gojango

# Inicializar git si no está inicializado
git init

# Crear archivo .gitignore
echo "*.exe" > .gitignore
echo "*.log" >> .gitignore
echo ".DS_Store" >> .gitignore
echo "app.db" >> .gitignore

# Agregar archivos
git add .
git commit -m "Initial commit: GoJango framework with mock database support"

# Agregar remote (reemplaza tu-usuario con tu usuario de GitHub)
git remote add origin https://github.com/tu-usuario/gojango.git

# Push al repositorio
git push -u origin main
```

## 2. Crear Release con Tag

```bash
# Crear tag de versión
git tag v1.0.1
git push origin v1.0.1
```

## 3. Usar el Paquete Publicado

### go.mod para un nuevo proyecto:
```go
module mi-proyecto

go 1.22

require github.com/tu-usuario/gojango v1.0.1
```

### Ejemplo de uso:
```go
package main

import (
    "log"
    "github.com/tu-usuario/gojango"
    "github.com/tu-usuario/gojango/models"
)

type User struct {
    models.Model
    Name  string `json:"name" db:"name,not_null"`
    Email string `json:"email" db:"email,unique,not_null"`
}

func main() {
    app := gojango.New()
    
    // Para desarrollo con mock (sin CGO)
    app.GetConfig().DatabaseURL = "mock://"
    
    // Para producción con SQLite (requiere CGO)
    // app.GetConfig().DatabaseURL = "sqlite://./app.db"
    
    if err := app.InitDB(); err != nil {
        log.Fatalf("Failed to initialize database: %v", err)
    }
    
    if err := app.AutoMigrate(&User{}); err != nil {
        log.Fatalf("Failed to migrate database: %v", err)
    }
    
    app.RegisterCRUD("/api/users", &User{})
    
    app.Run(":8000")
}
```

## 4. Instalación para Usuarios

```bash
# Crear nuevo proyecto
mkdir mi-app-gojango
cd mi-app-gojango

# Inicializar módulo Go
go mod init mi-app

# Instalar GoJango
go get github.com/tu-usuario/gojango@v1.0.1

# Crear main.go con el código del ejemplo anterior
# Ejecutar
go run main.go
```

## 5. Distribución con Dockerfile

```dockerfile
FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build with CGO enabled for SQLite
ENV CGO_ENABLED=1
RUN apk add --no-cache gcc musl-dev sqlite-dev
RUN go build -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates sqlite
WORKDIR /root/

COPY --from=builder /app/main .

CMD ["./main"]
```

## 6. Comandos de Desarrollo

```bash
# Ejecutar tests
go test ./...

# Build para múltiples plataformas
GOOS=linux GOARCH=amd64 go build -o gojango-linux-amd64
GOOS=windows GOARCH=amd64 go build -o gojango-windows-amd64.exe
GOOS=darwin GOARCH=amd64 go build -o gojango-darwin-amd64
```

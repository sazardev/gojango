# 📦 Cómo usar GoJango como Paquete

Este ejemplo muestra cómo usar GoJango después de publicarlo.

## Instalación

```bash
# Crear nuevo proyecto
mkdir mi-app-gojango
cd mi-app-gojango

# Inicializar módulo Go
go mod init mi-app

# Instalar GoJango (reemplaza con tu repositorio)
go get github.com/tu-usuario/gojango@latest
```

## Uso Básico

Ver `main.go` para un ejemplo completo.

## Comandos

```bash
# Instalar dependencias
go mod tidy

# Ejecutar (development mode con mock database)
go run main.go

# Build para producción
CGO_ENABLED=1 go build -o app

# Ejecutar tests
go test
```

## Nota

Este ejemplo usa `github.com/sazardev/gojango` pero debes reemplazarlo con tu repositorio real después de publicar el framework.

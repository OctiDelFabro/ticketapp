# Testing

## Objetivo

El proyecto incorpora pruebas automatizadas para validar la lógica principal del sistema y los endpoints HTTP del backend.

- Objetivo de regularidad: alcanzar una cobertura inicial cercana o superior al 40%.
- Objetivo final: alcanzar una cobertura cercana al 80%.

## Estrategia

- Tests unitarios de services para los flujos principales de autenticación, eventos y tickets.
- Tests de integración HTTP de controllers/routes con `httptest` y el router real de Gin.
- Base de datos SQLite en memoria solo para tests con GORM y un driver pure-Go.
- Los tests no llaman a `config.ConnectDatabase()`, no leen el `.env` real y no requieren MySQL corriendo.
- Los tests no requieren CGO ni gcc, por lo que pueden ejecutarse en Windows, macOS y Linux sin compiladores C.
- Los tests actuales cubren services y controllers principales del flujo cliente.

## Comandos

Desde la raíz del repositorio:

```bash
cd backend
go test ./...
go test ./... -coverpkg=./...
go test ./... -coverpkg=./... -coverprofile=coverage.out
go tool cover -func=coverage.out
```

También se recomienda ejecutar el orden completo usado para validar cambios de backend:

```bash
cd backend
gofmt -w tests/*.go
go mod tidy
go test ./...
go test ./... -coverpkg=./...
go test ./... -coverpkg=./... -coverprofile=coverage.out
go tool cover -func=coverage.out
```

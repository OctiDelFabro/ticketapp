# Testing

## Objetivo

El proyecto incorpora pruebas automatizadas para validar la lógica principal del sistema y los endpoints HTTP del backend. Para la entrega final se validó localmente una cobertura de **84.2%**.

## Estrategia

- Tests unitarios de services para los flujos principales de autenticación, eventos, tickets y métricas.
- Tests de integración HTTP de controllers/routes con `httptest` y el router real de Gin.
- Base de datos SQLite in-memory solo para tests con GORM.
- Los tests no llaman a `config.ConnectDatabase()`, no leen el `.env` real y no requieren MySQL corriendo.
- Los tests no requieren Docker ni Docker Compose.
- Los tests actuales cubren services, controllers, middlewares, rutas, DAOs y utilidades principales.

## Comandos finales

Desde la raíz del repositorio:

```bash
cd backend
go test ./...
go test ./... -cover
go test ./tests "-coverpkg=github.com/OctiDelFabro/ticketapp/backend/controllers,github.com/OctiDelFabro/ticketapp/backend/services,github.com/OctiDelFabro/ticketapp/backend/dao,github.com/OctiDelFabro/ticketapp/backend/routes,github.com/OctiDelFabro/ticketapp/backend/middlewares,github.com/OctiDelFabro/ticketapp/backend/utils" "-coverprofile=coverage.out"
go tool cover "-func=coverage.out"
```

## Coverage

- Cobertura final validada localmente: **84.2%**.
- `coverage.out` es un artefacto local de medición y no debe commitearse.
- La suite usa SQLite in-memory, por lo que no depende de MySQL ni de Docker para ejecutarse.

## Testing frontend

El frontend incorpora tests unitarios y de componentes con **Vitest**, **React Testing Library** y **jsdom**. La suite valida componentes de UI, guards de navegación admin y utilidades de imágenes sin depender del backend real.

Comandos:

```bash
cd frontend
npm install
npm run test:run
```

Estos tests no requieren backend, MySQL, Docker ni Docker Compose. Cuando un componente necesita datos externos, se usan mocks de testing en lugar de servicios reales.

## Relación con Docker

Docker Compose se usa para levantar el sistema completo en modo demo/desarrollo, pero no es requisito para correr los tests automatizados. Esto permite validar la suite en entornos locales o CI sin levantar servicios externos.

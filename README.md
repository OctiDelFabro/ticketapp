# TicketApp

Sistema de Gestión de Eventos y Entradas tipo Ticketek desarrollado para la materia Desarrollo de Software 2026.

## Descripción

TicketApp permitirá gestionar eventos y entradas mediante un backend en Go, un frontend en React y una base de datos MySQL.

El proyecto se trabajará con Spec Driven Development, manteniendo la documentación de especificación, API, base de datos, frontend, testing y roadmap dentro de `docs/`.

## Estructura inicial

```txt
ticketapp/
├── backend/
│   └── .gitkeep
├── frontend/
│   └── .gitkeep
├── docs/
│   ├── SPEC.md
│   ├── API.md
│   ├── DATABASE.md
│   ├── FRONTEND.md
│   ├── TESTING.md
│   └── ROADMAP.md
├── .gitignore
└── README.md
```

## Setup local

Para levantar el proyecto localmente con base MySQL, migraciones y datos demo reproducibles, seguí la guía de [setup local](docs/LOCAL_SETUP.md).


## Docker

### Requisitos

- Docker
- Docker Compose

### Levantar el proyecto

Desde la raíz del repositorio:

```bash
docker compose up --build
```

Este comando levanta MySQL, el backend Go/Gin y el frontend React/Vite servido con Nginx. No hace falta tener MySQL instalado localmente para usar Docker.

### URLs y puertos

- Frontend: `http://localhost:5173`
- Backend health: `http://localhost:8080/api/health`
- MySQL desde el host: `localhost:3307`
- MySQL dentro de Docker: `mysql:3306`

Docker expone MySQL en el puerto `3307` del host para evitar conflictos si ya usás MySQL local en `3306`. Dentro de la red de Docker, el backend se conecta a `mysql:3306`.

### Credenciales demo

- `admin@test.com` / `123456`
- `octavio@test.com` / `123456`
- `lorenzo@test.com` / `123456`
- `pablo@test.com` / `123456`

### Comandos útiles

```bash
docker compose down
docker compose down -v
docker compose logs -f backend
docker compose logs -f frontend
docker compose logs -f mysql
```

`docker compose down -v` borra el volumen local de MySQL del contenedor y resetea la base de datos. Al iniciar, el backend crea, migra y siembra la base con datos demo reproducibles.

## Backend - ejecución local

```bash
cd backend
go mod tidy
go run main.go
```

Al iniciar, el backend conecta con MySQL usando las variables de entorno configuradas, crea la base `ticketapp` si no existe, ejecuta las migraciones GORM de los modelos `User`, `Event` y `Ticket`, y carga datos demo idempotentes: eventos, usuarios cliente y tickets de prueba.

### Variables JWT

El backend usa `JWT_SECRET` para firmar tokens y `JWT_EXPIRATION_HOURS` para configurar su expiración. Si `JWT_SECRET` no está definido, se usa una clave simple de desarrollo que debe cambiarse antes de producción.

```env
JWT_SECRET=super-secret-dev-key
JWT_EXPIRATION_HOURS=24
```

## Testing backend

Los tests automatizados del backend usan SQLite in-memory mediante helpers de GORM. No requieren MySQL ni Docker. El archivo `coverage.out` se genera localmente para analizar cobertura y no debe commitearse.

```bash
cd backend
go test ./...
go test ./... -cover
```

Para medir cobertura real de los paquetes del backend desde la suite de tests:

```bash
cd backend
go test ./tests -coverpkg=./controllers,./services,./dao,./routes,./middlewares,./utils -coverprofile=coverage.out
go tool cover -func=coverage.out
go tool cover -html=coverage.out
```

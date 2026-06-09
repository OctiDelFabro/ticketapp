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

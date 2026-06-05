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

## Backend - ejecución local

```bash
cd backend
go mod tidy
go run main.go
```

Al iniciar, el backend conecta con MySQL usando las variables de entorno configuradas, ejecuta las migraciones GORM de los modelos `User`, `Event` y `Ticket`, y carga cuatro eventos iniciales si la tabla `events` está vacía.

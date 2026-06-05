# API Planificada

## Auth

- `POST /api/auth/register`
- `POST /api/auth/login`

## Eventos públicos

- `GET /api/events`
- `GET /api/events/:id`

## Cliente

- `POST /api/tickets/purchase`
- `GET /api/tickets/me`
- `PATCH /api/tickets/:id/cancel`
- `PATCH /api/tickets/:id/transfer`

## Admin

Funcionalidades planificadas para una etapa posterior:

- `POST /api/admin/events`
- `PUT /api/admin/events/:id`
- `DELETE /api/admin/events/:id`
- `GET /api/admin/events/:id/report`

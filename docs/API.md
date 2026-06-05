# API

## Health

### `GET /api/health`

Verifica que el backend esté activo.

#### Response `200 OK`

```json
{
  "status": "ok",
  "message": "TicketApp backend running"
}
```

## Auth

### `POST /api/auth/register`

Registra un usuario nuevo. El backend guarda la contraseña con bcrypt y asigna siempre el rol `CLIENT`.

#### Request

```json
{
  "name": "Ada Lovelace",
  "email": "ada@example.com",
  "password": "secret123"
}
```

#### Response `201 Created`

```json
{
  "token": "jwt-token",
  "user": {
    "id": 1,
    "name": "Ada Lovelace",
    "email": "ada@example.com",
    "role": "CLIENT"
  }
}
```

#### Errores

- `400 Bad Request`: request inválido, campos incompletos o password menor a 6 caracteres.
- `409 Conflict`: email ya registrado.
- `500 Internal Server Error`: error interno.

### `POST /api/auth/login`

Autentica un usuario existente con email y contraseña.

#### Request

```json
{
  "email": "ada@example.com",
  "password": "secret123"
}
```

#### Response `200 OK`

```json
{
  "token": "jwt-token",
  "user": {
    "id": 1,
    "name": "Ada Lovelace",
    "email": "ada@example.com",
    "role": "CLIENT"
  }
}
```

#### Errores

- `400 Bad Request`: request inválido o campos incompletos.
- `401 Unauthorized`: credenciales inválidas.
- `500 Internal Server Error`: error interno.

> Las respuestas de autenticación nunca incluyen `password_hash`.

## Eventos públicos

Los endpoints públicos de eventos no requieren JWT y devuelven solamente eventos activos.

### `GET /api/events`

Lista el catálogo público de eventos activos.

#### Query params opcionales

- `search`: filtra por título usando coincidencia parcial.
- `category`: filtra por categoría exacta.
- `available_only`: si es `true`, devuelve solamente eventos con cupo disponible.

#### Ejemplos

- `GET /api/events`
- `GET /api/events?search=rock`
- `GET /api/events?category=Música`
- `GET /api/events?available_only=true`

#### Response `200 OK`

```json
[
  {
    "id": 1,
    "title": "Rock Nacional",
    "description": "Evento de música rock nacional.",
    "image_url": "",
    "category": "Música",
    "location": "Córdoba",
    "start_date": "2026-09-12T21:00:00Z",
    "duration_minutes": 120,
    "capacity": 100,
    "available_capacity": 100,
    "active": true
  }
]
```

#### Errores

- `400 Bad Request`: valor inválido para `available_only`.
- `500 Internal Server Error`: error interno.

### `GET /api/events/:id`

Devuelve el detalle público de un evento activo.

#### Response `200 OK`

```json
{
  "id": 1,
  "title": "Rock Nacional",
  "description": "Evento de música rock nacional.",
  "image_url": "",
  "category": "Música",
  "location": "Córdoba",
  "start_date": "2026-09-12T21:00:00Z",
  "duration_minutes": 120,
  "capacity": 100,
  "available_capacity": 100,
  "active": true
}
```

#### Errores

- `400 Bad Request`: id inválido.
- `404 Not Found`: el evento no existe o no está activo.
- `500 Internal Server Error`: error interno.

> `available_capacity` se calcula como `capacity` menos la cantidad de tickets asociados al evento con status `ACTIVE`.

## Cliente

Funcionalidades planificadas para una etapa posterior:

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

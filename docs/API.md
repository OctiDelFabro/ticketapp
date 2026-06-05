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

Funcionalidades planificadas para una etapa posterior:

- `GET /api/events`
- `GET /api/events/:id`

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

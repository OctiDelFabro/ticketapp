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

Los endpoints públicos de eventos no requieren JWT y devuelven solamente eventos activos. El campo `price` representa el valor de una entrada general expresado en ARS.

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
    "price": 15000,
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
  "price": 15000,
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

Los endpoints de tickets requieren autenticación JWT. Enviar el token obtenido en `POST /api/auth/register` o `POST /api/auth/login` usando el header:

```http
Authorization: Bearer <token>
```

Si el header no existe, no usa el formato `Bearer <token>`, o el token es inválido/expiró, la API responde `401 Unauthorized`.

### `POST /api/tickets/purchase`

Compra una entrada para el usuario autenticado.

#### Request

```json
{
  "event_id": 1
}
```

#### Response `201 Created`

```json
{
  "id": 1,
  "event_id": 1,
  "event_title": "Rock Nacional",
  "event_start_date": "2026-09-12T21:00:00Z",
  "event_location": "Córdoba",
  "event_price": 15000,
  "status": "ACTIVE",
  "purchase_date": "2026-06-05T12:00:00Z",
  "user_id": 1,
  "user_email": "ada@example.com"
}
```

#### Reglas

- `event_id` es requerido.
- El evento debe existir y estar activo.
- El usuario autenticado no puede tener otro ticket `ACTIVE` para el mismo evento.
- Debe haber cupo disponible: `event.capacity - tickets ACTIVE del evento`.

#### Errores

- `400 Bad Request`: request inválido o `event_id` faltante.
- `401 Unauthorized`: token faltante, inválido o expirado.
- `404 Not Found`: el evento no existe o no está activo.
- `409 Conflict`: no hay cupo o el usuario ya tiene un ticket activo para ese evento.
- `500 Internal Server Error`: error interno.

### `GET /api/tickets/me`

Lista las entradas del usuario autenticado, ordenadas por `purchase_date` descendente.

#### Response `200 OK`

```json
[
  {
    "id": 1,
    "event_id": 1,
    "event_title": "Rock Nacional",
    "event_start_date": "2026-09-12T21:00:00Z",
    "event_location": "Córdoba",
    "event_price": 15000,
    "status": "ACTIVE",
    "purchase_date": "2026-06-05T12:00:00Z",
    "user_id": 1,
    "user_email": "ada@example.com"
  }
]
```

#### Errores

- `401 Unauthorized`: token faltante, inválido o expirado.
- `500 Internal Server Error`: error interno.

### `PATCH /api/tickets/:id/cancel`

Cancela una entrada activa propia. No elimina físicamente el registro: cambia `status` a `CANCELLED`.

#### Response `200 OK`

```json
{
  "id": 1,
  "event_id": 1,
  "event_title": "Rock Nacional",
  "event_start_date": "2026-09-12T21:00:00Z",
  "event_location": "Córdoba",
  "event_price": 15000,
  "status": "CANCELLED",
  "purchase_date": "2026-06-05T12:00:00Z",
  "user_id": 1,
  "user_email": "ada@example.com"
}
```

#### Errores

- `400 Bad Request`: id inválido.
- `401 Unauthorized`: token faltante, inválido o expirado.
- `403 Forbidden`: el ticket no pertenece al usuario autenticado.
- `404 Not Found`: el ticket no existe.
- `409 Conflict`: el ticket no está activo o ya fue cancelado.
- `500 Internal Server Error`: error interno.

### `PATCH /api/tickets/:id/transfer`

Transfiere una entrada activa propia a otro usuario existente. Mantiene el `status` en `ACTIVE` y actualiza el `user_id` del ticket.

#### Request

```json
{
  "target_email": "otro@test.com"
}
```

#### Response `200 OK`

```json
{
  "id": 1,
  "event_id": 1,
  "event_title": "Rock Nacional",
  "event_start_date": "2026-09-12T21:00:00Z",
  "event_location": "Córdoba",
  "event_price": 15000,
  "status": "ACTIVE",
  "purchase_date": "2026-06-05T12:00:00Z",
  "user_id": 2,
  "user_email": "otro@test.com"
}
```

#### Reglas

- `target_email` es requerido.
- El usuario destino debe existir.
- El usuario destino no puede ser el mismo usuario autenticado.
- El usuario destino no puede tener ya un ticket `ACTIVE` para el mismo evento.

#### Errores

- `400 Bad Request`: id inválido, request inválido o `target_email` faltante.
- `401 Unauthorized`: token faltante, inválido o expirado.
- `403 Forbidden`: el ticket no pertenece al usuario autenticado.
- `404 Not Found`: el ticket o el usuario destino no existe.
- `409 Conflict`: el ticket no está activo, el destino es el mismo usuario, o el destino ya tiene ticket activo para el evento.
- `500 Internal Server Error`: error interno.

> Las respuestas de tickets no devuelven objetos GORM completos ni `password_hash`. `event_price` corresponde al precio actual del evento asociado y está expresado en ARS.


## Admin endpoints

Los endpoints de administración requieren el header `Authorization: Bearer <token>` con un JWT de un usuario con rol `ADMIN`. Para desarrollo local existe el usuario demo `admin@test.com` / `123456`.

> El usuario admin demo es solo para desarrollo local.

### `POST /api/admin/events`

Crea un evento nuevo.

#### Request

```json
{
  "title": "Nombre del evento",
  "description": "Descripción",
  "image_url": "https://...",
  "category": "Música",
  "location": "Córdoba",
  "start_date": "2026-12-10T21:00:00Z",
  "duration_minutes": 120,
  "capacity": 100,
  "price": 15000,
  "active": true
}
```

#### Reglas

- `title`, `description`, `category`, `location` y `start_date` son requeridos.
- `duration_minutes` debe ser mayor a `0`.
- `capacity` debe ser mayor a `0`.
- `price` es opcional, queda en `0` si se omite y debe ser mayor o igual a `0`.

#### Response `201 Created`

Devuelve el evento creado con `available_capacity`.

#### Errores

- `400 Bad Request`: request inválido o validaciones fallidas.
- `401 Unauthorized`: falta autenticación.
- `403 Forbidden`: el usuario autenticado no tiene rol `ADMIN`.
- `500 Internal Server Error`: error interno.

### `PATCH /api/admin/events/:id`

Actualiza parcialmente un evento existente.

#### Campos editables

- `title`
- `description`
- `image_url`
- `category`
- `location`
- `start_date`
- `duration_minutes`
- `capacity`
- `price`
- `active`

#### Reglas

- Si se envía `duration_minutes`, debe ser mayor a `0`.
- Si se envía `capacity`, debe ser mayor a `0`.
- Si se envía `price`, debe ser mayor o igual a `0`.
- `capacity` no puede ser menor que la cantidad de tickets `ACTIVE` vendidos para el evento.

#### Response `200 OK`

Devuelve el evento actualizado.

#### Errores

- `400 Bad Request`: id inválido, request inválido o validaciones fallidas.
- `401 Unauthorized`: falta autenticación.
- `403 Forbidden`: el usuario autenticado no tiene rol `ADMIN`.
- `404 Not Found`: el evento no existe.
- `500 Internal Server Error`: error interno.

### `DELETE /api/admin/events/:id`

Deshabilita un evento con borrado lógico (`active = false`). No borra físicamente el registro.

#### Response `200 OK`

```json
{
  "message": "event disabled successfully"
}
```

#### Errores

- `400 Bad Request`: id inválido.
- `401 Unauthorized`: falta autenticación.
- `403 Forbidden`: el usuario autenticado no tiene rol `ADMIN`.
- `404 Not Found`: el evento no existe.
- `500 Internal Server Error`: error interno.

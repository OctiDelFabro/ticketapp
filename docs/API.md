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
    "tickets_sold": 0,
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

> `available_capacity` se calcula como `capacity` menos la cantidad de tickets asociados al evento con status `ACTIVE`. `tickets_sold` cuenta solamente tickets `ACTIVE`.

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

### `POST /api/tickets/gift`

Crea una entrada nueva como regalo para otro usuario registrado. A diferencia de transferir una entrada, regalar no mueve una entrada existente: compra/emite una entrada nueva directamente a nombre del destinatario.

- **Auth:** requerida (`Authorization: Bearer <token>`).
- **Body:**

```json
{
  "event_id": 1,
  "target_email": "lorenzo@test.com",
  "message": "Feliz cumple!"
}
```

- `event_id` es obligatorio.
- `target_email` es obligatorio y debe pertenecer a un usuario registrado.
- `message` es opcional, se recorta en backend y acepta hasta 250 caracteres.

#### Response `201 Created`

```json
{
  "id": 10,
  "event_id": 1,
  "event_title": "Nombre del evento",
  "image_url": "https://...",
  "event_start_date": "2026-06-25T20:00:00Z",
  "event_location": "Buenos Aires",
  "event_price": 15000,
  "status": "ACTIVE",
  "purchase_date": "2026-06-25T12:00:00Z",
  "user_id": 2,
  "user_email": "lorenzo@test.com",
  "is_gift": true,
  "gifted_by_id": 1,
  "gifted_by_email": "octavio@test.com",
  "gift_message": "Feliz cumple!",
  "gifted_at": "2026-06-25T12:00:00Z"
}
```

#### Errores

- `400`: body inválido, `event_id`/`target_email` vacío, mensaje demasiado largo o intento de regalarse a sí mismo.
- `401`: falta token o token inválido.
- `404`: evento o usuario destino inexistente.
- `409`: no hay capacidad disponible o el destinatario ya tiene una entrada `ACTIVE` para ese evento.
- `500`: error interno.


> Las respuestas de tickets no devuelven objetos GORM completos ni `password_hash`. `event_price` corresponde al precio actual del evento asociado y está expresado en ARS.


## Admin endpoints

Los endpoints de administración requieren el header `Authorization: Bearer <token>` con un JWT de un usuario con rol `ADMIN`. Para desarrollo local existe el usuario demo `admin@test.com` / `123456`.

> El usuario admin demo es solo para desarrollo local.

### `GET /api/admin/events`

Lista eventos para administración y mantiene los campos públicos de evento. Además incluye la cantidad de entradas vendidas y cupos disponibles.

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
    "available_capacity": 99,
    "tickets_sold": 1,
    "active": true
  }
]
```

- `tickets_sold` cuenta tickets `ACTIVE` del evento.
- `available_capacity` se calcula como `capacity - tickets_sold`; si da negativo, devuelve `0`.
- Los tickets `CANCELLED` no se cuentan como vendidos.

#### Errores

- `401 Unauthorized`: falta autenticación.
- `403 Forbidden`: el usuario autenticado no tiene rol `ADMIN`.
- `500 Internal Server Error`: error interno.

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
- `category` debe ser una de: `Música`, `Teatro`, `Deportes`, `Tecnología`, `Otros`.
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
- Si se envía `category`, debe ser una de: `Música`, `Teatro`, `Deportes`, `Tecnología`, `Otros`.
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

## Admin stats endpoints

Los endpoints administrativos de estadísticas requieren autenticación JWT con un usuario de rol `ADMIN`:

```http
Authorization: Bearer <token>
```

Usuario demo admin: `admin@test.com` / `123456`.

> `estimated_revenue` está expresado en ARS. Es un valor estimado porque se calcula con los tickets en estado `ACTIVE` y el precio actual del evento (`event.price`); esta PR no guarda precio histórico por ticket. Los tickets `CANCELLED` no cuentan para ingresos.

### `GET /api/admin/stats/summary`

Devuelve estadísticas generales del sistema.

#### Response `200 OK`

```json
{
  "total_users": 4,
  "client_users": 3,
  "admin_users": 1,
  "total_events": 5,
  "active_events": 4,
  "inactive_events": 1,
  "total_tickets": 10,
  "active_tickets": 8,
  "cancelled_tickets": 2,
  "total_capacity": 490,
  "available_capacity": 482,
  "occupancy_rate_percent": 1.63,
  "estimated_revenue": 120000
}
```

#### Reglas

- `total_capacity` suma la capacidad de eventos activos.
- `available_capacity` suma los cupos disponibles de eventos activos.
- `occupancy_rate_percent` se calcula como `active_tickets / total_capacity * 100`; si `total_capacity` es `0`, devuelve `0`.
- `estimated_revenue` suma `event.price` por cada ticket `ACTIVE`.

#### Errores

- `401 Unauthorized`: token faltante, inválido o expirado.
- `403 Forbidden`: usuario autenticado sin rol `ADMIN`.
- `500 Internal Server Error`: error interno.

### `GET /api/admin/stats/events`

Devuelve estadísticas por evento. Incluye eventos activos e inactivos y ordena por `event_id` ascendente.

#### Response `200 OK`

```json
[
  {
    "event_id": 1,
    "title": "Rock Nacional",
    "category": "Música",
    "location": "Córdoba",
    "active": true,
    "capacity": 100,
    "price": 15000,
    "active_tickets": 20,
    "cancelled_tickets": 3,
    "total_tickets": 23,
    "available_capacity": 80,
    "occupancy_rate_percent": 20,
    "estimated_revenue": 300000
  }
]
```

#### Reglas

- `active_tickets` cuenta tickets `ACTIVE` del evento.
- `cancelled_tickets` cuenta tickets `CANCELLED` del evento.
- `total_tickets` cuenta todos los tickets del evento.
- `available_capacity` se calcula como `capacity - active_tickets`; si da negativo, devuelve `0`.
- `occupancy_rate_percent` se calcula como `active_tickets / capacity * 100`; si `capacity` es `0`, devuelve `0`.
- `estimated_revenue` se calcula como `active_tickets * event.price`.

#### Errores

- `401 Unauthorized`: token faltante, inválido o expirado.
- `403 Forbidden`: usuario autenticado sin rol `ADMIN`.
- `500 Internal Server Error`: error interno.

### `GET /api/admin/events/:id/report`

Devuelve el reporte detallado de un evento específico, activo o inactivo.

#### Response `200 OK`

```json
{
  "event_id": 1,
  "title": "Rock Nacional",
  "description": "Evento de música rock nacional.",
  "category": "Música",
  "location": "Córdoba",
  "start_date": "2026-09-12T21:00:00Z",
  "duration_minutes": 120,
  "active": true,
  "capacity": 100,
  "price": 15000,
  "active_tickets": 20,
  "cancelled_tickets": 3,
  "total_tickets": 23,
  "available_capacity": 80,
  "occupancy_rate_percent": 20,
  "estimated_revenue": 300000
}
```

#### Errores

- `400 Bad Request`: id inválido.
- `401 Unauthorized`: token faltante, inválido o expirado.
- `403 Forbidden`: usuario autenticado sin rol `ADMIN`.
- `404 Not Found`: evento inexistente.
- `500 Internal Server Error`: error interno.

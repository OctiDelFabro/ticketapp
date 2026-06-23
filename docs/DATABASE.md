# Modelo de Base de Datos

## Entidades mínimas

### User

- `id`
- `name`
- `email`
- `password_hash`
- `role`
- `created_at`
- `updated_at`

### Event

- `id`
- `title`
- `description`
- `image_url`
- `category`
- `location`
- `start_date`
- `duration_minutes`
- `capacity`
- `price` (valor de entrada general en ARS, `not null`, default `0`)
- `active`
- `created_at`
- `updated_at`

### Ticket

- `id`
- `user_id`
- `event_id`
- `status`
- `purchase_date`
- `created_at`
- `updated_at`

## Relaciones

- Un usuario puede tener muchas entradas.
- Un evento puede tener muchas entradas.
- Una entrada pertenece a un usuario y a un evento.

# Frontend

## Vistas planificadas

- Home / Catálogo de eventos.
- Detalle del evento.
- Login.
- Registro.
- Mis Entradas.
- Admin Eventos, para etapa posterior.
- Formulario Evento, para etapa posterior.
- Reportes Admin, para etapa posterior.

## Lineamientos de interfaz

La interfaz debe ser responsive y simple, priorizando la claridad del flujo del cliente y la facilidad de uso.


## Precios de eventos

El frontend consume el campo `price` enviado por la API para eventos y lo muestra con formato visual en ARS, por ejemplo `$15.000`. Si el precio es `0`, `null` o `undefined`, se muestra `Gratis`.

El precio aparece en:

- Catálogo de eventos.
- Detalle del evento.
- Checkout, usando `event.price` para el precio unitario y el total.
- Mis Entradas, usando `event_price` del ticket o un fallback compatible si no está disponible.

## Categorías de eventos

Los formularios y filtros de eventos usan una lista fija de categorías: Música, Teatro, Deportes, Tecnología y Otros. El formulario admin de creación/edición muestra estas opciones en un selector para evitar categorías libres.

## Panel admin de eventos

El panel de administración de eventos muestra la foto del evento, entradas vendidas, cupos disponibles y capacidad total usando los campos `tickets_sold`, `available_capacity` y `capacity` de la API.

## Panel admin de reportes

El panel admin de reportes consume endpoints autenticados y requiere un usuario con rol `ADMIN`.

Endpoints usados:

- `GET /api/admin/stats/summary`: métricas generales de usuarios, eventos, tickets, capacidad, ocupación e ingresos estimados.
- `GET /api/admin/stats/events`: métricas agregadas por evento para listar reportes.
- `GET /api/admin/events/:id/report`: detalle de estadísticas de un evento específico.

Los ingresos mostrados en reportes son estimados y se calculan con tickets `ACTIVE` y el precio actual del evento. Los tickets `CANCELLED` no suman ingresos.

## Regalar entradas

Desde el detalle de un evento autenticado, el usuario puede elegir **Regalar entrada** para iniciar un checkout en modo regalo. En este modo siempre se procesa una sola entrada, se solicita el email del destinatario y se permite agregar un mensaje opcional de hasta 250 caracteres.

El destinatario debe ser un usuario registrado de TicketApp. Al confirmar, el frontend llama a `POST /api/tickets/gift`; la entrada queda asociada a la cuenta del destinatario y aparece en **Mis Entradas** con badge **Regalada**, el email de quien regaló y el mensaje si existe.

TicketApp no envía emails reales para los regalos y no integra pasarelas de pago reales en este flujo.

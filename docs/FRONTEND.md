# Frontend

El frontend de TicketApp está construido con React, Vite, Tailwind CSS y CSS propio. Consume la API REST del backend y separa la experiencia en vistas públicas, vistas autenticadas de cliente y vistas protegidas de administrador.

## Vistas Cliente

- **Home / Catálogo de eventos:** listado público de eventos activos con búsqueda, filtros y disponibilidad.
- **Detalle del evento:** muestra descripción, imagen, categoría, ubicación, fecha, precio y cupo disponible.
- **Login y Registro:** obtienen el JWT y guardan la sesión del usuario.
- **Checkout / Compra:** permite comprar una entrada para el usuario autenticado.
- **Mis Entradas:** lista tickets propios, estado, datos del evento, precio, acciones disponibles y datos de regalo si corresponde.
- **Cancelación:** acción desde Mis Entradas para cambiar un ticket propio a `CANCELLED`.
- **Transferencia:** acción para mover una entrada activa existente a otro usuario registrado.

## Vistas Admin

Las vistas administrativas requieren un usuario con rol `ADMIN` y consumen endpoints protegidos por JWT.

- **Admin Eventos:** panel de gestión con foto del evento, entradas vendidas, cupos disponibles y capacidad total usando `tickets_sold`, `available_capacity` y `capacity`.
- **Formulario Evento:** creación y edición de eventos con título, descripción, imagen, categoría, ubicación, fecha, duración, capacidad, precio y estado.
- **Reportes Admin:** métricas generales y por evento, incluyendo usuarios, eventos, tickets activos, tickets cancelados, capacidad, ocupación e ingresos estimados.

## Lineamientos de interfaz

La interfaz debe ser responsive y simple, priorizando la claridad del flujo del cliente y la facilidad de uso durante la defensa final.

## Precios de eventos

El frontend consume el campo `price` enviado por la API para eventos y lo muestra con formato visual en ARS, por ejemplo `$15.000`. Si el precio es `0`, `null` o `undefined`, se muestra `Gratis`.

El precio aparece en:

- Catálogo de eventos.
- Detalle del evento.
- Checkout, usando `event.price` para el precio unitario y el total.
- Mis Entradas, usando `event_price` del ticket o un fallback compatible si no está disponible.

## Categorías de eventos

Los formularios y filtros de eventos usan una lista fija de categorías: Música, Teatro, Deportes, Tecnología y Otros. El formulario admin de creación/edición muestra estas opciones en un selector para evitar categorías libres.

## Bonus: Regalar entradas

Desde el detalle de un evento autenticado, el usuario puede elegir **Regalar entrada** para iniciar un checkout en modo regalo. En este modo siempre se procesa una sola entrada, se solicita el email del destinatario y se permite agregar un mensaje opcional de hasta 250 caracteres.

El destinatario debe ser un usuario registrado de TicketApp. Al confirmar, el frontend llama a `POST /api/tickets/gift`; la entrada queda asociada a la cuenta del destinatario y aparece en **Mis Entradas** con badge **Regalada**, el email de quien regaló y el mensaje si existe.

La diferencia con transferencia es que regalar crea/emite una entrada nueva para otro usuario, mientras que transferir mueve una entrada existente.

TicketApp no envía emails reales para los regalos y no integra pasarelas de pago reales en este flujo.

## Docker

En Docker Compose, el frontend se construye con Vite y se sirve con Nginx. El servicio queda expuesto en `http://localhost:5173` y se comunica con el backend publicado en `http://localhost:8080`.

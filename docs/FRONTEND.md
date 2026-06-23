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

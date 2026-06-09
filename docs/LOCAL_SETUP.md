# Setup local de TicketApp

Esta guía explica cómo levantar TicketApp en una computadora local con datos demo reproducibles. No hace falta compartir bases de datos, dumps SQL ni archivos `.env` reales entre integrantes del equipo.

## Requisitos

- Go
- Node
- MySQL
- Git

## 1. Clonar el repositorio

```bash
git clone https://github.com/OctiDelFabro/ticketapp.git
cd ticketapp
```

## 2. Crear el archivo de entorno del backend

Copiá el ejemplo versionado y creá tu propio `.env` local:

```bash
cp backend/.env.example backend/.env
```

Luego ajustá `DB_USER` y `DB_PASSWORD` según tu instalación local de MySQL.

Ejemplo de variables esperadas:

```env
PORT=8080
DB_USER=root
DB_PASSWORD=root
DB_HOST=localhost
DB_PORT=3306
DB_NAME=ticketapp
JWT_SECRET=super-secret-dev-key
JWT_EXPIRATION_HOURS=24
```

> Las credenciales demo y el `JWT_SECRET` del ejemplo son solo para desarrollo local.

## 3. Correr el backend

```bash
cd backend
go mod tidy
go run main.go
```

Al iniciar, el backend:

- Crea la base `ticketapp` si no existe.
- Ejecuta las migraciones de GORM.
- Carga eventos demo.
- Carga usuarios demo.
- Carga tickets demo.
- Levanta la API en `http://localhost:8080`.

Los seeds son idempotentes: podés ejecutar `go run main.go` más de una vez sin duplicar usuarios, eventos ni tickets demo.

## Datos demo

Podés iniciar sesión con estos usuarios cliente:

| Nombre | Email | Password |
| --- | --- | --- |
| Lorenzo | `lorenzo@test.com` | `123456` |
| Pablo | `pablo@test.com` | `123456` |
| Octavio | `octavio@test.com` | `123456` |

Estas credenciales son solo para desarrollo local.

Los tickets demo incluyen:

- Octavio: ticket `ACTIVE` para Rock Nacional.
- Lorenzo: ticket `ACTIVE` para Festival Tech.
- Pablo: ticket `CANCELLED` para Stand Up Night.

## 4. Correr el frontend

En otra terminal, desde la raíz del repositorio:

```bash
cd frontend
npm install
npm run dev
```

Abrí la app en:

```txt
http://localhost:5173
```

## Troubleshooting

- Si no aparecen eventos, probá abrir `http://localhost:8080/api/events` y verificá que el backend esté respondiendo.
- Si no funciona el login, revisá que el backend esté corriendo en el puerto `8080`.
- Si falla MySQL, revisá `DB_USER` y `DB_PASSWORD` en `backend/.env`.
- Si no existe la base, el backend debería crearla automáticamente al iniciar.

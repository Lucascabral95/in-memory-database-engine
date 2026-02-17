# In-Memory Database Engine: E-Commerce API

API REST para gestion de e-commerce con arquitectura por capas, JWT, PostgreSQL y carrito en memoria con expiracion de 24 horas.

## Tabla de contenidos

- [Descripcion general](#descripcion-general)
- [Caracteristicas principales](#caracteristicas-principales)
- [Tecnologias utilizadas](#tecnologias-utilizadas)
- [Arquitectura del sistema](#arquitectura-del-sistema)
- [Flujo de checkout y pago con control de stock](#flujo-de-checkout-y-pago-con-control-de-stock)
- [Carrito en memoria (24h) y Redis TCP opcional](#carrito-en-memoria-24h-y-redis-tcp-opcional)
- [Estructura del proyecto](#estructura-del-proyecto)
- [Modelos de datos](#modelos-de-datos)
- [Catalogo de endpoints](#catalogo-de-endpoints)
- [Guia de instalacion y ejecucion](#guia-de-instalacion-y-ejecucion)
- [Guia Docker (build y run)](#guia-docker-build-y-run)
- [Guia de pruebas paso a paso con cURL](#guia-de-pruebas-paso-a-paso-con-curl)
- [CI: testing y build Docker](#ci-testing-y-build-docker)
- [Variables de entorno](#variables-de-entorno)
- [Scripts y comandos](#scripts-y-comandos)
- [Convenciones de Commits](#convenciones-de-commits)
- [Contribuciones](#contribuciones)
- [Licencia](#licencia)
- [Contacto](#contacto)

## Descripcion general

Este proyecto implementa una API de e-commerce con:

- Persistencia en PostgreSQL para usuarios, productos, ordenes y stock movements.
- Carrito en memoria RAM por usuario autenticado (JWT) con TTL de 24 horas.
- Checkout que crea orden y ejecuta pago con validacion transaccional de stock.
- Servidor TCP compatible con comandos basicos RESP (`PING`, `SET`, `GET`, `DEL`) para un "redis casero" opcional.

La aplicacion sigue el patron `Handler -> Service -> Repository -> Database`.

## Caracteristicas principales

- Arquitectura limpia y separacion de responsabilidades.
- Autenticacion JWT.
- CRUD de categorias y productos.
- Ordenes con estados `PENDING`, `PAID`, `SHIPPED`, `CANCELLED`.
- Pago de orden con control estricto de stock.
- Descuento de stock en transaccion y registro de `StockMovement` con razon `SALE`.
- Carrito en memoria con TTL 24h renovado en escrituras.
- Endpoints de carrito para agregar, actualizar, remover items, limpiar carrito y checkout.
- Ownership enforcement en ordenes: cada usuario solo puede ver/modificar/pagar sus ordenes.
- Redis TCP opcional compartiendo el mismo `MemoryStore`.

## 🛠️ Tecnologias utilizadas

- Go (Golang)
- Gin (HTTP framework)
- GORM (ORM)
- PostgreSQL
- JWT (autenticacion)
- Docker y Docker Compose
- GitHub Actions (CI)
- Redis protocol (RESP) server custom en TCP

## Arquitectura del sistema

Capas:

1. `internal/handler`: entrada HTTP y validaciones de request/response.
2. `internal/service`: reglas de negocio.
3. `internal/repository`: acceso a datos con GORM y transacciones.
4. `internal/database`: conexion y migraciones.
5. `internal/storage` + `internal/server`: memoria en RAM + protocolo RESP opcional.

Dependencias inyectadas desde `cmd/api/main.go`.

## Flujo de checkout y pago con control de stock

Flujo actual de `POST /cart/checkout`:

1. Lee carrito desde RAM del usuario autenticado.
2. Revalida productos y calcula total.
3. Crea orden `PENDING`.
4. Ejecuta pago (`OrderUpdatePay`) dentro de transaccion:
   - Bloquea filas de productos (`FOR UPDATE`).
   - Verifica stock por producto.
   - Descuenta stock.
   - Inserta `StockMovement` con `Reason=SALE` y `Quantity` negativa.
   - Cambia orden a `PAID`.
5. Si no hay stock, responde `409` e intenta cancelar la orden a `CANCELLED`.
6. Solo si paga correctamente, limpia el carrito en RAM.

## Carrito en memoria (24h) y Redis TCP opcional

### Carrito en memoria

- Clave por usuario: `cart:user:<user_id>`.
- TTL: `24h` (`86400` segundos).
- TTL se refresca en escrituras del carrito.
- Si una clave expira, se elimina al leer y tambien por cleanup periodico.

### Redis TCP opcional

Si activas:

- `REDIS_TCP_ENABLED=true`
- `REDIS_TCP_PORT=6379`

la app levanta un servidor TCP RESP en paralelo a la API HTTP.

Comandos soportados:

- `PING`
- `SET key value [EX seconds]`
- `GET key`
- `DEL key`

El servidor RESP usa el mismo `MemoryStore` en RAM que el carrito.

## Estructura del proyecto

```text
in-memory-database-engine/
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── database/
│   │   └── posgresql.go
│   ├── handler/
│   │   ├── cart.go
│   │   ├── category.go
│   │   ├── order.go
│   │   ├── product.go
│   │   ├── stockmovement.go
│   │   └── user.go
│   ├── model/
│   │   ├── base.go
│   │   ├── cart.go
│   │   ├── category.go
│   │   ├── order.go
│   │   ├── orderitem.go
│   │   ├── product.go
│   │   ├── stockmovement.go
│   │   └── user.go
│   ├── repository/
│   │   ├── category.go
│   │   ├── order.go
│   │   ├── product.go
│   │   ├── stockmovement.go
│   │   └── user.go
│   ├── routes/
│   │   └── routes.go
│   ├── server/
│   │   └── server.go
│   ├── service/
│   │   ├── category.go
│   │   ├── order.go
│   │   ├── product.go
│   │   ├── stockmovement.go
│   │   └── user.go
│   └── storage/
│       └── store.go
├── pkg/
│   ├── middleware/
│   │   └── auth.go
│   ├── response/
│   │   └── response.go
│   └── utils/
├── .env.template
├── go.mod
├── go.sum
└── README.md
```

## Modelos de datos

### Persistidos en PostgreSQL

- `User`
- `Category`
- `Product` (incluye `stock`)
- `Order`
- `OrderItem`
- `StockMovement` (`Reason`: `SALE`, `RESTOCK`, `ADJUSTMENT`)

### En memoria (no persistido en DB)

- `Cart`
- `CartItem`

## Catalogo de endpoints

Base URL local: `http://localhost:8080`

### Health

- `GET /health` (requiere JWT en el estado actual)

### Usuarios

- `POST /users/register`
- `POST /users/login`
- `DELETE /users/:email` (auth)

### Categorias

- `GET /categories`
- `GET /categories/:id`
- `POST /categories`
- `PATCH /categories/:id`
- `DELETE /categories/:id`

### Productos

- `GET /products`
- `GET /products/:id`
- `POST /products`
- `PATCH /products/:id`
- `DELETE /products/:id`

### Ordenes (auth)

- `GET /orders` (solo ordenes del usuario autenticado)
- `GET /orders/:id` (solo si es owner)
- `POST /orders`
- `POST /orders/:id/pay` (solo si es owner, con control de stock)
- `PATCH /orders/:id` (solo si es owner)

### Carrito (auth)

- `POST /cart/items`
- `PATCH /cart/items/:product_id`
- `DELETE /cart/items/:product_id`
- `GET /cart`
- `DELETE /cart`
- `POST /cart/checkout`

## Guia de instalacion y ejecucion

### Prerrequisitos

- Go 1.25+
- PostgreSQL (local o cloud)
- Git

### 1) Clonar repositorio

```bash
git clone <tu-repo-url>
cd in-memory-database-engine
```

### 2) Instalar dependencias

```bash
go mod download
```

### 3) Configurar entorno

```bash
cp .env.template .env
```

Completa `.env`:

```env
DATABASE_URL=postgresql://...
PORT=8080
ENV=development
JWT_SECRET=tu_clave_secreta
REDIS_TCP_ENABLED=false
REDIS_TCP_PORT=6379
```

### 4) Ejecutar API

```bash
go run cmd/api/main.go
```

Opcional en Windows:

```powershell
.\run.ps1
```

### 5) (Opcional) Activar Redis TCP

Setea en `.env`:

```env
REDIS_TCP_ENABLED=true
REDIS_TCP_PORT=6379
```

y vuelve a ejecutar la app.

## 🐳 Guia Docker (build y run)

### Opcion A: Docker Compose (recomendado)

1. Configura variables:

```bash
cp .env.template .env
```

2. Levanta servicios:

```bash
docker compose up --build
```

3. Verifica que la API responde:

```bash
curl http://localhost:8080/health
```

4. Detener servicios:

```bash
docker compose down
```

### Opcion B: Dockerfile (solo API)

1. Construir imagen:

```bash
docker build -t in-memory-db-api:latest .
```

2. Ejecutar contenedor:

```bash
docker run --rm -p 8080:8080 \
  -e DATABASE_URL="postgresql://usuario:password@host:5432/dbname?sslmode=disable" \
  -e PORT="8080" \
  -e ENV="production" \
  -e JWT_SECRET="tu_jwt_secret" \
  -e REDIS_TCP_ENABLED="false" \
  -e REDIS_TCP_PORT="6379" \
  in-memory-db-api:latest
```

3. Probar endpoint:

```bash
curl http://localhost:8080/health
```

## Guia de pruebas paso a paso con cURL

### Paso 0: definir base

```bash
BASE=http://localhost:8080
```

### Paso 1: registrar usuario

```bash
curl -X POST $BASE/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "email":"cart.test@example.com",
    "password":"password123",
    "first_name":"Cart",
    "last_name":"Tester"
  }'
```

### Paso 2: login y guardar token

```bash
curl -X POST $BASE/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "email":"cart.test@example.com",
    "password":"password123"
  }'
```

Guarda `token` de la respuesta en `TOKEN`.

### Paso 3: crear producto

```bash
curl -X POST $BASE/products \
  -H "Content-Type: application/json" \
  -d '{
    "name":"Mouse Gamer",
    "description":"RGB",
    "price":99.9,
    "stock":5
  }'
```

Guarda `id` del producto en `PRODUCT_ID`.

### Paso 4: agregar item al carrito

```bash
curl -X POST $BASE/cart/items \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"product_id\":\"$PRODUCT_ID\",
    \"quantity\":2
  }"
```

### Paso 5: actualizar cantidad del item

```bash
curl -X PATCH $BASE/cart/items/$PRODUCT_ID \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"quantity":3}'
```

### Paso 6: ver carrito

```bash
curl $BASE/cart \
  -H "Authorization: Bearer $TOKEN"
```

### Paso 7: checkout (crea orden + intenta pagar)

```bash
curl -X POST $BASE/cart/checkout \
  -H "Authorization: Bearer $TOKEN"
```

Ejemplo de exito:

```json
{
  "status": "success",
  "order_id": "uuid",
  "order_status": "PAID",
  "message": "Checkout completado y orden pagada",
  "order": { "...": "..." }
}
```

Ejemplo de fallo por stock:

```json
{
  "status": "failed",
  "order_id": "uuid",
  "order_status": "CANCELLED",
  "error": "Stock insuficiente para completar el checkout"
}
```

### Paso 8: listar ordenes del usuario

```bash
curl $BASE/orders \
  -H "Authorization: Bearer $TOKEN"
```

### Paso 9: pago manual de una orden puntual

```bash
curl -X POST $BASE/orders/<ORDER_ID>/pay \
  -H "Authorization: Bearer $TOKEN"
```

### Paso 10: remover item del carrito

```bash
curl -X DELETE $BASE/cart/items/$PRODUCT_ID \
  -H "Authorization: Bearer $TOKEN"
```

### Paso 11: limpiar carrito completo

```bash
curl -X DELETE $BASE/cart \
  -H "Authorization: Bearer $TOKEN"
```

### Paso 12: probar Redis TCP (opcional)

Con app corriendo y `REDIS_TCP_ENABLED=true`:

```bash
redis-cli -p 6379 PING
redis-cli -p 6379 SET demo hola EX 30
redis-cli -p 6379 GET demo
redis-cli -p 6379 DEL demo
```

## 🚦 CI: testing y build Docker

El workflow de CI esta en `.github/workflows/ci.yml` y ejecuta tres checks:

1. `CI / Go Build & Test`
2. `CI / Docker Image Build`
3. `CI / Required Checks`

### Que valida el CI

- Descarga dependencias (`go mod download`).
- Verifica `go.mod`/`go.sum` limpios (`go mod tidy` + diff).
- Ejecuta analisis estatico (`go vet ./...`).
- Ejecuta tests (`go test ./...` con cobertura).
- Compila la API (`go build ./cmd/api/main.go`).
- Construye imagen Docker (`docker build -f Dockerfile .`).

### Variables de entorno en GitHub Actions

El pipeline lee secretos del repositorio:

- `DATABASE_URL`
- `PORT`
- `ENV`
- `JWT_SECRET`
- `REDIS_TCP_ENABLED`
- `REDIS_TCP_PORT`

### Reproducir validacion en local

```bash
go mod download
go mod tidy
go vet ./...
go test ./... -v -covermode=atomic -coverprofile=coverage.out
go build -o bin/api-ci ./cmd/api/main.go
docker build --file Dockerfile --tag in-memory-db-api:ci .
```

### Bloquear merge si falla CI

Configura en GitHub (`Settings -> Branches` o `Rulesets`) para `main`:

- Activar `Require status checks to pass before merging`.
- Agregar como requerido: `CI / Required Checks`.
- Si no quieres aprobacion manual, dejar `Required approvals = 0`.

## ⚙️ Variables de entorno

| Variable | Descripcion | Requerido | Default |
|---|---|---|---|
| `DATABASE_URL` | DSN de PostgreSQL | Si | - |
| `PORT` | Puerto HTTP | No | `8080` |
| `ENV` | Entorno (`development`/`production`) | No | `development` |
| `JWT_SECRET` | Secreto JWT | No | `sd,fsdnlfksdmlkf` |
| `REDIS_TCP_ENABLED` | Habilita servidor RESP TCP | No | `false` |
| `REDIS_TCP_PORT` | Puerto TCP RESP | No | `6379` |

## 🧪 Scripts y comandos

```bash
# ejecutar app
go run cmd/api/main.go

# compilar
go build -o api.exe cmd/api/main.go

# dependencias
go mod download
go mod tidy

# tests
go test ./...

# seed de datos
go run cmd/seed/main.go
go run cmd/seed/main.go --force

# analisis estatico
go vet ./...
```

## Contribuciones

Contribuciones bienvenidas.

Recomendado:

1. Fork del repositorio.
2. Rama por feature/fix.
3. Cambios + pruebas.
4. Pull Request con descripcion tecnica clara.

### Convenciones de Commits

Este proyecto sigue [Conventional Commits](https://www.conventionalcommits.org/):

- `feat:` Nueva funcionalidad
- `fix:` Correccion de bugs
- `docs:` Cambios en documentacion
- `style:` Cambios de formato (no afectan la logica)
- `refactor:` Refactorizacion de codigo
- `test:` Anadir o modificar tests
- `chore:` Tareas de mantenimiento

---

## Licencia

Este proyecto esta bajo la licencia **MIT**.

---

<a id="contact-anchor"></a>
## 📬 Contacto

- **Autor:** Lucas Cabral
- **Email:** lucassimple@hotmail.com
- **LinkedIn:** [https://www.linkedin.com/in/lucas-gastón-cabral/](https://www.linkedin.com/in/lucas-gastón-cabral/)
- **Portfolio:** [https://portfolio-web-dev-git-main-lucascabral95s-projects.vercel.app/](https://portfolio-web-dev-git-main-lucascabral95s-projects.vercel.app/)
- **Github:** [https://github.com/Lucascabral95](https://github.com/Lucascabral95/)

---

<p align="center">
  Desarrollado con ❤️ por Lucas Cabral
</p>

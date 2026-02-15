<h1 align="center">In-Memory Database Engine: E-Commerce API</h1>

<p align="center">
  API REST robusta para gestión de inventario y órdenes de e-commerce con arquitectura limpia, autenticación JWT y PostgreSQL.
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.25+-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go Version"/>
  <img src="https://img.shields.io/badge/Gin-Framework-6BA338?style=for-the-badge&logo=gin&logoColor=white" alt="Gin Framework"/>
  <img src="https://img.shields.io/badge/PostgreSQL-Neon-4169E1?style=for-the-badge&logo=postgresql&logoColor=white" alt="PostgreSQL"/>
  <img src="https://img.shields.io/badge/JWT-Authentication-FF6600?style=for-the-badge&logo=JSON%20web%20tokens&logoColor=white" alt="JWT"/>
</p>

***

## Table of contents

- [Descripción general](#descripción-general)
- [⚙️ Características principales](#️-características-principales)
- [🏛️ Arquitectura del sistema](#️-arquitectura-del-sistema)
  - [Flujo de datos](#flujo-de-datos)
  - [Diagrama de secuencia](#diagrama-de-secuencia)
- [Estructura del proyecto](#estructura-del-proyecto)
- [📦 Modelos de datos](#-modelos-de-datos)
  - [Relaciones entre entidades](#relaciones-entre-entidades)
- [🛠️ Catálogo de endpoints](#️-catálogo-de-endpoints)
  - [🔐 Autenticación](#-autenticación)
  - [📦 Categorías](#-categorías)
  - [🎯 Productos](#-productos)
  - [🛒 Órdenes](#️-órdenes)
- [🚀 Guía de instalación y ejecución](#-guía-de-instalación-y-ejecución)
- [🧪 Guía de pruebas con cURL](#-guía-de-pruebas-con-curl)
- [🔒 Variables de entorno](#-variables-de-entorno)
- [🛠️ Scripts y comandos](#️-scripts-y-comandos)
- [Contribuciones](#contribuciones)
  - [Convenciones de Commits](#convenciones-de-commits)
- [Licencia](#licencia)
- [📬 Contacto](#-contacto)

## Descripción general

**In-Memory DB** es una API REST completa para la gestión de inventario y procesamiento de órdenes de un e-commerce, construida con **Go** y el framework **Gin**. A pesar del nombre del repositorio, la aplicación utiliza **PostgreSQL** como base de datos persistente con **GORM** como ORM, implementando una arquitectura limpia y escalable.

La solución proporciona una separación clara de responsabilidades mediante el patrón **Repository-Service-Handler**, con autenticación basada en **JWT**, gestión de stock, y un sistema de órdenes con múltiples estados (PENDING, PAID, SHIPPED, CANCELLED).

El sistema está diseñado para ser fácilmente extensible, mantenible y testeable, siguiendo los principios de **Clean Architecture** y **Dependency Injection**.

***

<a id="️-características-principales"></a>
## ⚙️ Características principales

- **Arquitectura limpia**: Separación clara en capas (Handler → Service → Repository → Database) con Dependency Injection
- **Autenticación JWT**: Tokens seguros con expiración de 30 días para proteger endpoints sensibles
- **Gestión completa del inventario**: CRUD de categorías y productos con control de stock
- **Sistema de órdenes**: Flujo completo desde PENDING hasta SHIPPED con múltiples estados
- **Autenticación de usuarios**: Registro, login con bcrypt y eliminación de cuentas
- **Soft deletes**: Eliminación lógica de registros para auditoría y recuperación
- **UUID v4**: Identificadores únicos y seguros para todas las entidades
- **Migraciones automáticas**: Sincronización automática del esquema con GORM
- **Validación de datos**: Validación robusta en handlers y servicios
- **Logging configurable**: Diferentes niveles de log según entorno (development/production)
- **Código modular**: Paquetes reutilizables (JWT, password hashing, middleware)
- **Base de datos serverless**: Compatible con PostgreSQL en la nube (Neon Tech)

***

<a id="️-arquitectura-del-sistema"></a>
## 🏛️ Arquitectura del sistema

<p align="center">
  <img src="https://go.dev/blog/go-brand/Go-Logo/PNG/Go-Logo_Aqua.png"
       alt="Go Gopher"
       width="120"/>
</p>

### Patrón de arquitectura por capas

```
┌─────────────────────────────────────────────────────────────┐
│                         CLIENT (HTTP)                        │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                     GIN ROUTER (LAYER)                       │
│  /categories  │  /users  │  /products  │  /orders           │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                   MIDDLEWARE (AUTH JWT)                      │
│         Validación de token y extracción de claims          │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                      HANDLER LAYER                           │
│  CategoryHandler │ UserHandler │ ProductHandler │ OrderHandler│
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                       SERVICE LAYER                          │
│   Lógica de negocio │ Validaciones │ Reglas de dominio      │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                    REPOSITORY LAYER                          │
│     AbstractDataAccess │ Interfaces │ Implementaciones GORM │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                     POSTGRESQL (DB)                          │
│           Tablas: users, categories, products, orders        │
└─────────────────────────────────────────────────────────────┘
```

### Flujo de Dependency Injection

```go
// En cmd/api/main.go - Cadena de inyección de dependencias

db := database.InitDB(cfg)
         │
         ▼
categoryRepo := repository.NewCategoryRepository(db)
categoryService := service.NewCategoryService(categoryRepo)
categoryHandler := handler.NewCategoryHandler(categoryService)
         │
         ▼
routes.SetupRoutes(r, cfg, categoryService, ...)
```

<a id="flujo-de-datos"></a>
## Flujo de datos

1. **Autenticación**: El cliente se registra o inicia sesión → Recibe token JWT
2. **Request autenticada**: Cliente envía request con `Authorization: Bearer <token>`
3. **Middleware**: Valida token y extrae `userID` del contexto
4. **Handler**: Procesa request, valida datos y llama al service
5. **Service**: Ejecuta lógica de negocio y llama al repository
6. **Repository**: Ejecuta operaciones en base de datos con GORM
7. **Response**: Datos retornados en formato JSON con preloads de relaciones

<a id="diagrama-de-secuencia"></a>
## Diagrama de secuencia

```
Cliente ──────► Router ──────► Middleware ──────► Handler ──────► Service ──────► Repository ──────► Database
   │              │                │                 │                │                  │                  │
   │──POST /orders│                │                 │                │                  │                  │
   │              │                │                 │                │                  │                  │
   │              │──Validate JWT──│                 │                │                  │                  │
   │              │                │──userID exists──│                │                  │                  │
   │              │                │                 │──CreateOrder───│                  │                  │
   │              │                │                 │                │──CreateOrder─────│                  │
   │              │                │                 │                │                  │──INSERT ORDER───│
   │              │                │                 │                │                  │──PRELOAD User───│
   │              │                │                 │                │──RETURN order────│                  │
   │              │                │                 │──201 JSON──────│                  │                  │
   │◄─201 + order─│◄───────────────│◄────────────────│                │                  │                  │
```

## Estructura del proyecto

```
in-memory-database-engine/
├── cmd/
│   └── api/
│       └── main.go              # Punto de entrada de la aplicación
├── internal/
│   ├── config/
│   │   └── config.go           # Carga de variables de entorno
│   ├── database/
│   │   └── posgresql.go        # Conexión a PostgreSQL y migraciones
│   ├── handler/
│   │   ├── category.go         # HTTP Handlers para categorías
│   │   ├── order.go            # HTTP Handlers para órdenes
│   │   ├── product.go          # HTTP Handlers para productos
│   │   ├── stockmovement.go    # HTTP Handlers para movimientos de stock
│   │   └── user.go             # HTTP Handlers para usuarios
│   ├── model/
│   │   ├── base.go             # BaseModel (ID, CreatedAt, UpdatedAt, DeletedAt)
│   │   ├── category.go         # Entity: Category
│   │   ├── order.go            # Entities: Order, OrderStatus
│   │   ├── orderitem.go        # Entity: OrderItem
│   │   ├── product.go          # Entity: Product
│   │   ├── stockmovement.go    # Entity: StockMovement
│   │   └── user.go             # Entity: User
│   ├── repository/
│   │   ├── category.go         # Repository: Category
│   │   ├── order.go            # Repository: Order
│   │   ├── product.go          # Repository: Product
│   │   ├── stockmovement.go    # Repository: StockMovement
│   │   └── user.go             # Repository: User
│   ├── routes/
│   │   └── routes.go           # Configuración de rutas y middlewares
│   └── service/
│       ├── category.go         # Service: Category
│       ├── order.go            # Service: Order
│       ├── product.go          # Service: Product
│       ├── stockmovement.go    # Service: StockMovement
│       └── user.go             # Service: User
├── pkg/
│   ├── middleware/
│   │   └── auth.go             # Middleware de autenticación JWT
│   ├── response/
│   │   └── response.go         # Utilidades de respuesta HTTP
│   └── utils/
│       ├── convert-strings.go  # Conversión de tipos
│       ├── generate-sku.go     # Generador de SKU único
│       ├── jwt.go              # Generación y validación de JWT
│       └── utils.go            # Hash de contraseñas (bcrypt)
├── .env                         # Variables de entorno
├── .gitignore                   # Archivos ignorados por Git
├── go.mod                       # Módulos y dependencias de Go
├── go.sum                       # Checksums de dependencias
├── run.ps1                      # Script de ejecución (PowerShell)
└── README.md                    # Esta documentación
```

<a id="-modelos-de-datos"></a>
## 📦 Modelos de datos

### BaseModel (Herencia automática)

Todos los modelos extienden de `BaseModel`:

```go
type BaseModel struct {
    ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
```

### Entidades principales

| Modelo | Descripción | Campos principales |
|--------|-------------|-------------------|
| **User** | Usuario del sistema | email, password (bcrypt), first_name, last_name |
| **Category** | Categoría de productos | name (unique) |
| **Product** | Producto del inventario | name, sku (unique), price, stock, category_id |
| **Order** | Orden de compra | user_id, total_amount, status (PENDING/PAID/SHIPPED/CANCELLED) |
| **OrderItem** | Item de una orden | order_id, product_id, quantity, price_at_moment |
| **StockMovement** | Movimiento de stock | product_id, quantity, reason (SALE/RESTOCK/ADJUSTMENT) |

<a id="relaciones-entre-entidades"></a>
## Relaciones entre entidades

```
User (1) ─────< (N) Order (1) ─────< (N) OrderItem (N) >───── (1) Product
 │                                                                   │
 │                                                                   │
 └───────────────────── (N) >───── (1) Category <───────────────────┘
                                                          │
                                                          │
                                                          ▼
                                                 StockMovement (N)
```

**Detalles de relaciones:**

- **User → Order**: Un usuario puede tener múltiples órdenes
- **Order → OrderItem**: Una orden contiene múltiples items
- **OrderItem → Product**: Un item pertenece a un producto (snapshot del precio)
- **Product → Category**: Un producto pertenece a una categoría (opcional)
- **Product → StockMovement**: Un producto tiene múltiples movimientos de stock

<a id="️-catálogo-de-endpoints"></a>
## 🛠️ Catálogo de endpoints

<a id="-autenticación"></a>
### 🔐 Autenticación

| Método | Endpoint | Descripción | Auth |
|--------|----------|-------------|------|
| `POST` | `/users/register` | Registrar nuevo usuario | ❌ |
| `POST` | `/users/login` | Iniciar sesión (devuelve JWT) | ❌ |
| `DELETE` | `/users/:email` | Eliminar usuario | ✅ |

**Request - Register:**
```json
POST /users/register
{
  "email": "usuario@ejemplo.com",
  "password": "password123",
  "first_name": "Juan",
  "last_name": "Pérez"
}
```

**Response - Login:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "message": "Login exitoso"
}
```

<a id="-categorías"></a>
### 📦 Categorías

| Método | Endpoint | Descripción | Auth |
|--------|----------|-------------|------|
| `GET` | `/categories` | Listar todas las categorías | ❌ |
| `GET` | `/categories/:id` | Obtener categoría por ID | ❌ |
| `POST` | `/categories` | Crear nueva categoría | ❌ |
| `PATCH` | `/categories/:id` | Actualizar categoría | ❌ |
| `DELETE` | `/categories/:id` | Eliminar categoría | ❌ |

<a id="-productos"></a>
### 🎯 Productos

| Método | Endpoint | Descripción | Auth |
|--------|----------|-------------|------|
| `GET` | `/products` | Listar todos los productos | ❌ |
| `GET` | `/products/:id` | Obtener producto por ID | ❌ |
| `POST` | `/products` | Crear nuevo producto | ❌ |
| `PATCH` | `/products/:id` | Actualizar producto | ❌ |
| `DELETE` | `/products/:id` | Eliminar producto | ❌ |

<a id="-órdenes"></a>
### 🛒 Órdenes

| Método | Endpoint | Descripción | Auth |
|--------|----------|-------------|------|
| `GET` | `/orders` | Listar todas las órdenes | ✅ |
| `GET` | `/orders/:id` | Obtener orden por ID | ✅ |
| `POST` | `/orders` | Crear nueva orden | ✅ |
| `PATCH` | `/orders/:id` | Actualizar orden | ✅ |
| `DELETE` | `/orders/:id` | Eliminar orden | ✅ |

**Request - Create Order:**
```json
POST /orders
Authorization: Bearer <token_jwt>

{
  "total_amount": 150.00,
  "items": [
    {
      "product_id": "uuid-producto-1",
      "quantity": 2,
      "price_at_moment": 50.00
    },
    {
      "product_id": "uuid-producto-2",
      "quantity": 1,
      "price_at_moment": 50.00
    }
  ]
}
```

**Nota:** El `user_id` se obtiene automáticamente del token JWT, no debe enviarse en el body.

**Response - Create Order:**
```json
{
  "id": "4d21d1e1-fc65-4529-90b4-0ab75ea10e6d",
  "created_at": "2026-02-15T14:44:24.744776-03:00",
  "updated_at": "2026-02-15T14:44:24.744776-03:00",
  "user_id": "29724283-95ed-49cc-a119-5944cbb96d01",
  "user": {
    "id": "29724283-95ed-49cc-a119-5944cbb96d01",
    "email": "usuario@ejemplo.com",
    "first_name": "Juan",
    "last_name": "Pérez"
  },
  "total_amount": 150.00,
  "status": "PENDING",
  "items": [
    {
      "id": "...",
      "product_id": "uuid-producto-1",
      "quantity": 2,
      "price_at_moment": 50.00,
      "product": {
        "id": "uuid-producto-1",
        "name": "Laptop Gaming",
        "sku": "LP-GAMING-001",
        "price": 1200.00,
        "stock": 10
      }
    }
  ]
}
```

<a id="-guía-de-instalación-y-ejecución"></a>
## 🚀 Guía de instalación y ejecución

### Prerrequisitos

- **Go** 1.25+ ([Instalar Go](https://golang.org/dl/))
- **PostgreSQL** 14+ o base de datos en la nube (Neon, Supabase, etc.)
- **Git** para clonar el repositorio

### 1. Clonar el repositorio

```bash
git clone https://github.com/tu-usuario/in-memory-database-engine.git
cd in-memory-database-engine
```

### 2. Instalar dependencias

```bash
go mod download
```

### 3. Configurar variables de entorno

Copiar el template de configuración y completar los valores:

```bash
# Copiar el template
cp .env.template .env

# Editar el archivo .env con tus valores reales
# Reemplaza: <tu_url_de_postgresql_aqui> y <cambia_esto_por_una_clave_secreta_muy_segura>
```

**Variables requeridas en `.env`:**

```env
DATABASE_URL=<tu_url_de_postgresql_aqui>
PORT=8080
ENV=development
JWT_SECRET=<cambia_esto_por_una_clave_secreta_muy_segura>
```

### 4. Ejecutar la aplicación

#### Opción A: Con `go run` (recomendado para desarrollo)

```bash
go run cmd/api/main.go
```

#### Opción B: Script PowerShell (evita bloqueo de Windows Defender)

```powershell
.\run.ps1
```

#### Opción C: Compilar y ejecutar

```bash
go build -o api.exe cmd/api/main.go
.\api.exe
```

### 5. Verificar que funciona

```bash
curl http://localhost:8080/health
```

**Respuesta esperada:**
```json
{
  "status": "ok",
  "message": "Server is running",
  "environment": "development",
  "port": "8080"
}
```

<a id="-guía-de-pruebas-con-curl"></a>
## 🧪 Guía de pruebas con cURL

### 1. Registrar un usuario

```bash
curl -X POST http://localhost:8080/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123",
    "first_name": "Test",
    "last_name": "User"
  }'
```

### 2. Iniciar sesión (obtener token)

```bash
curl -X POST http://localhost:8080/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

**Guardar el token** de la respuesta para los siguientes pasos.

### 3. Crear una categoría

```bash
curl -X POST http://localhost:8080/categories \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Electrónica"
  }'
```

### 4. Crear un producto

```bash
curl -X POST http://localhost:8080/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Laptop Gaming",
    "sku": "LP-GAMING-001",
    "description": "Laptop de alto rendimiento",
    "price": 1200.00,
    "stock": 10,
    "category_id": "uuid-de-la-categoria"
  }'
```

### 5. Crear una orden (requiere token)

```bash
curl -X POST http://localhost:8080/orders \
  -H "Authorization: Bearer TU_TOKEN_AQUI" \
  -H "Content-Type: application/json" \
  -d '{
    "total_amount": 2400.00,
    "items": [
      {
        "product_id": "uuid-del-producto",
        "quantity": 2,
        "price_at_moment": 1200.00
      }
    ]
  }'
```

### 6. Listar todas las órdenes

```bash
curl http://localhost:8080/orders \
  -H "Authorization: Bearer TU_TOKEN_AQUI"
```

<a id="-variables-de-entorno"></a>
## 🔒 Variables de entorno

| Variable | Descripción | Requerido | Default |
|----------|-------------|-----------|---------|
| `DATABASE_URL` | URL de conexión a PostgreSQL | ✅ | - |
| `PORT` | Puerto donde escucha el servidor | ❌ | `8080` |
| `ENV` | Entorno (development/production) | ❌ | `development` |
| `JWT_SECRET` | Clave secreta para firmar tokens JWT | ❌ | `sd,fsdnlfksdmlkf` |

**Ejemplo de `.env`:**

```env
# Base de datos PostgreSQL (Neon, Supabase, etc.)
DATABASE_URL=postgresql://neondb_owner:password@ep-host.aws.neon.tech/neondb?sslmode=require

# Configuración del servidor
PORT=8080
ENV=development

# Seguridad JWT
JWT_SECRET=mi_clave_muy_segura_cambiar_en_produccion
```

<a id="-scripts-y-comandos"></a>
## 🛠️ Scripts y comandos

### Comandos básicos de Go

```bash
# Ejecutar la aplicación
go run cmd/api/main.go

# Compilar binario
go build -o api.exe cmd/api/main.go

# Descargar dependencias
go mod download

# Actualizar dependencias
go mod tidy

# Ejecutar tests (cuando estén implementados)
go test ./...

# Formatear código
go fmt ./...

# Verificar código
go vet ./...
```

### Script PowerShell (Windows)

El archivo `run.ps1` evita el bloqueo de Windows Defender:

```powershell
# Establece directorio de compilación temporal
$env:GOTMPDIR = "$PWD\.go-build"

# Ejecuta la aplicación
go run cmd/api/main.go
```

## Contribuciones

¡Las contribuciones son bienvenidas! Seguí estos pasos:

1. Hacé un fork del repositorio.
2. Creá una rama para tu feature o fix (`git checkout -b feature/nueva-funcionalidad`).
3. Realizá tus cambios y escribí pruebas si es necesario.
4. Hacé commit y push a tu rama.
5. Abrí un Pull Request describiendo tus cambios.

### Convenciones de Commits

Este proyecto sigue [Conventional Commits](https://www.conventionalcommits.org/):

- `feat:` Nueva funcionalidad
- `fix:` Corrección de bugs
- `docs:` Cambios en documentación
- `style:` Cambios de formato (no afectan la lógica)
- `refactor:` Refactorización de código
- `test:` Añadir o modificar tests
- `chore:` Tareas de mantenimiento

## Licencia

Este proyecto está bajo la licencia **MIT**.

---

<a id="-contacto"></a>
## 📬 Contacto

- **Autor:** Lucas
- **Email:** lucassimple@hotmail.com
- **LinkedIn:** [Lucas Gastón Cabral](https://www.linkedin.com/in/lucas-gastón-cabral/)
- **Portfolio:** [Ver Portfolio](https://portfolio-web-dev-git-main-lucascabral95s-projects.vercel.app/)
- **GitHub:** [@Lucascabral95](https://github.com/Lucascabral95)

---

<p align="center">
  <img src="https://go.dev/blog/go-brand/Go-Logo/PNG/Go-Logo_Aqua.png" width="40"/>
  <br>
  Desarrollado con <span style="color: #e74c3c;">❤️</span> usando Go, Gin y GORM
</p>

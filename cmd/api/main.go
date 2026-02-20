package main

import (
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/lucas-dev/in-memory-db/docs"
	"github.com/lucas-dev/in-memory-db/internal/config"
	"github.com/lucas-dev/in-memory-db/internal/database"
	"github.com/lucas-dev/in-memory-db/internal/repository"
	"github.com/lucas-dev/in-memory-db/internal/routes"
	"github.com/lucas-dev/in-memory-db/internal/server"
	"github.com/lucas-dev/in-memory-db/internal/service"
	"github.com/lucas-dev/in-memory-db/internal/storage"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func registerSwaggerRoute(r *gin.Engine) {
	// Swagger UI: http://localhost:<PORT>/swagger/index.html
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func serverAddr(port string) string {
	return ":" + port
}

// @title In-Memory Database Engine API
// @version 1.0.0
// @description API REST de e-commerce con arquitectura por capas, JWT, PostgreSQL, carrito en memoria (TTL 24h), checkout atomico y control transaccional de stock.
// @contact.name Lucas Cabral
// @contact.url https://github.com/Lucascabral95
// @contact.email lucassimple@hotmail.com
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:8080
// @BasePath /
// @schemes http https
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Formato: Bearer <token>
// @tag.name users
// @tag.description Gestion de usuarios y autenticacion JWT.
// @tag.name products
// @tag.description CRUD de productos con control de stock.
// @tag.name categories
// @tag.description Gestion de categorias de productos.
// @tag.name orders
// @tag.description Ordenes de compra con estados y ownership enforcement.
// @tag.name cart
// @tag.description Carrito en memoria RAM con TTL de 24h y renovacion automatica.
// @tag.name stock-movements
// @tag.description Registro de movimientos de inventario (SALE, RESTOCK, ADJUSTMENT).
func main() {
	cfg := config.LoadConfig()
	log.Printf("=> Environment: %s", cfg.Environment)
	log.Printf("=> Server port: %s", cfg.Port)

	db := database.InitDB(cfg)

	defer func() {
		log.Println("=> Closing database connection...")
		if err := database.CloseDB(db); err != nil {
			log.Printf("Warning: failed to close database: %v", err)
		}
	}()

	// Category
	categoryRepo := repository.NewCategoryRepository(db)
	categoryService := service.NewCategoryService(categoryRepo)

	// User
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo, cfg.JWTSecret)

	// Product
	productRepo := repository.NewProductRepository(db)
	productService := service.NewProductService(productRepo)

	// Order
	orderRepo := repository.NewOrderRepository(db)
	orderService := service.NewOrderService(orderRepo)

	memoryStore := storage.NewMemoryStore()

	if cfg.RedisTCPEnabled {
		redisTCPServer := server.NewServer(memoryStore, cfg.GetRedisTCPAddr())
		go func() {
			log.Printf("=> Redis TCP enabled on %s", cfg.GetRedisTCPAddr())
			if err := redisTCPServer.Start(); err != nil {
				log.Printf("Warning: Redis TCP server stopped: %v", err)
			}
		}()
	}

	r := gin.Default()

	registerSwaggerRoute(r)

	routes.SetupRoutes(r, cfg, categoryService, userService, productService, orderService, memoryStore)

	if err := r.Run(serverAddr(cfg.Port)); err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
	log.Println("=> Server started on port " + cfg.Port)
}

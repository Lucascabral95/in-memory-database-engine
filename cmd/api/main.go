package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/lucas-dev/in-memory-db/internal/config"
	"github.com/lucas-dev/in-memory-db/internal/database"
	"github.com/lucas-dev/in-memory-db/internal/repository"
	"github.com/lucas-dev/in-memory-db/internal/routes"
	"github.com/lucas-dev/in-memory-db/internal/service"
)

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

	r := gin.Default()

	routes.SetupRoutes(r, cfg, categoryService, userService, productService, orderService)

	r.Run(":" + cfg.Port)
	log.Println("=> Server started on port " + cfg.Port)
}

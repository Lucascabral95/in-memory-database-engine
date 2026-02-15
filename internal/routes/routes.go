package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lucas-dev/in-memory-db/internal/config"
	"github.com/lucas-dev/in-memory-db/internal/handler"
	"github.com/lucas-dev/in-memory-db/internal/service"
	"github.com/lucas-dev/in-memory-db/pkg/middleware"
)

func SetupRoutes(
	r *gin.Engine,
	cfg *config.Config,
	categoryService *service.CategoryService,
	userService *service.UserService,
	productService *service.ProductService,
	orderService *service.OrderService,
) {
	categoryHandler := handler.NewCategoryHandler(categoryService)
	userHandler := handler.NewUserHandler(userService)
	productHandler := handler.NewProductHandler(productService)
	orderHandler := handler.NewOrderHandler(orderService)

	r.GET("/health", middleware.AuthMiddleware(cfg.JWTSecret), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":      "ok",
			"message":     "Server is running",
			"environment": cfg.Environment,
			"port":        cfg.Port,
		})
	})

	categories := r.Group("/categories")
	{
		categories.GET("", categoryHandler.FindAllCategories)
		categories.GET("/:id", categoryHandler.FindCategoryByID)
		categories.POST("", categoryHandler.CreateCategory)
		categories.PATCH("/:id", categoryHandler.UpdateCategory)
		categories.DELETE("/:id", categoryHandler.DeleteCategory)
	}

	users := r.Group("/users")
	{
		users.POST("/register", userHandler.RegisterUser)
		users.POST("/login", userHandler.LoginUser)

		users.DELETE("/:email", middleware.AuthMiddleware(cfg.JWTSecret), userHandler.DeleteUser)
	}

	products := r.Group("/products")
	{
		products.GET("", productHandler.GetAllProducts)
		products.GET("/:id", productHandler.GetProductByID)
		products.POST("", productHandler.CreateProduct)
		products.PATCH("/:id", productHandler.UpdateProduct)
		products.DELETE("/:id", productHandler.DeleteProduct)
	}

	orders := r.Group("/orders")
	{
		orders.GET("", middleware.AuthMiddleware(cfg.JWTSecret), orderHandler.GetAllOrders)
		orders.GET("/:id", middleware.AuthMiddleware(cfg.JWTSecret), orderHandler.GetOrderByID)
		orders.POST("", middleware.AuthMiddleware(cfg.JWTSecret), orderHandler.CreateOrder)
		orders.PATCH("/:id", middleware.AuthMiddleware(cfg.JWTSecret), orderHandler.UpdateOrder)
		orders.DELETE("/:id", middleware.AuthMiddleware(cfg.JWTSecret), orderHandler.DeleteOrder)
	}
}

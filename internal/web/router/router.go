package router

import (
	"net/http"

	"github.com/ernestngugi/sil-backend/internal/controller"
	"github.com/ernestngugi/sil-backend/internal/db"
	"github.com/ernestngugi/sil-backend/internal/repos"
	"github.com/ernestngugi/sil-backend/internal/web/auth"
	"github.com/ernestngugi/sil-backend/providers"
	"github.com/gin-gonic/gin"
)

type AppRouter struct {
	*gin.Engine
}

func BuildRouter(
	dB db.DB,
	oidcProvider providers.OpenID,
) *AppRouter {

	customerRepository := repos.NewCustomerRepository()
	orderRepository := repos.NewOrderRepository()

	customerController := controller.NewCustomerController(customerRepository)
	orderController := controller.NewOrderController(customerRepository, orderRepository)

	router := gin.Default()
	appRouter := router.Group("/v1")
	unauthenticatedUser := appRouter.Group("")
	appRouter.Use(authMiddleware(auth.NewAuthenticator(oidcProvider)))

	appRouter.POST("/orders", createOrder(dB, orderController))
	appRouter.GET("/orders/:id", orderByID(dB, orderController))
	appRouter.GET("/customers/:name", customerByName(dB, customerController))

	unauthenticatedUser.POST("/callback", handleLogin(dB, customerController, oidcProvider))

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error_message": "Endpoint not found"})
	})

	return &AppRouter{router}
}

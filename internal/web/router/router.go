package router

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/ernestngugi/sil-backend/internal/controller"
	"github.com/ernestngugi/sil-backend/internal/db"
	"github.com/ernestngugi/sil-backend/internal/forms"
	"github.com/ernestngugi/sil-backend/internal/model"
	"github.com/ernestngugi/sil-backend/internal/repos"
	"github.com/ernestngugi/sil-backend/internal/web/auth"
	"github.com/ernestngugi/sil-backend/providers"
	"github.com/gin-gonic/contrib/sessions"
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

	customerController := controller.NewCustomerRepository(customerRepository)
	orderController := controller.NewOrderController(customerRepository, orderRepository)

	router := gin.Default()
	appRouter := router.Group("/v1")
	appRouter.Use(authMiddleware(auth.NewAuthenticator(oidcProvider)))

	appRouter.POST("/orders", createOrder(dB, orderController))
	appRouter.GET("/orders/:id", orderByID(dB, orderController))
	appRouter.GET("/customers/:name", customerByName(dB, customerController))

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error_message": "Endpoint not found"})
	})

	return &AppRouter{router}
}

func customerByName(dB db.DB, customerController controller.CustomerController) func(c *gin.Context) {
	return func(c *gin.Context) {

		name := strings.TrimSpace(c.Param("name"))
		if name == "" {
			c.JSON(http.StatusBadRequest, gin.H{"success": false})
			return
		}

		customer, err := customerController.CustomerByName(c.Request.Context(), dB, name)
		if err != nil {

			c.JSON(http.StatusBadRequest, gin.H{"success": false})
			return
		}

		c.JSON(http.StatusOK, customer)
	}
}

func createOrder(dB db.DB, orderController controller.OrderController) func(c *gin.Context) {
	return func(c *gin.Context) {

		var form forms.CreateOrderForm

		err := c.BindJSON(&form)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false})
			return
		}

		order, err := orderController.CreateOrder(c.Request.Context(), dB, &form)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false})
			return
		}

		c.JSON(http.StatusOK, order)
	}
}

func orderByID(dB db.DB, orderController controller.OrderController) func(c *gin.Context) {
	return func(c *gin.Context) {

		orderID, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false})
			return
		}

		order, err := orderController.OrderByID(c.Request.Context(), dB, orderID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false})
			return
		}

		c.JSON(http.StatusOK, order)
	}
}

func authMiddleware(authAuthenticator auth.Authenticator) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		userInfo, err := authAuthenticator.TokenFromRequest(ctx, c.Request)
		if err != nil {
			fmt.Println("authentication error")
		}

		ctx = context.WithValue(ctx, model.CustomerKeyName, userInfo.Email)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func loginSession(
	oidcProvider providers.OpenID,
) func(c *gin.Context) {
	return func(c *gin.Context) {
		state, err := generateRandomState()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false})
			return
		}

		session := sessions.Default(c)
		session.Set("state", state)

		err = session.Save()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false})
			return
		}

		c.Redirect(http.StatusTemporaryRedirect, oidcProvider.AuthCodeURL(state))
	}
}

func generateRandomState() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	state := base64.StdEncoding.EncodeToString(b)

	return state, nil
}

func handleLogin(
	oidcProvider providers.OpenID,
) func(c *gin.Context) {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		if c.Query("state") != session.Get("state") {
			c.JSON(http.StatusBadRequest, gin.H{"success": false})
			return
		}

		ctx := c.Request.Context()

		token, err := oidcProvider.Exchange(ctx, c.Query("code"))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false})
			return
		}

		//add verify auth middleware
		idToken, err := oidcProvider.VerifyIDToken(ctx, token)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false})
			return
		}

		var profile map[string]interface{}
		if err := idToken.Claims(&profile); err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		session.Set("access_token", token.AccessToken)
		session.Set("profile", profile)
		if err := session.Save(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false})
			return
		}

		//redirect
	}
}

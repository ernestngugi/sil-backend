package router

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/ernestngugi/sil-backend/internal/controller"
	"github.com/ernestngugi/sil-backend/internal/db"
	"github.com/ernestngugi/sil-backend/internal/forms"
	"github.com/ernestngugi/sil-backend/internal/model"
	"github.com/ernestngugi/sil-backend/internal/web/auth"
	"github.com/ernestngugi/sil-backend/providers"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

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
		c.Redirect(http.StatusTemporaryRedirect, oidcProvider.AuthCodeURL(os.Getenv("OAUTH_STATE_STRING")))
	}
}

func handleLogin(
	dB db.DB,
	customerController controller.CustomerController,
	oidcProvider providers.OpenID,
) func(c *gin.Context) {
	return func(c *gin.Context) {

		if c.Query("state") != os.Getenv("OAUTH_STATE_STRING") {
			c.JSON(http.StatusBadRequest, gin.H{"success": false})
			return
		}

		ctx := c.Request.Context()

		token, err := oidcProvider.Exchange(ctx, c.Query("code"))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false})
			return
		}

		userInfo, err := oidcProvider.UserInfo(ctx, oauth2.StaticTokenSource(token))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false})
			return
		}

		customer, err := customerController.CreateCustomer(ctx, dB, &forms.CustomerCreateForm{Name: userInfo.Email})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false})
			return
		}

		c.Writer.Header().Set("X-SIL-TOKEN", token.AccessToken)
		c.JSON(http.StatusOK, customer)
	}
}

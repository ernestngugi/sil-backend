package router

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/coreos/go-oidc"
	"github.com/ernestngugi/sil-backend/internal/controller"
	"github.com/ernestngugi/sil-backend/internal/db"
	"github.com/ernestngugi/sil-backend/internal/forms"
	"github.com/ernestngugi/sil-backend/internal/model"
	"github.com/ernestngugi/sil-backend/internal/repos"
	"github.com/ernestngugi/sil-backend/internal/web/auth"
	"github.com/ernestngugi/sil-backend/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
	"syreclabs.com/go/faker"
)

func TestApplicationEndpoints(t *testing.T) {

	t.Setenv("DATABASE_URL", "postgres://savannah:password@localhost:5433/savannah?sslmode=disable")
	t.Setenv("OAUTH_STATE_STRING", "123456")

	ctx := context.Background()

	dB := db.InitDB()
	defer dB.Close()

	atProvider := mocks.NewMockATProvider()
	oidcProvider := mocks.NewMockOpenID()

	customerRepository := repos.NewCustomerRepository()
	orderRepository := repos.NewOrderRepository()

	customerController := controller.NewCustomerController(customerRepository)
	orderController := controller.NewOrderController(atProvider, customerRepository, orderRepository)

	testRouter := gin.Default()
	appRouter := testRouter.Group("/v1")
	unAuthenticatedUser := testRouter.Group("")
	appRouter.Use(authMiddleware(auth.NewAuthenticator(oidcProvider)))

	appRouter.POST("/orders", createOrder(dB, orderController))
	appRouter.GET("/orders/:id", orderByID(dB, orderController))
	appRouter.GET("/customers/:name", customerByName(dB, customerController))

	unAuthenticatedUser.POST("/callback", handleLogin(dB, customerController, oidcProvider))

	t.Run("can process oauth2 callback", func(t *testing.T) {

		token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDAwNzk1ODAsImlzX2FkbWluIjpmYWxzZSwibGFuZyI6ImVuLVVTIiwicmVmcmVzaCI6MTcwODQ2MDc4MCwic2Vzc2dGl2ZSIsInVzZXJfaWQiOjE0MDg4fQ.UxiW1NZ4daIGY9roHQat7Be9IfujzeAq-c5nH8Fkgz0"

		oidcProvider.User = &oidc.UserInfo{Email: faker.Internet().Email()}
		oidcProvider.AuthToken = &oauth2.Token{AccessToken: token}

		w := httptest.NewRecorder()

		queryParam := fmt.Sprintf("/callback?state=%v&code=%v", url.QueryEscape(os.Getenv("OAUTH_STATE_STRING")), token)

		req, err := http.NewRequest(http.MethodPost, queryParam, nil)
		assert.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")

		testRouter.ServeHTTP(w, req)

		data, err := io.ReadAll(w.Body)
		assert.NoError(t, err)

		var customer model.Customer

		err = json.Unmarshal(data, &customer)
		assert.NoError(t, err)

		assert.NotZero(t, customer.ID)
		assert.NotZero(t, customer.DateCreated)
		assert.NotZero(t, customer.DateModified)

		assert.Equal(t, w.Code, http.StatusOK)
		tokenHeader := w.Header().Get("X-SIL-TOKEN")
		assert.NotEmpty(t, tokenHeader)

		clearCustomerTable(ctx, dB)

	})
	t.Run("can create an order", func(t *testing.T) {

		customer := model.BuildCustomer()

		err := customerRepository.Save(ctx, dB, customer)
		assert.NoError(t, err)

		oidcProvider.User = &oidc.UserInfo{Email: customer.Name}

		oidcProvider.UserInfo(ctx, nil)

		w := httptest.NewRecorder()

		form := &forms.CreateOrderForm{
			Amount: 100,
			Item:   "item",
		}

		b, err := json.Marshal(form)
		assert.NoError(t, err)

		req, err := http.NewRequest(http.MethodPost, "/v1/orders", bytes.NewReader(b))
		assert.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")

		req.Header.Set("X-SIL-TOKEN", customer.Name)

		testRouter.ServeHTTP(w, req)

		data, err := io.ReadAll(w.Body)
		assert.NoError(t, err)

		var order model.Order

		err = json.Unmarshal(data, &order)
		assert.NoError(t, err)

		assert.NotZero(t, order.ID)
		assert.NotZero(t, order.DateCreated)
		assert.Equal(t, order.Amount, form.Amount)
		assert.Equal(t, order.Item, form.Item)
		assert.Equal(t, order.CustomerID, customer.ID)

		assert.Equal(t, w.Code, http.StatusOK)

		clearCustomerTable(ctx, dB)
		clearOrderTable(ctx, dB)
	})
}

func clearCustomerTable(ctx context.Context, dB db.DB) {
	dB.ExecContext(ctx, "DELETE FROM customers")
	dB.ExecContext(ctx, "ALTER SEQUENCE customers_id_seq RESTART WITH 1")
}

func clearOrderTable(ctx context.Context, dB db.DB) {
	dB.ExecContext(ctx, "DELETE FROM orders")
	dB.ExecContext(ctx, "ALTER SEQUENCE orders_id_seq RESTART WITH 1")
}

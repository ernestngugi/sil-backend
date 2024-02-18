package router

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/ernestngugi/sil-backend/internal/model"
	"github.com/ernestngugi/sil-backend/internal/web/auth"
	"github.com/ernestngugi/sil-backend/providers"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

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

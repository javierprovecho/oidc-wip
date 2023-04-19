package server

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/javierprovecho/oidc-wip/src/pkg/parse"
	"github.com/javierprovecho/oidc-wip/src/pkg/verify"
)

func Server(issuer, audience string) error {
	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	r.GET("/auth", func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"status": http.StatusUnauthorized, "message": "missing token"})
			return
		}

		token := strings.Split(authHeader, " ")

		if len(token) != 2 {
			c.JSON(http.StatusUnauthorized, gin.H{"status": http.StatusUnauthorized, "message": "missing token"})
			return
		}

		if issuer == "" {
			issuer, _ = parse.GetIssuer(token[1])
		}

		if claims, ok := verify.VerifyTokenWithIssuer(issuer, audience, token[1]); ok {
			c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "authorized", "claims": claims})
			return
		}

		c.JSON(http.StatusUnauthorized, gin.H{"status": http.StatusUnauthorized, "message": "missing token"})

	})

	return r.Run(":8080")
}

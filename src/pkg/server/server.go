package server

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/javierprovecho/oidc-wip/src/pkg/parse"
	"github.com/javierprovecho/oidc-wip/src/pkg/verify"
)

func Server(issuer, audience, namespace, serviceAccount, pod string) error {
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
			issuer, _ = parse.GetIssuerFromToken(token[1])
		}

		claims, ok := verify.VerifyTokenWithIssuer(token[1], issuer, audience)

		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"status": http.StatusUnauthorized, "message": "missing token"})
			return
		}

		if !verify.VerifyTokenWithSub(token[1], namespace, serviceAccount, pod) {
			c.JSON(http.StatusUnauthorized, gin.H{"status": http.StatusUnauthorized, "message": "missing token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "authorized", "claims": claims})

	})

	return r.Run(":8080")
}

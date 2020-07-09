// middleware.auth.go

package main

import (
	"net/http"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/thoas/go-funk"
)

func AuthenticationRequired(auths ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get("user")
		if user == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User needs to be signed in to access this service"})
			c.Abort()
			return
		}
		if len(auths) != 0 {
			authType := session.Get("authType")
			if authType == nil || !funk.ContainsString(auths, authType.(string)) {
				c.JSON(http.StatusForbidden, gin.H{"error": "Invalid request, restricted endpoint"})
				c.Abort()
				return
			}
		}
		// check if user exists
		c.Next()
	}
}

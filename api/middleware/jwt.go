// middleware to verify jwt token

package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nickhansel/nucleus/api/utils/token"
)

func JwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := token.TokenValid(c)
		id, _ := token.ExtractTokenID(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Unauthorized",
			})
			c.Abort()
			return
		}
		// send id to next middleware
		c.Set("id", id)
		c.Next()
	}
}

package api

import (
	"net/http"

	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func Test(c *gin.Context) {

	err := godotenv.Load("../.env")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	url := os.Getenv("DB_URL")

	c.JSON(http.StatusOK, gin.H{
		"customers": url,
	})
}

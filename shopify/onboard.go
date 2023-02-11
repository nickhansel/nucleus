package shopify

import (
	"github.com/gin-gonic/gin"
	"github.com/nickhansel/nucleus/model"
)

func GetInitialData(c *gin.Context) {
	org := c.MustGet("orgs").(model.Organization)

}

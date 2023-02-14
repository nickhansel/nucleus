package customers

import (
	"net/http"

	"fmt"
	"strconv"

	"github.com/nickhansel/nucleus/config"
	"github.com/nickhansel/nucleus/model"

	"github.com/gin-gonic/gin"
)

func GetCustomers(c *gin.Context) {
	orgId := c.Param("orgId")

	orgIdInt, err := strconv.Atoi(orgId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid URL parameter",
		})
		return
	}

	// query the db and find the customers where organizationId = orgId
	customers := []model.Customer{}
	res := config.DB.Where(&model.Customer{OrganizationID: int64(orgIdInt)}).Find(&customers).Error

	if res != nil {
		fmt.Println("res is not nil")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": res.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"customers": customers,
	})
}

package groups

import (
	"github.com/gin-gonic/gin"
	"github.com/nickhansel/nucleus/config"
	"github.com/nickhansel/nucleus/model"
	"net/http"
)

type Response struct {
	Group         model.CustomerGroup `json:"group"`
	CustomerCount int                 `json:"customer_count"`
}

func ListCustomerGroups(c *gin.Context) {
	org := c.MustGet("orgs").(model.Organization)

	var response []Response

	var groups []model.CustomerGroup
	err := config.DB.Where(&model.CustomerGroup{OrganizationID: org.ID}).Find(&groups).Error

	for i := 0; i < len(groups); i++ {
		var customersToCustomerGroups []model.CustomersToCustomerGroups
		// preload the customers that are in the customer group
		config.DB.Preload("Customer").Find(&customersToCustomerGroups)

		var customerCount int
		for _, customerToCustomerGroup := range customersToCustomerGroups {
			if customerToCustomerGroup.B == groups[i].ID {
				customerCount++
			}
		}

		response = append(response, Response{
			Group:         groups[i],
			CustomerCount: customerCount,
		})

	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": response,
	})
}

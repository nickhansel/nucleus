package customers

import (
	"net/http"

	"time"

	"github.com/gin-gonic/gin"
	"github.com/nickhansel/nucleus/config"
	"github.com/nickhansel/nucleus/model"
)

func CreateDefaultGroup(c *gin.Context) {
	var Customers []model.Customer
	err := config.DB.Find(&Customers).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	// create a customer group and connect all of the customers that have organizationId = 19
	var customerGroup model.CustomerGroup
	customerGroup.Name = "Test Group"
	customerGroup.OrganizationID = 19
	customerGroup.CreatedAt = time.Now()
	customerGroup.UpdatedAt = time.Now()
	config.DB.Create(&customerGroup)

	// add all of the customers to the customer group
	for _, customer := range Customers {
		if customer.OrganizationID == 19 {
			// add the customer to the Customers field of the customer group and connect them
			var CustomersToCustomerGroups model.CustomersToCustomerGroups

			CustomersToCustomerGroups.A = customer.ID
			CustomersToCustomerGroups.B = customerGroup.ID

			config.DB.Create(&CustomersToCustomerGroups)
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

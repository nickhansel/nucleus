package customers

import (
	"net/http"

	"github.com/nickhansel/nucleus/config"
	"github.com/nickhansel/nucleus/model"

	"github.com/gin-gonic/gin"
)

type Customers struct {
	Customers []model.Customer `json:"customers"`
}

func GetCustomerGroup(c *gin.Context) {
	// groupId := c.Param("groupId")

	// CustomersToCustomerGroups is a many-to-many relationship, so we need to get the customer group and then get the customers
	var customerGroup model.CustomerGroup
	// check if the customer group exists

	err := config.DB.First(&customerGroup, c.Param("groupId"))

	if err.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	var CustomersToCustomerGroups []model.CustomersToCustomerGroups
	// preload the customers that are in the customer group
	config.DB.Preload("Customer").Find(&CustomersToCustomerGroups)

	for _, customerToCustomerGroup := range CustomersToCustomerGroups {
		if customerToCustomerGroup.B == customerGroup.ID {
			customerGroup.Customers = append(customerGroup.Customers, customerToCustomerGroup.Customer)
		}
	}

	// check length of customerGroup.Customers
	if len(customerGroup.Customers) == 0 {
		customerGroup.Customers = []model.Customer{}
	}

	// get the customers that are in the customer group
	// config.DB.Model(&customerGroup).Association("Customers").Find(&customerGroup.Customers)

	c.JSON(http.StatusOK, customerGroup)

}

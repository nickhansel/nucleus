package groups

import (
	"fmt"
	"net/http"

	"time"

	"github.com/gin-gonic/gin"
	"github.com/nickhansel/nucleus/config"
	fb "github.com/nickhansel/nucleus/fb/audiences"
	"github.com/nickhansel/nucleus/model"
)

type Body struct {
	IDs  []int64 `json:"ids"`
	Name string  `json:"name"`
}

func CreateCustomerGroup(c *gin.Context) {
	org := c.MustGet("orgs").(model.Organization)
	if org.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	var body Body
	bindErr := c.BindJSON(&body)
	if bindErr != nil {
		return
	}

	if len(body.IDs) == 0 {
		// abort if the custom audience doesn't exist
		c.JSON(http.StatusBadRequest, gin.H{"error": "No customers selected!"})
		return
	}

	var Customers []model.Customer
	// find customers with the ids in the body
	err := config.DB.Where("id IN (?)", body.IDs).Find(&Customers).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	if body.Name == "" || body.Name == "Square customers" || body.Name == "Shopify customers" || body.Name == "Default Group" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name can not be empty, \"Square customers\", \"Shopify customers\" or \"Default Group\"!"})
		return
	}

	var customerGroup model.CustomerGroup
	customerGroup.Name = body.Name
	customerGroup.OrganizationID = org.ID
	customerGroup.CreatedAt = time.Now()
	customerGroup.UpdatedAt = time.Now()
	config.DB.Create(&customerGroup)

	// add all the customers to the customer group
	for _, customer := range Customers {
		if customer.OrganizationID == org.ID {
			// add the customer to the Customers field of the customer group and connect them
			var CustomersToCustomerGroups model.CustomersToCustomerGroups

			CustomersToCustomerGroups.A = customer.ID
			CustomersToCustomerGroups.B = customerGroup.ID

			config.DB.Create(&CustomersToCustomerGroups)
		}
	}
	id, err := fb.CreateAudience(c, customerGroup.ID)
	if err != nil {
		config.DB.Delete(&customerGroup)
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Error creating audience: %s", err.Error())})
		return
	}

	if id == "" {
		config.DB.Delete(&customerGroup)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error creating audience!"})
		return
	}

	// update the customer group with the facebook audience id
	customerGroup.FbCustomAudienceID = id

	if id == "" || err != nil {
		config.DB.Delete(&customerGroup)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error creating audience!"})
		return
	}
	config.DB.Save(&customerGroup)

	c.JSON(http.StatusOK, gin.H{"result": customerGroup})
}

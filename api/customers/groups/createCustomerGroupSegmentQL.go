package groups

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nickhansel/nucleus/config"
	fb "github.com/nickhansel/nucleus/fb/audiences"
	"github.com/nickhansel/nucleus/model"
	"github.com/nickhansel/nucleus/segmentQL"
	"net/http"
	"time"
)

type SegmentQLRequestBody struct {
	ItemID           int64   `json:"item_id"`
	StartDate        string  `json:"start_date"`
	EndDate          string  `json:"end_date"`
	MinPurchasePrice float64 `json:"min_purchase_price"`
	MaxPurchasePrice float64 `json:"max_purchase_price"`
	Name             string  `json:"name"`
	IsVariation      bool    `json:"is_variation"`
}

func CreateCustomerGroupSegmentQL(c *gin.Context) {
	org := c.MustGet("orgs").(model.Organization)
	if org.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	query := SegmentQLRequestBody{}
	err := c.BindJSON(&query)
	if err != nil {
		return
	}

	if query.Name == "" || query.Name == "Square customers" || query.Name == "Shopify customers" || query.Name == "Default Group" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name can not be empty, \"Square customers\", \"Shopify customers\" or \"Default Group\"!"})
		return
	}

	ql := segmentQL.Parse(int64(query.ItemID), org.ID, query.StartDate, query.EndDate, query.MinPurchasePrice, query.MaxPurchasePrice)

	var body Body
	//convert the query to a slice of ints
	for _, id := range ql {
		body.IDs = append(body.IDs, int64(id))
	}

	body.Name = query.Name

	if len(ql) == 0 {
		// abort if the custom audience doesn't exist
		c.JSON(http.StatusBadRequest, gin.H{"error": "No customers selected!"})
		return
	}

	var Customers []model.Customer
	// find customers with the ids in the body
	err = config.DB.Where("id IN (?)", body.IDs).Find(&Customers).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// update the customer group with the facebook audience id
	customerGroup.FbCustomAudienceID = id
	config.DB.Save(&customerGroup)

	id, err = fb.CreateAudience(c, customerGroup.ID)
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

	c.JSON(http.StatusOK, gin.H{"result": customerGroup})
}

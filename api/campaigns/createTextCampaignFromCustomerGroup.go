package campaign

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nickhansel/nucleus/config"
	"github.com/nickhansel/nucleus/model"
)

type Body struct {
	Type            string  `json:"type"`
	Budget          float64 `json:"budget"`
	Name            string  `json:"name"`
	SendTime        string  `json:"send_time"`
	TextBody        string  `json:"text_body"`
	CustomerGroupID int     `json:"customer_group_id"`
}

func CreateTextCampaign(c *gin.Context) {
	org := c.MustGet("orgs").(model.Organization)

	var body Body
	if err := c.ShouldBindJSON(&body); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// make sure there are no values in the body that are empty
	if body.Type == "" || body.Budget == 0 || body.Name == "" || body.SendTime == "" || body.TextBody == "" || body.CustomerGroupID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields!"})
		return
	}

	var campaign model.Campaign
	campaign.OrganizationID = org.ID
	campaign.Budget = body.Budget
	campaign.Type = body.Type
	campaign.Name = body.Name
	campaign.CreatedAt = time.Now().String()
	campaign.IsTextCampaign = true

	config.DB.Create(&campaign)

	// find all of the customers in the customer group
	var customerGroup model.CustomerGroup
	// check if the customer group exists

	err := config.DB.First(&customerGroup, body.CustomerGroupID)

	if err.Error != nil {
		c.AbortWithError(http.StatusBadRequest, err.Error)
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

	var TargetCustomers []string
	for _, customer := range customerGroup.Customers {
		if customer.PhoneNumber != "" {
			TargetCustomers = append(TargetCustomers, customer.PhoneNumber)
		}
	}

	// create the text campaign

	var TextCampaign model.TextCampaign
	TextCampaign.CampaignID = campaign.ID
	TextCampaign.Name = body.Name
	TextCampaign.TargetNumbers = TargetCustomers
	TextCampaign.SendTime = body.SendTime
	TextCampaign.Body = body.TextBody
	TextCampaign.SendTime = body.SendTime
	TextCampaign.From = org.TwilioNumber

	config.DB.Create(&TextCampaign)

	c.JSON(http.StatusOK, gin.H{
		"data":          campaign,
		"text_campaign": TextCampaign,
	})
}

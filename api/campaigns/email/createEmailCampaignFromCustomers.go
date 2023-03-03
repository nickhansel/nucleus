package email

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nickhansel/nucleus/config"
	"github.com/nickhansel/nucleus/cron/email"
	"github.com/nickhansel/nucleus/model"
	"net/http"
	"strconv"
	"time"
)

type EmailCampaignBody struct {
	Type            string  `json:"type"`
	Budget          float64 `json:"budget"`
	Name            string  `json:"name"`
	SendTime        string  `json:"send_time"`
	Subject         string  `json:"subject"`
	HtmlContent     string  `json:"htmlContent"`
	CustomerGroupID string  `json:"customer_group_id"`
	TemplateId      int64   `json:"template_id"`
}

type EmailBody struct {
	Sender struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	} `json:"sender"`
	To []struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	} `json:"to"`
	Subject     string `json:"subject"`
	HtmlContent string `json:"htmlContent"`
}

func CreateEmailCampaign(c *gin.Context) {
	org := c.MustGet("orgs").(model.Organization)

	if org.IsSendinblueAuthed == false {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email provider not authorized"})
		return
	}

	// get the campaign body
	var campaignBody EmailCampaignBody
	if err := c.ShouldBindJSON(&campaignBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// make sure there are no values in the body that are empty
	if campaignBody.Type == "" || campaignBody.Budget == 0 || campaignBody.Name == "" || campaignBody.SendTime == "" || campaignBody.Subject == "" || campaignBody.HtmlContent == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields!"})
		return
	}

	localTimeZone := time.Now().Location()
	// convert the send time to central time
	centralTime, err := time.ParseInLocation("2006-01-02 15:04:05", campaignBody.SendTime, localTimeZone)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("send time: ", centralTime)
	fmt.Println("central time: ", time.Now().In(time.FixedZone("America/Chicago", -6*60*60)))
	// compare central time to the current time in central time
	if centralTime.Before(time.Now().In(time.FixedZone("America/Chicago", -6*60*60))) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Send time must be in the future!"})
		return
	}

	var CustomersToCustomerGroups []model.CustomersToCustomerGroups
	// preload the customers that are in the customer group
	config.DB.Preload("Customer").Find(&CustomersToCustomerGroups)
	// find all the customers in the customer group
	var customerGroup model.CustomerGroup

	//convert customer group id to int64 from string
	convertedId, err := strconv.ParseInt(campaignBody.CustomerGroupID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = config.DB.First(&customerGroup, convertedId).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for _, customerToCustomerGroup := range CustomersToCustomerGroups {
		if customerToCustomerGroup.B == customerGroup.ID {
			customerGroup.Customers = append(customerGroup.Customers, customerToCustomerGroup.Customer)
		}
	}

	var TargetCustomers []string
	for _, customer := range customerGroup.Customers {
		if customer.EmailAddress != "" && customer.EmailUnsubscribed == false {
			TargetCustomers = append(TargetCustomers, customer.EmailAddress)
			//var currentCustomer model.Customer
			//config.DB.First(&currentCustomer, customer.ID)
			//currentCustomer.DatesReceivedEmail = append(currentCustomer.DatesReceivedEmail, time.Now().String())
			//config.DB.Save(&currentCustomer)
		}
	}

	var campaign model.Campaign
	campaign.OrganizationID = org.ID
	campaign.Budget = campaignBody.Budget
	campaign.Type = campaignBody.Type
	campaign.Name = campaignBody.Name
	campaign.CreatedAt = time.Now().String()
	campaign.IsEmailCampaign = true
	campaign.CustomersTargeted = int32(len(TargetCustomers))
	campaign.CustomerGroupID = convertedId

	// save the campaign
	config.DB.Create(&campaign)

	// create the email body
	var emailBody EmailBody
	emailBody.Sender.Name = org.Name
	emailBody.Sender.Email = org.SendinblueEmail
	emailBody.Subject = campaignBody.Subject
	emailBody.HtmlContent = campaignBody.HtmlContent

	var EmailCampaign model.EmailCampaign
	EmailCampaign.CampaignID = campaign.ID
	EmailCampaign.SendTime = campaignBody.SendTime
	EmailCampaign.From = org.SendinblueEmail
	EmailCampaign.Subject = campaignBody.Subject
	EmailCampaign.HTML = campaignBody.HtmlContent
	EmailCampaign.Text = campaignBody.Subject
	EmailCampaign.TargetEmails = TargetCustomers

	// save the email campaign
	err = config.DB.Create(&EmailCampaign).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var emailCampaignAnalytics model.EmailCampaignAnalytics
	emailCampaignAnalytics.EmailCampaignID = EmailCampaign.ID
	//date in format YYYY-MM-DD
	now := time.Now()
	emailCampaignAnalytics.Date = now.Format("2006-01-02")

	// save the email campaign analytics
	err = config.DB.Create(&emailCampaignAnalytics).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for _, customer := range TargetCustomers {
		var to struct {
			Email string `json:"email"`
			Name  string `json:"name"`
		}
		to.Email = customer
		emailBody.To = append(emailBody.To, to)
	}

	// send the email
	email.ScheduleEmailTasks(campaignBody.SendTime, EmailCampaign, org)

	c.JSON(http.StatusOK, gin.H{
		"message":        "Email campaign created!",
		"email_campaign": EmailCampaign,
	})
}

package fb

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"time"

	"github.com/gin-gonic/gin"
	"github.com/nickhansel/nucleus/config"
	"github.com/nickhansel/nucleus/model"
)

type CampaignBody struct {
	Name                string   `json:"name"`
	Objective           string   `json:"objective"`
	Status              string   `json:"status"`
	SpecialAdCategories []string `json:"special_ad_categories"`
	AccessToken         string   `json:"access_token"`
	CampaignType        string   `json:"campaign_type"`
}

func CreateCampaign(c *gin.Context) {
	// add the campaignbody to the form data of the request
	var campaignBody CampaignBody
	c.BindJSON(&campaignBody)

	// get the org from the context
	org := c.MustGet("orgs").(model.Organization)

	// get the org from the db
	var organization model.Organization
	config.DB.Where("id = ?", org.ID).First(&organization)

	// get the access token from the org
	accessToken := organization.FbAccessToken

	// get the page id from the org
	adId := organization.FbAdAccountID

	// create the url
	url := fmt.Sprintf("https://graph.facebook.com/v15.0/act_%s/campaigns", adId)

	// create the request
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		log.Fatal("wtf")
	}

	// add the form data to the request
	q := req.URL.Query()
	q.Add("name", campaignBody.Name)
	q.Add("objective", campaignBody.Objective)
	q.Add("status", campaignBody.Status)
	q.Add("special_ad_categories", "[]")
	q.Add("access_token", accessToken)
	req.URL.RawQuery = q.Encode()

	// make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("something went wrong")
	}
	defer resp.Body.Close()

	// convert the response to json
	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	var Campaign model.Campaign
	Campaign.CampaignId = result["id"].(string)
	// created at should be the time right now but as a string
	Campaign.CreatedAt = time.Now().String()
	Campaign.Type = campaignBody.CampaignType
	Campaign.OrganizationID = org.ID
	Campaign.Budget = 0

	// save the campaign to the db
	config.DB.Save(&Campaign)

	// get the id from the DB of the campaign
	var campaign model.Campaign
	config.DB.Where("campaign_id = ?", Campaign.CampaignId).First(&campaign)

	// add the campaign to the org
	config.DB.Model(&organization).Association("Campaigns").Append(&campaign)

	var FbCampaign model.FbCampaign
	FbCampaign.CampaignID = campaign.ID
	FbCampaign.Name = campaignBody.Name
	FbCampaign.Objective = campaignBody.Objective
	FbCampaign.Status = campaignBody.Status

	// save the fb campaign to the db
	config.DB.Save(&FbCampaign)

	// return the response
	c.JSON(200, gin.H{
		"message":    "Campaign created",
		"fbcampaign": FbCampaign,
		"campaign":   Campaign,
	})

}

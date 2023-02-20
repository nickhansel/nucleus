package fb

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nickhansel/nucleus/config"
	"github.com/nickhansel/nucleus/model"
	"io"
	"log"
	"net/http"
	"time"
)

type AdsetBody struct {
	Name             string `json:"name"`
	OptimizationGoal string `json:"optimization_goal"`
	BillingEvent     string `json:"billing_event"`
	BidAmount        int    `json:"bid_amount"`
	DailyBudget      int    `json:"daily_budget"`
	CampaignId       string `json:"campaign_id"`
	FbCampaignId     int64  `json:"fb_campaign_id"`
	Status           string `json:"status"`
	StartTime        string `json:"start_time"`
	EndTime          string `json:"end_time"`
	Targeting        struct {
		AgeMin int `json:"age_min"`
		AgeMax int `json:"age_max"`

		CustomAudiences []struct {
			ID string `json:"id"`
		} `json:"custom_audiences"`

		GeoLocations struct {
			Countries []string `json:"countries"`
			Regions   []string `json:"regions"`
			ZipCodes  []string `json:"zip_codes"`
		} `json:"geo_locations"`
	} `json:"targeting"`
}

func CreateAdSet(c *gin.Context) {
	var adsetBody AdsetBody
	err := c.BindJSON(&adsetBody)
	if err != nil {
		fmt.Println(err)
		return
	}

	org := c.MustGet("orgs").(model.Organization)
	fbAccId := org.FbAdAccountID
	fbAccessToken := org.FbAccessToken

	if org.FbAccessToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No Facebook Access Token"})
		return
	}

	if fbAccId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No Facebook Ad Account ID"})
		return
	}

	if fbAccessToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No Facebook Access Token"})
		return
	}

	url := fmt.Sprintf("https://graph.facebook.com/v15.0/act_%s/adsets", fbAccId)

	// create the request
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		log.Fatal("wtf")
	}

	if adsetBody.DailyBudget < 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Daily Budget must be greater than $1"})
		return
	}

	q := req.URL.Query()
	q.Add("name", adsetBody.Name)
	q.Add("start_time", adsetBody.StartTime)
	q.Add("end_time", adsetBody.EndTime)
	q.Add("optimization_goal", adsetBody.OptimizationGoal)
	q.Add("billing_event", adsetBody.BillingEvent)
	q.Add("bid_amount", fmt.Sprintf("%d", adsetBody.BidAmount))
	q.Add("daily_budget", fmt.Sprintf("%d", adsetBody.DailyBudget))
	q.Add("campaign_id", adsetBody.CampaignId)
	q.Add("status", adsetBody.Status)
	q.Add("targeting", fmt.Sprintf(`{"age_min":%d,"age_max":%d,"custom_audiences":[{"id":"%s"}],"geo_locations":{"countries":["US"]}}`, adsetBody.Targeting.AgeMin, adsetBody.Targeting.AgeMax, adsetBody.Targeting.CustomAudiences[0].ID))
	q.Add("access_token", fbAccessToken)
	req.URL.RawQuery = q.Encode()

	// send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	// decode the response
	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	fmt.Println(result)

	//get num days between start and end time where they are ibn the form 2006-01-02 23:59:59-07:00
	startDateDate, err := time.Parse("2006-01-02 15:04:05-07:00", adsetBody.StartTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	endDateDate, err := time.Parse("2006-01-02 15:04:05-07:00", adsetBody.EndTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	numDays := endDateDate.Sub(startDateDate).Hours() / 24

	fmt.Println(result)

	// create the adset in the database
	var adset model.FbAdset
	adset.Name = adsetBody.Name
	adset.LifetimeBudget = float64(adsetBody.DailyBudget*int(numDays)) / 100
	adset.StartTime = adsetBody.StartTime
	adset.EndTime = adsetBody.EndTime
	adset.CampaignID = adsetBody.CampaignId
	adset.BidAmount = float64(adsetBody.BidAmount) / 100
	adset.OptimizationGoal = adsetBody.OptimizationGoal
	adset.Status = adsetBody.Status
	adset.FbTargetID = 12
	adset.FbCampaignID = adsetBody.FbCampaignId
	adset.FbAdsetID = result["id"].(string)

	// create the target in the database
	if dbc := config.DB.Create(&adset); dbc.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": dbc.Error})
		return
	}

	var fbTarget model.FbTarget
	fbTarget.AgeMin = int32(adsetBody.Targeting.AgeMin)
	fbTarget.AgeMax = int32(adsetBody.Targeting.AgeMax)
	// add the customaudience ids to an array
	var customAudienceIds []string
	for _, customAudience := range adsetBody.Targeting.CustomAudiences {
		customAudienceIds = append(customAudienceIds, customAudience.ID)
	}

	fbTarget.CustomAudiences = customAudienceIds
	fbTarget.FbAdsetID = adset.ID

	var customersTargeted int32

	var customerGroup model.CustomerGroup
	for _, customAudience := range adsetBody.Targeting.CustomAudiences {
		config.DB.Where("\"fb_custom_audience_id\" = ?", customAudience.ID).First(&customerGroup)
		//	find customers in the customer group
		var customers []model.CustomersToCustomerGroups
		config.DB.Where("\"B\" = ?", customerGroup.ID).Find(&customers)
		customersTargeted += int32(len(customers))
	}

	//update the fb_campaigns table with the number of customers targeted
	var campaign model.Campaign
	config.DB.Where("\"campaign_id\" = ?", adsetBody.CampaignId).First(&campaign)
	campaign.CustomersTargeted = customersTargeted

	if dbc := config.DB.Save(&campaign); dbc.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": dbc.Error})
		return
	}

	if dbc := config.DB.Create(&fbTarget); dbc.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": dbc.Error})
		return
	}

	//update the fb_adset with the fb_target_id
	adset.FbTargetID = fbTarget.ID
	if dbc := config.DB.Save(&adset); dbc.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": dbc.Error})
		return
	}

	var geoLocation model.GeoLocation
	geoLocation.Countries = adsetBody.Targeting.GeoLocations.Countries
	geoLocation.Regions = adsetBody.Targeting.GeoLocations.Regions
	geoLocation.ZipCodes = adsetBody.Targeting.GeoLocations.ZipCodes
	geoLocation.FbTargetID = fbTarget.ID

	if dbc := config.DB.Create(&geoLocation); dbc.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": dbc.Error})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": adset})

}

// create a test json body for adsetbody
// {
// 	"name": "test adset",
// 	"optimization_goal": "REACH",
// 	"billing_event": "IMPRESSIONS",
// 	"bid_amount": 100,
// 	"daily_budget": 1000,
// 	"campaign_id": "23843951054510123",
// 	"fb_campaign_id": 1,
// 	"status": "ACTIVE",
// 	"start_time": "2021-09-01T00:00:00-0700",
// 	"end_time": "2021-09-30T00:00:00-0700",
// 	"targeting": {
// 		"age_min": 18,
// 		"age_max": 65,
// 		"custom_audiences": [{
// 			"id": "23843951054510123"
// 		}],
// 		"geo_locations": {
// 			"countries": ["US"],
// 			"regions": [],
// 			"zip_codes": []
// 		}
// 	}
// }

package shopify

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/nickhansel/nucleus/config"
	"github.com/nickhansel/nucleus/model"
	"net/http"
	"strconv"
)

type Body struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

func Oauth(c *gin.Context) {
	org := c.MustGet("orgs").(model.Organization)

	var body Body
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	reqUrl := "https://" + body.Name + ".myshopify.com/admin/oauth/access_token?client_id=7636322da1ba1eadb809b738ca4c0129&client_secret=16c5e8ec013752ec84be0924df78532e&code=" + body.Code

	var client = &http.Client{}
	req, err := http.NewRequest("POST", reqUrl, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	var result map[string]interface{}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return
	}

	//	check if result["access_token"] is nil
	//	if it is, return an error
	//	if it isn't, return the access token
	if result["access_token"] == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "error in getting access token",
		})
		return
	}

	shopReqUrl := "https://" + body.Name + ".myshopify.com/admin/api/2023-01/shop.json"

	req, err = http.NewRequest("GET", shopReqUrl, nil)
	//	set the header
	req.Header.Set("X-Shopify-Access-Token", result["access_token"].(string))

	resp, err = client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	var shop map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&shop)
	if err != nil {
		return
	}

	//	check if shop["shop"] is nil
	//	if it is, return an error
	//	if it isn't, return the shop id
	if shop["shop"] == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "error in getting shop id",
		})
		return
	}

	// check if the org is connected to a shopify pos
	var orgPos []model.Pos
	config.DB.Where("\"organizationId\" = ?", org.ID).Find(&orgPos)

	for _, pos := range orgPos {
		if pos.Name == "Shopify" {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "organization already connected to a shopify pos",
			})
			return
		}
	}

	var pos model.Pos
	pos.Name = "Shopify"
	pos.MerchantID = strconv.FormatInt(int64(shop["shop"].(map[string]interface{})["id"].(float64)), 10)
	pos.AccessToken = result["access_token"].(string)
	pos.RefreshToken = ""
	pos.OrganizationID = org.ID

	err = config.DB.Create(&pos).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	var store model.StoreLocation
	store.Name = shop["shop"].(map[string]interface{})["name"].(string)
	store.Type = "Shopify"

	if shop["shop"].(map[string]interface{})["address1"] != nil {
		store.AddressLine1 = shop["shop"].(map[string]interface{})["address1"].(string)
	} else {
		store.AddressLine1 = ""
	}

	if shop["shop"].(map[string]interface{})["city"] != nil {
		store.Locality = shop["shop"].(map[string]interface{})["city"].(string)
	} else {
		store.Locality = ""
	}

	if shop["shop"].(map[string]interface{})["province"] != nil {
		store.AdministrativeDistrictLevel1 = shop["shop"].(map[string]interface{})["province"].(string)
	} else {
		store.AdministrativeDistrictLevel1 = ""
	}
	store.Country = "US"
	store.CreatedAt = shop["shop"].(map[string]interface{})["created_at"].(string)
	store.LanguageCode = shop["shop"].(map[string]interface{})["primary_locale"].(string)
	store.Currency = shop["shop"].(map[string]interface{})["currency"].(string)
	store.BusinessName = shop["shop"].(map[string]interface{})["domain"].(string)
	store.OrganizationID = org.ID
	store.MerchantID = pos.MerchantID
	store.PosID = strconv.FormatInt(pos.ID, 10)
	store.Status = "ACTIVE"
	var emptyArray []string
	store.Capabilities = emptyArray
	store.Timezone = shop["shop"].(map[string]interface{})["iana_timezone"].(string)
	if shop["shop"].(map[string]interface{})["zip"] != nil {
		store.PostalCode = shop["shop"].(map[string]interface{})["zip"].(string)
	} else {
		store.PostalCode = ""
	}

	err = config.DB.Create(&store).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"pos":     pos,
	})
}

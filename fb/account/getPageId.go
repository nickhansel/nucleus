package fb

import (
	"encoding/json"
	"fmt"

	// "log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/nickhansel/nucleus/config"
	"github.com/nickhansel/nucleus/model"
)

type Body struct {
	PageName string `json:"name"`
}

func GetPageID(c *gin.Context) {
	org := c.MustGet("orgs").(model.Organization)
	fbAccessToken := org.FbAccessToken

	if fbAccessToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No Facebook Access Token"})
		return
	}

	var body Body
	err := c.BindJSON(&body)
	if err != nil {
		return
	}

	// get the page id
	url := fmt.Sprintf("https://graph.facebook.com/v15.0/%s?fields=id,name,fan_count,picture,is_verified&access_token=%s", body.PageName, fbAccessToken)

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error creating request"})
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error sending request"})
		return
	}

	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	if result["error"] != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Page does not exist."})
		return
	}

	// var orgData model.Organization
	// config.DB.Where("id = ?", org.ID).First(&orgData)

	org.FbPageID = result["id"].(string)
	org.FbPageName = result["name"].(string)
	org.FbPageFanCount = int32(result["fan_count"].(float64))
	org.FbPageImgURL = result["picture"].(map[string]interface{})["data"].(map[string]interface{})["url"].(string)

	config.DB.Save(&org)

	c.JSON(http.StatusOK, gin.H{
		"org": org,
		"res": result,
	})
}

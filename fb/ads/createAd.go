// {
// 	"name": "Otis",
// 	"object_story_spec": {
// 	  "image_url": "https://www.thesprucepets.com/thmb/7TDhfkK5CAKBWEaJfez6607J48Y=/1500x0/filters:no_upscale():max_bytes(150000):strip_icc()/chinese-dog-breeds-4797219-hero-2a1e9c5ed2c54d00aef75b05c5db399c.jpg",
// 	  "link_data": {
// 		"link": "https://www.akc.org/",
// 		"message": "This is a test"
// 	  },
// 	  "page_id": "104298420934949"
// 	}
//   }

// create an ad creative and use the reponse from that to create an ad

// {
// 	"name": "Otis",
// 	"adset_id": "23841000000000000",
// 	"creative": {
// 			"creative_id": "23841000000000000"
// 	},
// 	"status": "PAUSED",
// 	"access_token": "token"
// }
// https://graph.facebook.com/v15.0/act_<AD_ACCOUNT_ID>/ads

// TODO: save the data from the ad creative to the database
// TODO: check if the ad creative already exists

package fb

import (
	"encoding/json"
	"fmt"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nickhansel/nucleus/aws"
	"github.com/nickhansel/nucleus/model"
	// "github.com/nickhansel/nucleus/config"
)

type AdBody struct {
	Name    string `json:"name"`
	Link    string `json:"link"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

func createAdCreative(c *gin.Context) (creativeId string) {
	org := c.MustGet("orgs").(model.Organization)

	url := aws.UploadImage(c, "test", "test")

	apiUrl := fmt.Sprintf("https://graph.facebook.com/v15.0/act_%s/adcreatives", org.FbAdAccountID)

	// get the query params
	name := c.Query("name")
	link := c.Query("link")
	message := c.Query("message")

	req, err := http.NewRequest("POST", apiUrl, nil)
	if err != nil {
		fmt.Println(err)
	}

	q := req.URL.Query()
	q.Add("name", name)
	q.Add("object_story_spec", fmt.Sprintf(`{"image_url": "%s", "link_data": {"link": "%s", "message": "%s"}, "page_id": "%s"}`, url, link, message, org.FbPageID))
	q.Add("access_token", org.FbAccessToken)
	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
	respsonse, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	defer respsonse.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(respsonse.Body).Decode(&result)

	c.JSON(200, gin.H{
		"result": result,
	})
	return "6136298909426"
}

func CreateAd(c *gin.Context) {
	creativeId := createAdCreative(c)
	org := c.MustGet("orgs").(model.Organization)

	if creativeId == "" {
		c.JSON(500, gin.H{
			"error": "Could not create ad creative",
		})
		return
	}

	apiUrl := fmt.Sprintf("https://graph.facebook.com/v15.0/act_%s/ads", org.FbAdAccountID)

	req, err := http.NewRequest("POST", apiUrl, nil)
	if err != nil {
		fmt.Println(err)
	}

	adName := c.Query("ad_name")
	adSetId := c.Query("adset_id")

	q := req.URL.Query()
	q.Add("name", adName)
	q.Add("adset_id", adSetId)
	q.Add("creative", fmt.Sprintf(`{"creative_id": "%s"}`, creativeId))
	q.Add("status", c.Query("status"))
	q.Add("access_token", org.FbAccessToken)

	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
	respsonse, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	defer respsonse.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(respsonse.Body).Decode(&result)

	c.JSON(200, gin.H{
		"result": result,
	})
}

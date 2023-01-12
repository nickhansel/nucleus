package sendgrid

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/nickhansel/nucleus/config"
	"github.com/nickhansel/nucleus/model"
	"github.com/sendgrid/sendgrid-go"
)

type Body struct {
	Nickname    string `json:"nickname" binding:"required"`
	FromEmail   string `json:"from_email" binding:"required"`
	FromName    string `json:"from_name" binding:"required"`
	ReplyTo     string `json:"reply_to" binding:"required"`
	ReplyToName string `json:"reply_to_name" binding:"required"`
	Address     string `json:"address" binding:"required"`
	Address2    string `json:"address2" binding:"required"`
	State       string `json:"state" binding:"required"`
	City        string `json:"city" binding:"required"`
	Country     string `json:"country" binding:"required"`
	Zip         string `json:"zip" binding:"required"`
}

type TwilioResponse struct {
	ID          int32  `json:"id"`
	Nickname    string `json:"nickname"`
	FromEmail   string `json:"from_email"`
	FromName    string `json:"from_name"`
	ReplyTo     string `json:"reply_to"`
	ReplyToName string `json:"reply_to_name"`
	Address     string `json:"address"`
	Address2    string `json:"address2"`
	State       string `json:"state"`
	City        string `json:"city"`
	Country     string `json:"country"`
	Zip         string `json:"zip"`
	Verified    bool   `json:"verified"`
	Locked      bool   `json:"locked"`
}

func SendVerificationEmail(nickname, fromEmail, fromName, replyTo, replyToName, address, address2, state, city, country, zip string) (int, TwilioResponse) {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file for sendgrid")
	}

	apiKey := os.Getenv("SENDGRID_API_KEY")

	host := "https://api.sendgrid.com"

	request := sendgrid.GetRequest(apiKey, "/v3/verified_senders", host)
	request.Method = "POST"

	request.Body = []byte(`{
		"nickname": "` + nickname + `",
		"from_email": "` + fromEmail + `",
		"from_name": "` + fromName + `",
		"reply_to": "` + replyTo + `",
		"reply_to_name": "` + replyToName + `",
		"address": "` + address + `",
		"address2": "` + address2 + `",
		"state": "` + state + `",
		"city": "` + city + `",
		"country": "` + country + `",
		"zip": "` + zip + `"
	}`)

	response, err := sendgrid.API(request)

	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}

	// convert body to struct
	var twilioResponse TwilioResponse
	json.Unmarshal([]byte(response.Body), &twilioResponse)

	return response.StatusCode, twilioResponse
}

func VerifySendgridEmail(c *gin.Context) {
	org, _ := c.MustGet("orgs").(model.Organization)

	var body Body

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err, res := SendVerificationEmail(
		body.Nickname,
		body.FromEmail,
		body.FromName,
		body.ReplyTo,
		body.ReplyToName,
		body.Address,
		body.Address2,
		body.State,
		body.City,
		body.Country,
		body.Zip,
	)

	if err != 201 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	var organization model.Organization
	orgErr := config.DB.Where(&model.Organization{ID: org.ID}).First(&organization).Error

	if orgErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": orgErr,
		})
		return
	}

	organization.SendgridEmail = body.FromEmail
	organization.SendgridID = res.ID

	config.DB.Save(&organization)

	c.JSON(200, gin.H{
		"message":  "Verification email sent to" + body.FromEmail,
		"response": res,
	})

}

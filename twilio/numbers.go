// Download the helper library from https://www.twilio.com/docs/go/install
package twilio

import (
	"fmt"

	"os"

	"log"

	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/nickhansel/nucleus/config"
	"github.com/nickhansel/nucleus/model"
	"github.com/twilio/twilio-go"
	api "github.com/twilio/twilio-go/rest/api/v2010"
)

func ListAvaliableNumbers(areaCode int) string {
	err := godotenv.Load("../.env")

	if err != nil {
		log.Fatal("Error loading .env file for twilio")
	}

	os.Setenv("TWILIO_ACCOUNT_SID", os.Getenv("TWILIO_ACCOUNT_SID"))
	os.Setenv("TWILIO_AUTH_TOKEN", os.Getenv("TWILIO_AUTH_TOKEN"))

	client := twilio.NewRestClient()

	params := &api.ListAvailablePhoneNumberLocalParams{}
	params.SetAreaCode(areaCode)
	params.SetLimit(1)

	phoneNumber := ""

	resp, err := client.Api.ListAvailablePhoneNumberLocal("US", params)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		for record := range resp {
			if resp[record].FriendlyName != nil {
				phoneNumber = *resp[record].PhoneNumber
			} else {
				phoneNumber = *resp[record].PhoneNumber
			}
		}
	}
	return phoneNumber
}

func BuyNumber(phoneNumber string) string {

	err := godotenv.Load("../.env")

	if err != nil {
		log.Fatal("Error loading .env file for twilio")
	}

	os.Setenv("TWILIO_ACCOUNT_SID", os.Getenv("TWILIO_ACCOUNT_SID"))
	os.Setenv("TWILIO_AUTH_TOKEN", os.Getenv("TWILIO_AUTH_TOKEN"))

	client := twilio.NewRestClient()

	params := &api.CreateIncomingPhoneNumberParams{}
	params.SetPhoneNumber(phoneNumber)

	resp, err := client.Api.CreateIncomingPhoneNumber(params)

	responsePhoneNumber := ""

	if err != nil {
		fmt.Println(err.Error())
	} else {
		if resp.Sid != nil {
			responsePhoneNumber = *resp.PhoneNumber
		} else {
			responsePhoneNumber = *resp.PhoneNumber
		}
	}
	return responsePhoneNumber
}

func addNumberToDB(phoneNumber string, c *gin.Context) {
	// get the user id from the jwt
	org := c.MustGet("orgs").(model.Organization)

	var organization model.Organization
	config.DB.Where("id = ?", org.ID).First(&organization)

	// create a new phone number
	organization.TwilioNumber = phoneNumber
	organization.IsTwilioAuthed = true

	// save the phone number to the database
	config.DB.Save(&organization)
}

func RegisterOrgTwilioNumber(c *gin.Context) {
	// get the area code from the url params
	areaCode := c.Query("area_code")
	// convert the area code to an int
	areaCodeInt, err := strconv.Atoi(areaCode)

	if err != nil {
		fmt.Println(err)
	}

	phoneNumber := ListAvaliableNumbers(areaCodeInt)

	res := BuyNumber(phoneNumber)

	addNumberToDB(res, c)

	if phoneNumber == "" || res == " " {
		c.JSON(400, gin.H{
			"message": "No phone number available",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": res,
	})

}

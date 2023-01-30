package twilio

import (
	"github.com/joho/godotenv"
	"github.com/twilio/twilio-go"
	api "github.com/twilio/twilio-go/rest/api/v2010"
	"os"
)

type Response struct {
	IsDeliverable bool `json:"is_deliverable"`
}

func GetMessages(number string) (Response, error) {
	err := godotenv.Load("../.env")

	response := Response{
		IsDeliverable: false,
	}

	if err != nil {
		return response, err
	}

	err = os.Setenv("TWILIO_ACCOUNT_SID", os.Getenv("TWILIO_ACCOUNT_SID"))
	if err != nil {
		return response, err
	}
	err = os.Setenv("TWILIO_AUTH_TOKEN", os.Getenv("TWILIO_AUTH_TOKEN"))
	if err != nil {
		return response, err
	}

	client := twilio.NewRestClient()

	params := &api.ListMessageParams{}
	params.SetTo(number)

	resp, err := client.Api.ListMessage(params)

	if len(resp) == 0 {
		response.IsDeliverable = true
		return response, nil
	}

	if err != nil {
		return response, err
	} else {
		for _, message := range resp {
			if *message.Status == "delivered" {
				response.IsDeliverable = true
				return response, nil
			}
		}
	}
	return response, nil
}

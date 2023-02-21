package sendinblue

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/nickhansel/nucleus/config"
	"github.com/nickhansel/nucleus/model"
	"io"
	"net/http"
	"os"
	"time"
)

type EmailBody struct {
	Sender struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	} `json:"sender"`
	To []struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	} `json:"to"`
	Subject     string   `json:"subject"`
	HtmlContent string   `json:"htmlContent"`
	Tags        []string `json:"tags"`
}

func ScheduleFlowEmails(ids []int64, org model.Organization, emailBody string, emailSubject string) {
	err := godotenv.Load("../.env")
	if err != nil {
		return
	}

	var body EmailBody
	body.Sender.Name = org.Name
	body.Sender.Email = org.SendinblueEmail
	body.Subject = emailSubject
	body.HtmlContent = emailBody
	body.Tags = []string{"flow"}
	body.To = make([]struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	}, 0)

	apiKey := os.Getenv("SENDINBLUE_API_KEY")

	reqUrl := "https://api.sendinblue.com/v3/smtp/email"

	client := &http.Client{Timeout: 10 * time.Second}

	var sentEmails int32

	//lopp through the ids
	for _, id := range ids {
		var customer model.Customer
		config.DB.Where("\"id\" = ?", id).First(&customer)
		fmt.Println(customer.EmailAddress)

		if customer.EmailUnsubscribed == true || customer.EmailAddress == "" {
			continue
		}

		sentEmails++

		body.To = append(body.To, struct {
			Email string `json:"email"`
			Name  string `json:"name"`
		}{Email: customer.EmailAddress, Name: customer.GivenName + " " + customer.FamilyName})
	}
	// make the request
	req, err := http.NewRequest("POST", reqUrl, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	//add the body to the request
	reqBody, err := json.Marshal(body)
	if err != nil {
		fmt.Println(err)
		return
	}

	// set the body
	req.Body = io.NopCloser(bytes.NewBuffer(reqBody))

	// set the headers
	req.Header.Set("content-type", "application/json")
	req.Header.Set("accept", "application/json")
	req.Header.Set("api-key", apiKey)

	// make the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	// close the response body
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
	}(resp.Body)

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return
	}

	for _, id := range ids {
		var customer model.Customer
		config.DB.Where("\"id\" = ?", id).First(&customer)
		fmt.Println(customer.EmailAddress)

		if customer.EmailUnsubscribed == true || customer.EmailAddress == "" {
			continue
		}
		// time in 2023-02-12T15:29:19.818Z format in string
		customer.DatesReceivedEmail = append(customer.DatesReceivedEmail, time.Now().Format("2006-01-02 15:04:05"))
		config.DB.Save(&customer)
	}

	org.EmailCount += sentEmails
	config.DB.Save(&org)

	fmt.Println(result)
}

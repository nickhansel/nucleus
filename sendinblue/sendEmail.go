package sendinblue

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/nickhansel/nucleus/config"
	"github.com/nickhansel/nucleus/model"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Body struct {
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

func SendSendinblueEmail(c *gin.Context) {
	err := godotenv.Load("../.env")
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	apiKey := os.Getenv("SENDINBLUE_API_KEY")

	reqUrl := "https://api.sendinblue.com/v3/smtp/email"

	client := &http.Client{Timeout: 10 * time.Second}

	var body Body
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// make sure there are no values in the body that are empty
	if body.Sender.Name == "" || body.Sender.Email == "" || body.To == nil || body.Subject == "" || body.HtmlContent == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields!"})
		return
	}

	// make the request
	req, err := http.NewRequest("POST", reqUrl, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//add the body to the request
	reqBody, err := json.Marshal(body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// close the response body
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}(resp.Body)

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return
	}

	// return the response
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"body":   result,
	})
}

func SendScheduledEmails(EmailCampaign model.EmailCampaign, org model.Organization) {
	err := godotenv.Load("../.env")
	if err != nil {
		return
	}

	//batch update all the customers that are in the target emails and add the time right now to the DatesReceivedEmail field which is a string array
	var Customers []model.Customer

	//find the customers that are in the target emails
	//goroutine to run the two queries at the same time

	//run the two bottom functions at the same time
	go func() {
		var emailAddresses []string
		for _, v := range EmailCampaign.TargetEmails {
			emailAddresses = append(emailAddresses, fmt.Sprint(v))
		}

		config.DB.Where("email_address IN (?)", emailAddresses).Find(&Customers)

	}()

	go func() {
		for _, customer := range Customers {
			customer.DatesReceivedEmail = append(customer.DatesReceivedEmail, time.Now().Format("2006-01-02 15:04:05"))
			config.DB.Save(&customer)
		}
	}()
	
	//var emailAddresses []string
	//for _, v := range EmailCampaign.TargetEmails {
	//	emailAddresses = append(emailAddresses, fmt.Sprint(v))
	//}
	//
	//config.DB.Where("email_address IN (?)", emailAddresses).Find(&Customers)
	//for _, customer := range Customers {
	//	customer.DatesReceivedEmail = append(customer.DatesReceivedEmail, time.Now().Format("2006-01-02 15:04:05"))
	//	config.DB.Save(&customer)
	//}

	apiKey := os.Getenv("SENDINBLUE_API_KEY")

	reqUrl := "https://api.sendinblue.com/v3/smtp/email"

	client := &http.Client{Timeout: 10 * time.Second}

	organization := org
	targetEmails := EmailCampaign.TargetEmails

	//emails is a hashmap
	emails := make(map[string]string)

	for _, id := range targetEmails {
		var customer model.Customer
		config.DB.Where("email_address = ?", id).First(&customer)
		if customer.EmailAddress == "" {
			continue
		}
		emails[customer.EmailAddress] = customer.GivenName
	}

	//	loop through target emails and send the email
	for email, name := range emails {
		type Body struct {
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

		//convert the EmailCampaign.ID to a string that is the int32 value
		idToString := strconv.Itoa(int(EmailCampaign.ID))

		var body Body
		body.Sender.Name = organization.Name
		body.Sender.Email = organization.SendinblueEmail
		body.Subject = EmailCampaign.Subject
		body.HtmlContent = EmailCampaign.HTML
		body.Tags = []string{idToString}

		//append Name: name and Email: email to the To array
		body.To = append(body.To, struct {
			Email string `json:"email"`
			Name  string `json:"name"`
		}{Email: email, Name: name})
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

		var result map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&result)
		if err != nil {
			fmt.Println(err)
			return
		}

		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				fmt.Println(err)
				return
			}
		}(resp.Body)

		// return the response
		fmt.Println(result)
	}

}

package email

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/madflojo/tasks"
	"github.com/nickhansel/nucleus/config"
	"github.com/nickhansel/nucleus/model"
	"io"
	"net/http"
	"os"
	"time"
)

//TODO: Get email bounces from Sendinblue and update the database field is_email_deliverable https://developers.sendinblue.com/reference/gettransacblockedcontacts

func GetEmailBounces() (map[string]interface{}, error) {
	err := godotenv.Load("../.env")

	if err != nil {
		fmt.Println("Error loading .env file for sendinblue")
	}

	apiKey := os.Getenv("SENDINBLUE_API_KEY")

	reqUrl := "https://api.sendinblue.com/v3/smtp/blockedContacts?limit=50&offset=0&sort=desc"

	client := &http.Client{Timeout: 10 * time.Second}

	req, err := http.NewRequest("GET", reqUrl, nil)
	if err != nil {
		fmt.Println("Error creating request to sendinblue")
	}

	req.Header.Add("api-key", apiKey)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request to sendinblue")
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Error closing response body")
		}
	}(resp.Body)

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		fmt.Println("Error decoding response body")
	}

	return result, nil
}

func ScheduleGetEmailBounces() {
	scheduler := tasks.New()

	id, err := scheduler.Add(&tasks.Task{
		Interval: time.Duration(1) * time.Hour,
		RunOnce:  false,
		TaskFunc: func() error {
			bounces, err := GetEmailBounces()
			if err != nil {
				return err
			}

			var customer []model.Customer
			config.DB.Where("\"is_email_deliverable\" = ?", true).Find(&customer)
			for _, bounce := range bounces["contacts"].([]interface{}) {
				//check if the email is in the customer array
				for _, c := range customer {
					if c.EmailAddress == bounce.(map[string]interface{})["email"] {
						config.DB.Model(&model.Customer{}).Where("\"email_address\" = ?", bounce.(map[string]interface{})["email"]).Update("\"is_email_deliverable\"", false)
						fmt.Println("Updated user with email: ", bounce.(map[string]interface{})["email"])
					}
				}
			}

			return nil
		},
	})
	if err != nil {
		fmt.Println("Error scheduling task")
	}

	fmt.Println("Scheduled task with ID: ", id, " to run in ", 86400, " seconds")
	fmt.Println(id)
}

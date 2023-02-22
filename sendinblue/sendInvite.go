package sendinblue

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/nickhansel/nucleus/model"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

func SendInviteEmail(emailAddress string, org model.Organization, inviteId int64) {
	err := godotenv.Load("../.env")
	if err != nil {
		return
	}

	apiKey := os.Getenv("SENDINBLUE_API_KEY")

	reqUrl := "https://api.sendinblue.com/v3/smtp/email"

	client := &http.Client{Timeout: 10 * time.Second}

	inviteIdString := strconv.FormatInt(inviteId, 10)

	var body = Body{
		Sender: struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		}{
			Name:  "Nucleus",
			Email: org.SendinblueEmail,
		},
		To: []struct {
			Email string `json:"email"`
			Name  string `json:"name"`
		}{
			{
				Email: emailAddress,
				Name:  "Nucleus",
			},
		},
		Subject:     "Nucleus Invite",
		HtmlContent: "<div><p>You have been invited to join Nucleus!</p><br><p>Click the link below to join!</p><br><p>https://rerach.co/invite/" + inviteIdString + "</p></div>",
	}

	// make the request
	req, err := http.NewRequest("POST", reqUrl, nil)
	if err != nil {
		return
	}

	//add the body to the request
	reqBody, err := json.Marshal(body)
	if err != nil {
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
		return
	}

	// close the response body
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return
	}

	fmt.Println(result)
}

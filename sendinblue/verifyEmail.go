package sendinblue

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/nickhansel/nucleus/config"
	"github.com/nickhansel/nucleus/model"
	"io"
	"net/http"
	"os"
	"time"
)

type VerifyEmail struct {
	Email string `json:"email"`
}

func SendVerifyEmail(c *gin.Context) {
	org := c.MustGet("orgs").(model.Organization)

	if org.IsSendinblueAuthed == true {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Organization is already verified!"})
		return
	}

	err := godotenv.Load("../.env")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error loading .env file!"})
		return
	}

	apiKey := os.Getenv("SENDINBLUE_API_KEY")

	reqUrl := "https://api.sendinblue.com/v3/smtp/email"

	client := &http.Client{Timeout: 10 * time.Second}

	//generate random 6 digit verification code
	verificationCode := generateRandomString(6)

	var body VerifyEmail
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error binding JSON!"})
		return
	}

	if body.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields!"})
		return
	}

	var organization model.Organization
	config.DB.Where("id = ?", org.ID).First(&organization)
	organization.EmailVerificationCode = verificationCode
	organization.SendinblueEmail = body.Email

	config.DB.Save(&organization)

	var reqBody Body
	reqBody.Sender.Name = "Rereach"
	reqBody.Sender.Email = "no-reply@rereach.co"
	reqBody.To = append(reqBody.To, struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	}{Email: body.Email, Name: org.Name})
	reqBody.Subject = "Verify your email address"
	reqBody.HtmlContent = fmt.Sprintf("<p>Hi %s,</p><p>Thanks for signing up for Rereach! This is your verification code: %s</p><p>Thanks,</p><p>Rereach Team</p>", org.Name, verificationCode)

	fmt.Println(reqBody)
	req, err := http.NewRequest("POST", reqUrl, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error creating request!"})
		return
	}

	//bind the body to the request
	c.BindJSON(&reqBody)

	requestBody, err := json.Marshal(reqBody)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error marshalling JSON!"})
		return
	}
	// send the email
	req.Body = io.NopCloser(bytes.NewBufferString(string(requestBody)))

	// set the headers
	req.Header.Set("content-type", "application/json")
	req.Header.Set("accept", "application/json")
	req.Header.Set("api-key", apiKey)

	// make the request
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error making request!"})
		return
	}

	// close the response body
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error closing response body!"})
			return
		}
	}(resp.Body)

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return
	}

	fmt.Println(result)

	// return the response
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"body":   result,
	})
}

type Code struct {
	Code string `json:"code"`
}

func Verify(c *gin.Context) {
	org := c.MustGet("orgs").(model.Organization)

	var code Code
	if err := c.ShouldBindJSON(&code); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var organization model.Organization
	config.DB.Where("id = ?", org.ID).First(&organization)

	if organization.EmailVerificationCode == code.Code {
		organization.IsSendinblueAuthed = true
		config.DB.Save(&organization)
		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"body":   "Email verified!",
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"body":   "Invalid code!",
		})
	}

}

func generateRandomString(i int) (res string) {
	b := make([]byte, i)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println(err)
	}
	return base64.URLEncoding.EncodeToString(b)
}

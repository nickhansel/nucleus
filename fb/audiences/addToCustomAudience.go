package fb

// url: /fb/:orgId/audiences/:customer_group_id
// method to take in a customer group id and create or update a custom audience

import (
	"fmt"
	"io"
	"io/ioutil"

	"bytes"
	"crypto/sha256"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nickhansel/nucleus/config"
	"github.com/nickhansel/nucleus/model"
)

type Body struct {
	IDs []string `json:"ids"`
}

func UpdateCustomAudience(c *gin.Context) {
	customerGroup, org := GetCustomersFromGroup(c)

	customAudienceId := customerGroup.FbCustomAudienceID

	// check if the custom audience exists
	if customAudienceId == "" {
		// abort if the custom audience doesn't exist
		c.JSON(http.StatusBadRequest, gin.H{"error": "Custom audience not found!"})
	}
	var body Body
	err := c.BindJSON(&body)
	if err != nil {
		return
	}

	// get the customers from the customer group
	var customers []model.Customer
	// find customers where the id is in body.IDs
	config.DB.Where("id IN (?)", body.IDs).Find(&customers)
	fmt.Println(customers)
	// create the schema
	schema := []string{"FN", "LN", "EMAIL"}

	// create the data
	var data [][]string

	for _, customer := range customers {
		if customer.GivenName != "" && customer.FamilyName != "" {
			// hash the data
			customer.GivenName = fmt.Sprintf("%x", sha256.Sum256([]byte(customer.GivenName)))
			customer.FamilyName = fmt.Sprintf("%x", sha256.Sum256([]byte(customer.FamilyName)))
			customer.EmailAddress = fmt.Sprintf("%x", sha256.Sum256([]byte(customer.EmailAddress)))
			data = append(data, []string{customer.GivenName, customer.FamilyName, customer.EmailAddress})
		}
	}

	payload := CustomAudiencePayload{
		Schema: schema,
		Data:   data,
	}

	// convert the payload to json
	payloadJson, err := json.Marshal(payload)

	url := "https://graph.facebook.com/v15.0/" + customAudienceId + "/users"

	// make the request
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// add payloadJson as body to the request
	req.Body = ioutil.NopCloser(bytes.NewBuffer(payloadJson))

	q := req.URL.Query()
	// q.Add("payload", fmt.Sprintf(`{"schema": %v, "data": %v}`, schema, data))
	// add payloadJson as body to the request
	q.Add("access_token", org.FbAccessToken)
	q.Add("payload", string(payloadJson))
	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

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

	c.JSON(http.StatusOK, gin.H{"result": result})

}

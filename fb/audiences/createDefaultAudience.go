package fb

// url: /fb/:orgId/audiences/:customer_group_id
// method to take in a customer group id and create or update a custom audience

import (
	"fmt"
	"io/ioutil"

	"bytes"
	"crypto/sha256"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nickhansel/nucleus/config"
	"github.com/nickhansel/nucleus/model"
)

func GetCustomersFromGroup(c *gin.Context) (resp model.CustomerGroup, orgresp model.Organization) {
	org := c.MustGet("orgs").(model.Organization)

	// get the customer group id from the url
	customerGroupId := c.Param("customer_group_id")

	// get the customer group from the database
	customerGroup := model.CustomerGroup{}
	config.DB.First(&customerGroup, customerGroupId)

	// check if the customer group exists
	if customerGroup.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	var customersToCustomerGroups []model.CustomersToCustomerGroups
	// preload the customers that are in the customer group
	config.DB.Preload("Customer").Find(&customersToCustomerGroups)

	for _, customerToCustomerGroup := range customersToCustomerGroups {
		if customerToCustomerGroup.B == customerGroup.ID {
			customerGroup.Customers = append(customerGroup.Customers, customerToCustomerGroup.Customer)
		}
	}

	if len(customerGroup.Customers) == 0 {
		customerGroup.Customers = []model.Customer{}
	}

	// get the customer group's customers
	return customerGroup, org

}

type CustomAudiencePayload struct {
	Schema []string   `json:"schema"`
	Data   [][]string `json:"data"`
}

func CreateCustomAudience(c *gin.Context) {
	customerGroup, org := GetCustomersFromGroup(c)

	customAudienceId := customerGroup.FbCustomAudienceID

	// check if the custom audience exists
	if customAudienceId == "" {
		// abort if the custom audience doesn't exist
		c.JSON(http.StatusBadRequest, gin.H{"error": "Custom audience not found!"})
	}

	// get the customers from the customer group
	customers := customerGroup.Customers

	// create the schema
	schema := []string{"FN", "LN", "EMAIL"}

	// create the data
	data := [][]string{}

	for _, customer := range customers {
		if customer.GivenName != "" && customer.FamilyName != "" && customer.EmailAddress != "" {
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

	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	c.JSON(http.StatusOK, gin.H{"result": result})

}

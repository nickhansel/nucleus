package customers

import (
	"net/http"

	"fmt"
	"strconv"

	"github.com/nickhansel/nucleus/config"
	"github.com/nickhansel/nucleus/model"

	"github.com/gin-gonic/gin"
)

type CustomerResponse struct {
	ID          int64   `json:"id"`
	PurchaseID  string  `json:"purchase_id"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
	AmountMoney float64 `json:"amount_money"`
	Currency    string  `json:"currency"`
	Status      string  `json:"status"`
	SourceType  string  `json:"source_type"`
	LocationID  string  `json:"location_id"`
	ProductType string  `json:"product_type"`
	Name        string  `json:"name"`
	Quantity    int64   `json:"quantity"`
	Integration string  `json:"integration"`
}

func GetCustomers(c *gin.Context) {
	orgId := c.Param("orgId")

	orgIdInt, err := strconv.Atoi(orgId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid URL parameter",
		})
		return
	}

	// query the db and find the customers where organizationId = orgId
	customers := []model.Customer{}
	res := config.DB.Where(&model.Customer{OrganizationID: int64(orgIdInt)}).Find(&customers).Error

	if res != nil {
		fmt.Println("res is not nil")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": res.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"customers": customers,
	})
}

func GetCustomerById(c *gin.Context) {
	org := c.MustGet("orgs").(model.Organization)
	customerId := c.Param("customerId")

	customerIdInt, err := strconv.Atoi(customerId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid URL parameter",
		})
		return
	}

	customer := model.Customer{}
	res := config.DB.Where(&model.Customer{OrganizationID: org.ID, ID: int64(customerIdInt)}).First(&customer).Error
	if res != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": res.Error(),
		})
		return
	}

	//query := "SELECT * FROM \"Purchase\" INNER JOIN \"purchased_item\" ON \"Purchase\".purchase_id = purchased_item.\"purchaseId\" WHERE \"customerId\" = 838211150098235393"
	//VariationsName string  `json:"name"`
	//
	//customerRespoonse can be any struct
	var customerResponse []CustomerResponse
	config.DB.Table("\"Purchase\"").Select("*").Joins("INNER JOIN \"purchased_item\" ON \"Purchase\".purchase_id = purchased_item.\"purchaseId\"").
		Where("\"customerId\" = ?", customer.ID).Scan(&customerResponse)

	c.JSON(http.StatusOK, gin.H{
		"customer":  customer,
		"purchases": customerResponse,
	})

}

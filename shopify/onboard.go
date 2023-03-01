package shopify

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nickhansel/nucleus/config"
	"github.com/nickhansel/nucleus/model"
	"net/http"
	"strconv"
)

func GetInitialData(c *gin.Context) {
	org := c.MustGet("orgs").(model.Organization)

	//	get items from shopify
	var orgPos []model.Pos
	config.DB.Where("\"organizationId\" = ?", org.ID).Find(&orgPos)

	var shopifyPos model.Pos

	//	find the pos where the name is Shopify
	for _, pos := range orgPos {
		if pos.Name == "Shopify" {
			shopifyPos = pos
		}
	}

	err := GetItems(org, shopifyPos)
	if err != nil {
		return
	}

	err = GetCustomers(org, shopifyPos)
	if err != nil {
		return
	}

	err = GetLocations(org, shopifyPos)
	if err != nil {
		return
	}

	err = GetOrders(org, shopifyPos)
	if err != nil {
		return
	}

	var location model.StoreLocation
	config.DB.Where("\"type\" = ? AND \"organizationId\" = ?", "Shopify", org.ID).First(&location)

	err = RegisterHooks(shopifyPos, location.BusinessName, org)
	if err != nil {
		return
	}

	//	get all of the items from the database
	var customers []model.Customer
	config.DB.Where("\"organizationId\" = ? AND \"pos_name\" = ?", org.ID, "Shopify").Find(&customers)

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}

func GetLocations(org model.Organization, shopifyPos model.Pos) error {
	accessToken := shopifyPos.AccessToken
	storeUrl := org.ShopifyUrl

	reqUrl := fmt.Sprintf("%s/admin/api/2022-07/locations.json", storeUrl)

	req, err := http.NewRequest("GET", reqUrl, nil)
	if err != nil {
		return err
	}

	req.Header.Add("X-Shopify-Access-Token", accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return err
	}

	if result["locations"] == nil {
		newErr := fmt.Errorf("error in getting locations")
		return newErr
	}

	locations := result["locations"].([]interface{})

	for _, location := range locations {
		locationMap := location.(map[string]interface{})

		var newLocation model.StoreLocation
		newLocation.Name = locationMap["name"].(string)
		newLocation.PosID = strconv.FormatInt(shopifyPos.ID, 10)
		newLocation.Type = "Shopify Physical"
		if locationMap["address1"] != nil {
			newLocation.AddressLine1 = locationMap["address1"].(string)
		} else {
			newLocation.AddressLine1 = ""
		}

		if locationMap["city"] != nil {
			newLocation.Locality = locationMap["city"].(string)
		} else {
			newLocation.Locality = ""
		}

		if locationMap["country"] != nil {
			if newLocation.Country == "United States" {
				newLocation.Country = "US"
				newLocation.Currency = "USD"
			}
		} else {
			newLocation.Country = ""
		}

		if locationMap["province"] != nil {
			newLocation.AdministrativeDistrictLevel1 = locationMap["province"].(string)
		} else {
			newLocation.AdministrativeDistrictLevel1 = ""
		}

		if locationMap["zip"] != nil {
			newLocation.PostalCode = locationMap["zip"].(string)
		} else {
			newLocation.PostalCode = ""
		}

		newLocation.OrganizationID = org.ID
		if locationMap["active"].(bool) {
			newLocation.Status = "ACTIVE"
		} else {
			newLocation.Status = "INACTIVE"
		}

		newLocation.CreatedAt = locationMap["created_at"].(string)
		newLocation.MerchantID = strconv.FormatInt(int64(locationMap["id"].(float64)), 10)
		newLocation.LanguageCode = "en-US"
		if locationMap["phone"] != nil {
			newLocation.PhoneNumber = locationMap["phone"].(string)
		} else {
			newLocation.PhoneNumber = ""
		}

		newLocation.BusinessName = locationMap["name"].(string)
		newLocation.Timezone = "None"
		newLocation.OrganizationID = org.ID

		config.DB.Create(&newLocation)
		fmt.Println("create new location:", newLocation.Name)
	}
	return nil
}

func GetOrders(org model.Organization, shopifyPos model.Pos) error {
	accessToken := shopifyPos.AccessToken
	storeUrl := org.ShopifyUrl

	reqUrl := fmt.Sprintf("%s/admin/api/2022-07/orders.json?status=any", storeUrl)

	req, err := http.NewRequest("GET", reqUrl, nil)
	if err != nil {
		return err
	}

	req.Header.Add("X-Shopify-Access-Token", accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return err
	}

	if result["orders"] == nil {
		newErr := fmt.Errorf("error in getting orders")
		return newErr
	}

	orders := result["orders"].([]interface{})

	for _, order := range orders {
		var newPurchase model.Purchase
		newPurchase.CreatedAt = order.(map[string]interface{})["created_at"].(string)
		newPurchase.PurchaseID = strconv.FormatFloat(order.(map[string]interface{})["id"].(float64), 'f', 0, 64)
		newPurchase.UpdatedAt = order.(map[string]interface{})["created_at"].(string)
		// convert the total price from string to float64
		newPurchase.AmountMoney, err = strconv.ParseFloat(order.(map[string]interface{})["current_subtotal_price"].(string), 64)
		if err != nil {
			newPurchase.AmountMoney = 0
		}
		newPurchase.Currency = order.(map[string]interface{})["currency"].(string)
		if order.(map[string]interface{})["financial_status"].(string) == "paid" {
			newPurchase.Status = "COMPLETED"
		} else {
			newPurchase.Status = "OPEN"
		}

		newPurchase.SourceType = "PURCHASE"
		if order.(map[string]interface{})["location_id"] != nil {
			var storeLocation model.StoreLocation
			// convert the location id from float64 to string
			convertedID := strconv.FormatFloat(order.(map[string]interface{})["location_id"].(float64), 'f', 0, 64)
			config.DB.Where("\"merchant_id\" = ?", convertedID).First(&storeLocation)
			fmt.Println("store location:", storeLocation.ID, convertedID)
			newPurchase.LocationID = storeLocation.ID
			newPurchase.Location_ID = convertedID
		} else {
			var shopifyLocation model.StoreLocation
			config.DB.Where("\"organizationId\" = ? AND \"type\" = ?", org.ID, "Shopify").First(&shopifyLocation)
			fmt.Println("shopify location:", shopifyLocation.ID)
			newPurchase.LocationID = shopifyLocation.ID
			newPurchase.Location_ID = shopifyLocation.MerchantID
		}

		fmt.Println("new purchase location:", newPurchase.LocationID)

		newPurchase.OrganizationID = org.ID

		if order.(map[string]interface{})["customer"] != nil {
			customerId := order.(map[string]interface{})["customer"].(map[string]interface{})["id"].(float64)
			var dbCustomer model.Customer
			config.DB.Where("\"pos_id\" = ?", strconv.FormatFloat(customerId, 'f', 0, 64)).First(&dbCustomer)
			newPurchase.CustomerID = dbCustomer.ID
		}

		if newPurchase.CustomerID != 0 {
			config.DB.Omit("attributedCampaignId", "itemsId").Create(&newPurchase)
		} else {
			config.DB.Omit("attributedCampaignId", "itemsId", "customerId").Create(&newPurchase)
		}

		if order.(map[string]interface{})["line_items"] != nil {
			lineItems := order.(map[string]interface{})["line_items"].([]interface{})
			for _, lineItem := range lineItems {
				if lineItem.(map[string]interface{})["product_exists"] != false {
					var purchasedItem model.PurchasedItem
					if lineItem.(map[string]interface{})["variant_id"] != nil && lineItem.(map[string]interface{})["variant_title"] != nil {
						purchasedItem.IsVaration = true
						var variation model.Variation
						variationID := strconv.FormatFloat(lineItem.(map[string]interface{})["variant_id"].(float64), 'f', 0, 64)
						config.DB.Where("\"varation_item_id\" = ?", variationID).First(&variation)
						fmt.Println("variation id:", variation.ID, variationID)
						purchasedItem.VariationID = variation.ID
					} else {
						purchasedItem.IsVaration = false
						itemId := strconv.FormatFloat(lineItem.(map[string]interface{})["product_id"].(float64), 'f', 0, 64)
						var item model.Item
						config.DB.Where("\"item_id\" = ?", itemId).First(&item)
						purchasedItem.ItemID = item.ID
					}
					purchasedItem.Quantity = int32(int64(lineItem.(map[string]interface{})["quantity"].(float64)))
					purchasedItem.Name = lineItem.(map[string]interface{})["name"].(string)
					//	convert price from string to float64
					purchasedItem.Cost, err = strconv.ParseFloat(lineItem.(map[string]interface{})["price"].(string), 64)
					if err != nil {
						purchasedItem.Cost = 0
					}
					purchasedItem.PurchaseID = newPurchase.PurchaseID

					if purchasedItem.IsVaration {
						config.DB.Omit("itemId").Create(&purchasedItem)
					} else {
						config.DB.Omit("variationId").Create(&purchasedItem)
					}
				}
			}
		}

	}
	return nil
}

func GetCustomers(org model.Organization, shopifyPos model.Pos) error {
	accessToken := shopifyPos.AccessToken
	storeUrl := org.ShopifyUrl

	reqUrl := fmt.Sprintf("%s/admin/api/2022-04/customers.json", storeUrl)

	req, err := http.NewRequest("GET", reqUrl, nil)
	if err != nil {
		return err
	}

	req.Header.Add("X-Shopify-Access-Token", accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return err
	}

	if result["customers"] == nil {
		newErr := fmt.Errorf("error in getting customers")
		return newErr
	}

	customers := result["customers"].([]interface{})

	var customerGroup model.CustomerGroup
	config.DB.Where("\"organizationId\" = ? AND \"name\" = ?", org.ID, "Default Group").Find(&customerGroup)

	var shopifyCustomerGroup model.CustomerGroup
	shopifyCustomerGroup.Name = "Shopify customers"
	shopifyCustomerGroup.OrganizationID = org.ID
	config.DB.Create(&shopifyCustomerGroup)

	if customerGroup.ID == 0 {
		newErr := fmt.Errorf("no customer group found for org %s", org.Name)
		return newErr
	}

	for _, customer := range customers {
		customerMap := customer.(map[string]interface{})
		var newCustomer model.Customer

		newCustomer.PosID = strconv.FormatInt(int64(customerMap["id"].(float64)), 10)
		newCustomer.CreatedAt = customerMap["created_at"].(string)
		newCustomer.UpdatedAt = customerMap["updated_at"].(string)
		newCustomer.GivenName = customerMap["first_name"].(string)
		newCustomer.FamilyName = customerMap["last_name"].(string)
		newCustomer.EmailAddress = customerMap["email"].(string)
		if customerMap["phone"] != nil {
			newCustomer.PhoneNumber = customerMap["phone"].(string)
		} else {
			newCustomer.PhoneNumber = ""
		}
		if customerMap["note"] != nil {
			newCustomer.Note = customerMap["note"].(string)
		} else {
			newCustomer.Note = ""
		}

		newCustomer.CreationSource = "FIRST_PARTY"

		if customerMap["addresses"] != nil {
			addresses := customerMap["addresses"].([]interface{})
			for _, address := range addresses {
				addressMap := address.(map[string]interface{})
				if addressMap["address1"] != nil {
					newCustomer.AddressLine1 = addressMap["address1"].(string)
				} else {
					newCustomer.AddressLine1 = ""
				}

				if addressMap["address2"] != nil {
					newCustomer.AddressLine2 = addressMap["address2"].(string)
				} else {
					newCustomer.AddressLine2 = ""
				}

				if addressMap["city"] != nil {
					newCustomer.Locality = addressMap["city"].(string)
				} else {
					newCustomer.Locality = ""
				}

				if addressMap["province"] != nil {
					newCustomer.AdministrativeDistrictLevel1 = addressMap["province"].(string)
				} else {
					newCustomer.AdministrativeDistrictLevel1 = ""
				}

				if addressMap["country"] != nil {
					newCustomer.Country = addressMap["country"].(string)
				} else {
					newCustomer.Country = ""
				}

				if addressMap["zip"] != nil {
					newCustomer.PostalCode = addressMap["zip"].(string)
				} else {
					newCustomer.PostalCode = ""
				}
			}
		}

		newCustomer.OrganizationID = org.ID
		newCustomer.PosName = "Shopify"
		// convert total spent to float
		totalSpent, err := strconv.ParseFloat(customerMap["total_spent"].(string), 64)
		if err != nil {
			totalSpent = 0
		}
		newCustomer.TotalSpent = totalSpent
		newCustomer.TotalPurhcases = int32(customerMap["orders_count"].(float64))

		var emailMarketingConsentState map[string]interface{}
		if customerMap["email_marketing_consent"] != nil {
			emailMarketingConsentState = customerMap["email_marketing_consent"].(map[string]interface{})
			newCustomer.EmailUnsubscribed = emailMarketingConsentState["state"].(string) == "subscribed"
		} else {
			newCustomer.EmailUnsubscribed = true
		}

		var smsMarketingConsentState map[string]interface{}

		if customerMap["sms_marketing_consent"] != nil {
			smsMarketingConsentState = customerMap["sms_marketing_consent"].(map[string]interface{})
			newCustomer.SmsUnsubscribed = smsMarketingConsentState["state"].(string) == "not_subscribed"
		} else {
			newCustomer.SmsUnsubscribed = true
		}

		config.DB.Create(&newCustomer)
		fmt.Println("customer created: ", newCustomer.GivenName)

		var customerToCustomerGroup model.CustomersToCustomerGroups
		customerToCustomerGroup.A = newCustomer.ID
		customerToCustomerGroup.B = customerGroup.ID
		config.DB.Create(&customerToCustomerGroup)
		customerToCustomerGroup.A = newCustomer.ID
		customerToCustomerGroup.B = shopifyCustomerGroup.ID
		config.DB.Create(&customerToCustomerGroup)
		fmt.Println("customer to customer group created: ", customerToCustomerGroup.A)

	}
	return nil
}

func GetItems(org model.Organization, shopifyPos model.Pos) error {
	var shopifyLocation model.StoreLocation
	config.DB.Where("\"pos_id\" = ?", strconv.FormatInt(shopifyPos.ID, 10)).Find(&shopifyLocation)

	if shopifyLocation.ID == 0 {
		newErr := fmt.Errorf("no shopify location found for org %s", org.Name)
		return newErr
	}

	accessToken := shopifyPos.AccessToken
	storeUrl := org.ShopifyUrl

	reqUrl := fmt.Sprintf("%s/admin/api/2022-07/products.json", storeUrl)

	req, err := http.NewRequest("GET", reqUrl, nil)
	if err != nil {
		return err
	}

	req.Header.Add("X-Shopify-Access-Token", accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return err
	}

	if result["products"] == nil {
		return err
	}

	products := result["products"].([]interface{})

	for _, product := range products {
		var item model.Item
		item.Type = product.(map[string]interface{})["product_type"].(string)
		item.ItemID = strconv.Itoa(int(product.(map[string]interface{})["id"].(float64)))
		item.CreatedAt = product.(map[string]interface{})["created_at"].(string)
		item.UpdatedAt = product.(map[string]interface{})["updated_at"].(string)
		item.PresentAtAllLocations = true
		item.IsDeleted = product.(map[string]interface{})["status"] != "active"
		item.Name = product.(map[string]interface{})["title"].(string)
		if product.(map[string]interface{})["body_html"] != nil {
			item.Description = product.(map[string]interface{})["body_html"].(string)
		} else {
			item.Description = ""
		}
		item.ProductType = "SHOPIFY-ITEM"
		item.OrganizationID = org.ID
		config.DB.Create(&item)
		fmt.Println("created item: ", item.Name)

		if product.(map[string]interface{})["variants"] != nil {
			for _, variant := range product.(map[string]interface{})["variants"].([]interface{}) {
				var newVariant model.Variation
				newVariant.VarationItemID = strconv.Itoa(int(variant.(map[string]interface{})["id"].(float64)))
				newVariant.IsDeleted = item.IsDeleted
				newVariant.Name = variant.(map[string]interface{})["title"].(string)
				//	convert price from string to float
				price, err := strconv.ParseFloat(variant.(map[string]interface{})["price"].(string), 64)
				if err != nil {
					return err
				}
				newVariant.Price = price
				newVariant.Currency = shopifyLocation.Currency
				newVariant.ItemID = item.ID

				config.DB.Create(&newVariant)
				fmt.Println("created variant: ", newVariant.Name)
			}
		}
	}

	return nil
}

package segmentQL

import (
	"fmt"
	"github.com/nickhansel/nucleus/config"
	"github.com/nickhansel/nucleus/model"
	"time"
)

func ParseSegmentQL(item int64, orgId int64, startDate string, endDate string, minPurchasePrice float64, maxPurchasePrice float64) []int64 {
	// find all customers who have a purchase with PurchasedItems.Variation.ID = item
	// return a slice of customer IDs
	// get all purchases
	var purchases []model.Purchase
	err := config.DB.Preload("PurchasedItems").Preload("PurchasedItems.Variation").Preload("Customer").Find(&purchases).Error
	// find the customers where organizationId = orgId
	for _, purchase := range purchases {
		customer := purchase.Customer
		if customer.OrganizationID != orgId {
			purchases = append(purchases[:0], purchases[1:]...)
		}
	}
	
	if err != nil {
		fmt.Println(err)
	}

	var customers []model.Customer
	//find the customers who made a purchase
	for _, purchase := range purchases {
		customer := purchase.Customer
		customers = append(customers, customer)
	}

	fmt.Println("Customers: ", len(customers))
	// get all customers who have a purchase with PurchasedItems.Variation.ID = item
	var customerIDs []int64
	if item != 0 {
		//	find customers who purchased item
		for _, purchase := range purchases {
			// remove the purchase where purchasedItems.Variation.ID != item
			for _, purchasedItem := range purchase.PurchasedItems {
				if purchasedItem.VariationID != item {
					//	remove all purchases where they share a customerID
					for _, purchased := range purchases {
						if purchased.CustomerID == purchase.CustomerID {
							purchases = append(purchases[:0], purchases[1:]...)
						}
					}
				} else {
					customers = append(customers, purchase.Customer)
				}
			}
		}
	}

	startDateCustomers := make([]model.Customer, 0)
	if startDate != "" && endDate != "" {
		//	find customers who purchased between startDate and endDate
		for _, purchase := range purchases {
			//convert createdAt to date
			// convert 2023-01-08T16:25:26.263Z to time
			createdDate, err := time.Parse(time.RFC3339, purchase.CreatedAt)
			if err != nil {
				fmt.Println(err)
			}
			startDateParsed, err := time.Parse("2006-01-02", startDate)
			if err != nil {
				fmt.Println(err)
			}
			endDateParsed, err := time.Parse("2006-01-02", endDate)
			if err != nil {
				fmt.Println(err)
			}
			if createdDate.After(startDateParsed) && createdDate.Before(endDateParsed) {
				for _, customer := range customers {
					if customer.ID == purchase.CustomerID {
						if purchase.Customer.TotalSpent > minPurchasePrice {
							startDateCustomers = append(startDateCustomers, purchase.Customer)
						}
					}
				}
			} else {
				for _, purchased := range purchases {
					if purchased.CustomerID == purchase.CustomerID {
						purchases = append(purchases[:0], purchases[1:]...)
					}
				}
			}
		}
		customers = startDateCustomers
	}

	newRangeCustomers := make([]model.Customer, 0)
	if minPurchasePrice != 0 || maxPurchasePrice != 0 {
		//	find customers who purchased between minPurchasePrice and maxPurchasePrice
		for _, purchase := range purchases {
			if purchase.Customer.TotalSpent > minPurchasePrice && purchase.Customer.TotalSpent < maxPurchasePrice {
				for _, customer := range customers {
					if customer.ID == purchase.CustomerID {
						if purchase.Customer.TotalSpent > minPurchasePrice {
							newRangeCustomers = append(newRangeCustomers, purchase.Customer)
						}
					}
				}
			} else {
				continue
			}
		}
		customers = newRangeCustomers
	}

	//print of all of the purchase AmountMoney
	for _, customer := range customers {
		if !Contains(customerIDs, customer.ID) && customer.ID != 0 {
			customerIDs = append(customerIDs, customer.ID)
		}
	}

	return customerIDs
}

func Contains(s []int64, e int64) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

//func FindCustomersWhoSpendMoreThan(amount float64, orgId int32) []int32 {
//
//}

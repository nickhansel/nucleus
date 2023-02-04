package groups

import (
	"fmt"
	"github.com/nickhansel/nucleus/config"
	"github.com/nickhansel/nucleus/model"
	"sort"
	"time"
)

func GetTopCustomers(orgId int32) {
	var customers []model.Customer

	config.DB.Where("\"organizationId\" = ?", orgId).Find(&customers)

	if len(customers) == 0 {
		return
	}

	totalCustomers := len(customers)
	topCustomers := float64(totalCustomers) * 0.2

	// sort the customers by total spent
	sort.Slice(customers, func(i, j int) bool {
		return customers[i].TotalSpent > customers[j].TotalSpent
	})

	// get the top 20% of customers
	topSpenders := customers[:int(topCustomers)]

	// remove customers where the total spent is 0
	for i, customer := range topSpenders {
		if customer.TotalSpent == 0 {
			topSpenders = topSpenders[:i]
			break
		}
	}

	// check if the customer group already exists
	var checkCustomerGroup model.CustomerGroup
	config.DB.Where("\"organizationId\" = ? AND \"name\" = ?", orgId, "Top 20% of customers").First(&checkCustomerGroup)

	var customerGroup model.CustomerGroup
	if checkCustomerGroup.ID != 0 {
		//	update the customer group
		checkCustomerGroup.UpdatedAt = time.Now()
		config.DB.Save(&checkCustomerGroup)

		//	check if the customers are already in the customer group and if they are not, add them
		// if there are customers in the customer group that are not in the topSpenders array, remove them
		for _, customer := range topSpenders {
			var checkCustomerToCustomerGroup model.CustomersToCustomerGroups
			config.DB.Where("\"A\" = ? AND \"B\" = ?", customer.ID, checkCustomerGroup.ID).First(&checkCustomerToCustomerGroup)

			if checkCustomerToCustomerGroup.A == 0 || checkCustomerToCustomerGroup.B == 0 {
				// add the customer to the customer group
				var CustomersToCustomerGroups model.CustomersToCustomerGroups

				CustomersToCustomerGroups.A = customer.ID
				CustomersToCustomerGroups.B = checkCustomerGroup.ID

				fmt.Println("Adding customer to customer group:", customer.ID)
				config.DB.Create(&CustomersToCustomerGroups)
			}
		}

		// if there are customers in the topSpenders array that are not in the customer group, add them
		var CustomersToCustomerGroups []model.CustomersToCustomerGroups
		config.DB.Where("\"B\" = ?", checkCustomerGroup.ID).Find(&CustomersToCustomerGroups)

		for _, customerToCustomerGroup := range CustomersToCustomerGroups {
			if customerToCustomerGroup.B == checkCustomerGroup.ID {
				var found bool
				for _, customer := range topSpenders {
					if customerToCustomerGroup.A == customer.ID {
						found = true
						break
					}
				}

				if !found {
					fmt.Println("Removing customer from customer group:", customerToCustomerGroup.A)
					config.DB.Where("\"A\" = ? AND \"B\" = ?", customerToCustomerGroup.A, customerToCustomerGroup.B).Delete(&model.CustomersToCustomerGroups{})
				}
			}
		}
	} else {
		// create a customer group and connect all of the customers that have organizationId = 19
		customerGroup.Name = "Top 20% of customers"
		customerGroup.OrganizationID = orgId
		customerGroup.CreatedAt = time.Now()
		customerGroup.UpdatedAt = time.Now()
		config.DB.Create(&customerGroup)

		// add the customers to the customer group
		for _, customer := range topSpenders {
			// add the customer to the Customers field of the customer group and connect them
			var CustomersToCustomerGroups model.CustomersToCustomerGroups

			CustomersToCustomerGroups.A = customer.ID
			CustomersToCustomerGroups.B = customerGroup.ID

			config.DB.Create(&CustomersToCustomerGroups)
		}
	}

	fmt.Println("Total customers:", totalCustomers)
	fmt.Println("Top spenders:", len(topSpenders))
	fmt.Println("Top 20% of customers created")

}

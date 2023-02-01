package text

import (
	"fmt"
	"github.com/madflojo/tasks"
	"github.com/nickhansel/nucleus/config"
	"github.com/nickhansel/nucleus/model"
	twilio "github.com/nickhansel/nucleus/twilio"
	"time"
)

//TODO: Get text bounces from Twilio and update the database field is_sms_deliverable https://www.twilio.com/docs/sms/tutorials/how-to-confirm-delivery-python

// GetMessages
func ScheduleGetTextBounces() {
	scheduler := tasks.New()

	id, err := scheduler.Add(&tasks.Task{
		//interval is every 1.3 hours
		Interval: time.Duration(1.3 * float64(time.Hour)),
		RunOnce:  false,
		TaskFunc: func() error {
			var customers []model.Customer
			config.DB.Where("\"is_sms_deliverable\" = ?", true).Find(&customers)
			//find customers where is_sms_deliverable is true
			for _, customer := range customers {
				isDeliverable, err := twilio.GetMessages(customer.PhoneNumber)
				if err != nil {
					fmt.Println("Error getting messages from twilio")
					return err
				}
				if !isDeliverable.IsDeliverable {
					fmt.Println("Customer: ", customer.ID, " is not deliverable")
					customer.IsSMSDeliverable = false
					config.DB.Save(&customer)
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

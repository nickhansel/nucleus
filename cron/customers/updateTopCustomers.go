package customers

import (
	"github.com/madflojo/tasks"
	"github.com/nickhansel/nucleus/api/customers/groups"
	"github.com/nickhansel/nucleus/config"
	"github.com/nickhansel/nucleus/model"
	"time"
)

func GetTopCustomersJob() {
	// GetTopCustomers(1)
	scheduler := tasks.New()
	// run every 2.3 hours
	duration := time.Duration(1)*time.Hour + time.Duration(23)*time.Minute

	_, err := scheduler.Add(&tasks.Task{
		Interval: duration,
		RunOnce:  false,
		TaskFunc: func() error {
			var customerGroups []model.CustomerGroup
			config.DB.Where("\name\" = ?", "Top 20% of customers").Find(&customerGroups)

			for _, customerGroup := range customerGroups {
				groups.GetTopCustomers(customerGroup.ID)
			}

			return nil
		},
	})
	if err != nil {
		return
	}
}

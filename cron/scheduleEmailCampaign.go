package cron

import (
	"time"

	"fmt"

	"github.com/madflojo/tasks"

	model "github.com/nickhansel/nucleus/model"
	sendinblue "github.com/nickhansel/nucleus/sendinblue"
)

func ScheduleEmailTasks(Date string, EmailCampaign model.EmailCampaign, org model.Organization) {
	scheduler := tasks.New()

	// print scheduled tasks

	howMany := secondsFromNow(Date)

	id, err := scheduler.Add(&tasks.Task{
		Interval: time.Duration(howMany) * time.Second,
		RunOnce:  true,
		TaskFunc: func() error {
			sendinblue.SendScheduledEmails(EmailCampaign, org)
			return nil
		},
	})
	if err != nil {
		fmt.Println(err)
		fmt.Println("Error scheduling task")
	}
	fmt.Println("Scheduled task with ID: ", id, " to run in ", howMany, " seconds at ", Date)
	fmt.Println(id)
}

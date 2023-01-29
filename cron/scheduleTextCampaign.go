package cron

import (
	"time"

	"fmt"

	"github.com/madflojo/tasks"

	model "github.com/nickhansel/nucleus/model"
	twilio "github.com/nickhansel/nucleus/twilio"
)

func secondsFromNow(dateString string) int {
	layout := "2006-01-02 15:04:05"
	t, _ := time.ParseInLocation(layout, dateString, time.Local)
	fmt.Println(t, t.Sub(time.Now()).Seconds())
	return int(t.Sub(time.Now()).Seconds())
}

func ScheduleTextTask(Date string, TextCampaign model.TextCampaign, org model.Organization) {
	scheduler := tasks.New()

	howMany := secondsFromNow(Date)

	id, err := scheduler.Add(&tasks.Task{
		Interval: time.Duration(howMany) * time.Second,
		RunOnce:  true,
		TaskFunc: func() error {
			twilio.SendScheduledTexts(TextCampaign, org)
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

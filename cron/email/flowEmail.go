package email

import (
	"fmt"
	"github.com/madflojo/tasks"
	"github.com/nickhansel/nucleus/model"
	"github.com/nickhansel/nucleus/sendinblue"
	"time"
)

func secondsFromNowUTC(dateString string) int {
	layout := "2006-01-02 15:04:05"
	t, _ := time.Parse(layout, dateString)
	fmt.Println(t, t.Sub(time.Now()).Seconds())
	return int(t.Sub(time.Now()).Seconds())
}

func ScheduleFlowEmails(Date string, ids []int64, org model.Organization, emailBody string, emailSubject string) {
	scheduler := tasks.New()

	if Date == "" {
		Date = time.Now().Add(time.Second * 30).Format("2006-01-02 15:04:05")
	}

	howMany := secondsFromNowUTC(Date)

	id, err := scheduler.Add(&tasks.Task{
		Interval: time.Duration(howMany) * time.Second,
		RunOnce:  true,
		TaskFunc: func() error {
			sendinblue.ScheduleFlowEmails(ids, org, emailBody, emailSubject)
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

package email

import (
	"fmt"
	"github.com/madflojo/tasks"
	"github.com/nickhansel/nucleus/config"
	"github.com/nickhansel/nucleus/model"
	"time"
)

type EmailCampaignAnalytics struct {
	Sent            int    `json:"sent"`
	Delivered       int    `json:"delivered"`
	Bounces         int    `json:"bounces"`
	Clicks          int    `json:"clicks"`
	UniqueClicks    int    `json:"unique_clicks"`
	Opens           int    `json:"opens"`
	UniqueOpens     int    `json:"unique_opens"`
	SpamReports     int    `json:"spamReports"`
	Blocked         int    `json:"blocked"`
	Unsubscribed    int    `json:"unsubscribed"`
	Invalid         int    `json:"invalid"`
	Date            string `json:"date"`
	EmailCampaignID int    `json:"email_campaign_id"`
}

// GetEmailCampaignAnalytics this cron job runs every day at 12:01am and gets the analytics for all email campaigns
func GetEmailCampaignAnalytics() {
	scheduler := tasks.New()
	now := time.Now()
	tomorrow := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 1, 0, 0, now.Location())
	duration := time.Until(tomorrow)

	_, err := scheduler.Add(&tasks.Task{
		Interval: duration,
		RunOnce:  false,
		TaskFunc: func() error {
			// get all email campaigns
			var emailCampaigns []model.EmailCampaign
			config.DB.Find(&emailCampaigns)

			for _, emailCampaign := range emailCampaigns {
				now := time.Now().Format("2006-01-02")

				//check if there is already an analytics entry for today
				var emailCampaignAnalytics []model.EmailCampaignAnalytics
				config.DB.Where("\"emailCampaignId\" = ? AND date = ?", emailCampaign.ID, now).Find(&emailCampaignAnalytics)
				if len(emailCampaignAnalytics) > 0 {
					continue
				}

				var EmailCampaignAnalytics model.EmailCampaignAnalytics
				EmailCampaignAnalytics.Date = now
				EmailCampaignAnalytics.EmailCampaignID = int64(int(emailCampaign.ID))

				err := config.DB.Create(&EmailCampaignAnalytics).Error

				if err != nil {
					return err
				}
				return nil
			}
			return nil
		},
	})

	if err != nil {
		fmt.Println(err)
	}
}

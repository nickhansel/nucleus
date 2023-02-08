package email

import (
	"fmt"
	"github.com/madflojo/tasks"
	"github.com/nickhansel/nucleus/config"
	"github.com/nickhansel/nucleus/model"
	"github.com/nickhansel/nucleus/sendinblue"
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
				analytics, err := sendinblue.GetEmailAnalytics(emailCampaign.ID)
				var EmailCampaignAnalytics model.EmailCampaignAnalytics
				EmailCampaignAnalytics.Sent = int32(analytics.Reports[0].Requests)
				EmailCampaignAnalytics.Delivered = int32(analytics.Reports[0].Delivered)
				EmailCampaignAnalytics.Bounces = int32(analytics.Reports[0].HardBounces + analytics.Reports[0].SoftBounces)
				EmailCampaignAnalytics.Clicks = int32(analytics.Reports[0].Clicks)
				EmailCampaignAnalytics.UniqueClicks = int32(analytics.Reports[0].UniqueClicks)
				EmailCampaignAnalytics.Opens = int32(analytics.Reports[0].Opens)
				EmailCampaignAnalytics.UniqueOpens = int32(analytics.Reports[0].UniqueOpens)
				EmailCampaignAnalytics.SpamReports = int32(analytics.Reports[0].SpamReports)
				EmailCampaignAnalytics.Blocked = int32(analytics.Reports[0].Blocked)
				EmailCampaignAnalytics.Unsubscribed = int32(analytics.Reports[0].Unsubscribed)
				EmailCampaignAnalytics.Invalid = int32(analytics.Reports[0].Invalid)
				EmailCampaignAnalytics.Date = analytics.Reports[0].Date
				EmailCampaignAnalytics.EmailCampaignID = int64(int(emailCampaign.ID))

				err = config.DB.Create(&EmailCampaignAnalytics).Error

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

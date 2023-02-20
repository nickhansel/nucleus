package fb

import (
	"fmt"
	"github.com/madflojo/tasks"
	fbMetrics "github.com/nickhansel/nucleus/api/analytics/fb"
	"github.com/nickhansel/nucleus/config"
	"github.com/nickhansel/nucleus/model"
	"strconv"
	"time"
)

//GetMetricsFromFB(org) returns the metrics from the facebook ad

func ScheduleGetFBMetrics() {
	scheduler := tasks.New()

	duration := time.Duration(10) * time.Second

	id, err := scheduler.Add(&tasks.Task{
		//interval is every 1.3 hours
		Interval: duration,
		RunOnce:  false,
		TaskFunc: func() error {
			var fbCampaigns []model.FbCampaign
			config.DB.Find(&fbCampaigns)

			for _, campaign := range fbCampaigns {
				fmt.Println("campaign", campaign.ID)

				var camp model.Campaign
				config.DB.Where("id = ?", campaign.CampaignID).First(&camp)

				var org model.Organization
				config.DB.Where("id = ?", camp.OrganizationID).First(&org)

				if org.FbAccessToken != "" && campaign.FbId != "" {
					err, metrics := fbMetrics.GetMetricsFromFB(org, "ad", campaign.FbId)
					if err != nil && metrics == nil {
						fmt.Println(err)
						return err
					}
					var FbAdMetrics model.FbAdMetrics
					//	convert each field in the struct to the correct type
					FbAdMetrics.StartDate = metrics["data"].([]interface{})[0].(map[string]interface{})["date_start"].(string)
					FbAdMetrics.EndDate = metrics["data"].([]interface{})[0].(map[string]interface{})["date_stop"].(string)
					FbAdMetrics.Clicks, _ = strconv.Atoi(metrics["data"].([]interface{})[0].(map[string]interface{})["clicks"].(string))
					FbAdMetrics.Reach, _ = strconv.Atoi(metrics["data"].([]interface{})[0].(map[string]interface{})["reach"].(string))
					FbAdMetrics.Impressions, _ = strconv.Atoi(metrics["data"].([]interface{})[0].(map[string]interface{})["impressions"].(string))
					FbAdMetrics.UniqueClicks, _ = strconv.Atoi(metrics["data"].([]interface{})[0].(map[string]interface{})["unique_clicks"].(string))
					FbAdMetrics.Cpc, _ = strconv.ParseFloat(metrics["data"].([]interface{})[0].(map[string]interface{})["cpc"].(string), 64)
					FbAdMetrics.Cpm, _ = strconv.ParseFloat(metrics["data"].([]interface{})[0].(map[string]interface{})["cpm"].(string), 64)
					FbAdMetrics.Ctr, _ = strconv.ParseFloat(metrics["data"].([]interface{})[0].(map[string]interface{})["ctr"].(string), 64)
					FbAdMetrics.Frequency, _ = strconv.Atoi(metrics["data"].([]interface{})[0].(map[string]interface{})["frequency"].(string))
					FbAdMetrics.OrganizationID = org.ID
					FbAdMetrics.Spend, _ = strconv.Atoi(metrics["data"].([]interface{})[0].(map[string]interface{})["spend"].(string))
					FbAdMetrics.AccountID = metrics["data"].([]interface{})[0].(map[string]interface{})["account_id"].(string)
					FbAdMetrics.CampaignID = metrics["data"].([]interface{})[0].(map[string]interface{})["campaign_id"].(string)
					FbAdMetrics.AdsetID = metrics["data"].([]interface{})[0].(map[string]interface{})["adset_id"].(string)
					FbAdMetrics.AdID = metrics["data"].([]interface{})[0].(map[string]interface{})["ad_id"].(string)
					FbAdMetrics.Level = "ad"
					config.DB.Create(&FbAdMetrics)

					err, metrics = fbMetrics.GetMetricsFromFB(org, "adset", campaign.FbId)
					if err != nil && metrics == nil {
						fmt.Println(err)
						return err
					}
					var FbAdsetMetrics model.FbAdMetrics
					//	convert each field in the struct to the correct type
					FbAdsetMetrics.StartDate = metrics["data"].([]interface{})[0].(map[string]interface{})["date_start"].(string)
					FbAdsetMetrics.EndDate = metrics["data"].([]interface{})[0].(map[string]interface{})["date_stop"].(string)
					FbAdsetMetrics.Clicks, _ = strconv.Atoi(metrics["data"].([]interface{})[0].(map[string]interface{})["clicks"].(string))
					FbAdsetMetrics.Reach, _ = strconv.Atoi(metrics["data"].([]interface{})[0].(map[string]interface{})["reach"].(string))
					FbAdsetMetrics.Impressions, _ = strconv.Atoi(metrics["data"].([]interface{})[0].(map[string]interface{})["impressions"].(string))
					FbAdsetMetrics.UniqueClicks, _ = strconv.Atoi(metrics["data"].([]interface{})[0].(map[string]interface{})["unique_clicks"].(string))
					//convert to float from string
					FbAdsetMetrics.Cpc, _ = strconv.ParseFloat(metrics["data"].([]interface{})[0].(map[string]interface{})["cpc"].(string), 64)
					FbAdsetMetrics.Cpm, _ = strconv.ParseFloat(metrics["data"].([]interface{})[0].(map[string]interface{})["cpm"].(string), 64)
					FbAdsetMetrics.Ctr, _ = strconv.ParseFloat(metrics["data"].([]interface{})[0].(map[string]interface{})["ctr"].(string), 64)
					FbAdsetMetrics.Frequency, _ = strconv.Atoi(metrics["data"].([]interface{})[0].(map[string]interface{})["frequency"].(string))
					FbAdsetMetrics.OrganizationID = org.ID
					FbAdsetMetrics.Spend, _ = strconv.Atoi(metrics["data"].([]interface{})[0].(map[string]interface{})["spend"].(string))
					FbAdsetMetrics.AdsetID = metrics["data"].([]interface{})[0].(map[string]interface{})["adset_id"].(string)
					FbAdsetMetrics.AccountID = metrics["data"].([]interface{})[0].(map[string]interface{})["account_id"].(string)
					FbAdsetMetrics.CampaignID = metrics["data"].([]interface{})[0].(map[string]interface{})["campaign_id"].(string)

					FbAdsetMetrics.Level = "adset"
					config.DB.Create(&FbAdsetMetrics)

					err, metrics = fbMetrics.GetMetricsFromFB(org, "campaign", campaign.FbId)
					if err != nil && metrics == nil {
						fmt.Println(err)
						return err
					}
					var FbCampaignMetrics model.FbAdMetrics
					//	convert each field in the struct to the correct type
					FbCampaignMetrics.StartDate = metrics["data"].([]interface{})[0].(map[string]interface{})["date_start"].(string)
					FbCampaignMetrics.EndDate = metrics["data"].([]interface{})[0].(map[string]interface{})["date_stop"].(string)
					FbCampaignMetrics.Clicks, _ = strconv.Atoi(metrics["data"].([]interface{})[0].(map[string]interface{})["clicks"].(string))
					FbCampaignMetrics.Reach, _ = strconv.Atoi(metrics["data"].([]interface{})[0].(map[string]interface{})["reach"].(string))
					FbCampaignMetrics.Impressions, _ = strconv.Atoi(metrics["data"].([]interface{})[0].(map[string]interface{})["impressions"].(string))
					FbCampaignMetrics.UniqueClicks, _ = strconv.Atoi(metrics["data"].([]interface{})[0].(map[string]interface{})["unique_clicks"].(string))
					FbCampaignMetrics.Cpc, _ = strconv.ParseFloat(metrics["data"].([]interface{})[0].(map[string]interface{})["cpc"].(string), 64)
					FbCampaignMetrics.Cpm, _ = strconv.ParseFloat(metrics["data"].([]interface{})[0].(map[string]interface{})["cpm"].(string), 64)
					FbCampaignMetrics.Ctr, _ = strconv.ParseFloat(metrics["data"].([]interface{})[0].(map[string]interface{})["ctr"].(string), 64)
					FbCampaignMetrics.Frequency, _ = strconv.Atoi(metrics["data"].([]interface{})[0].(map[string]interface{})["frequency"].(string))
					FbCampaignMetrics.OrganizationID = org.ID
					FbCampaignMetrics.Spend, _ = strconv.Atoi(metrics["data"].([]interface{})[0].(map[string]interface{})["spend"].(string))
					FbCampaignMetrics.AccountID = metrics["data"].([]interface{})[0].(map[string]interface{})["account_id"].(string)
					FbCampaignMetrics.CampaignID = metrics["data"].([]interface{})[0].(map[string]interface{})["campaign_id"].(string)
					FbCampaignMetrics.Level = "campaign"
					config.DB.Create(&FbCampaignMetrics)
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

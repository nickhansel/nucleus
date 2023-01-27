package model

const TableNameEmailCampaignAnalytics = "email_campaign_analytics"

type EmailCampaignAnalytics struct {
	ID              int32  `gorm:"column:id;primaryKey" json:"id"`
	Date            string `gorm:"column:date;not null" json:"date"`
	EmailCampaignID int32  `gorm:"column:emailCampaignId;not null" json:"emailCampaignId"`
	Sent            int32  `gorm:"column:sent;not null" json:"sent"`
	Delivered       int32  `gorm:"column:delivered;not null" json:"delivered"`
	Bounces         int32  `gorm:"column:bounces;not null" json:"bounces"`
	Clicks          int32  `gorm:"column:clicks;not null" json:"clicks"`
	UniqueClicks    int32  `gorm:"column:unique_clicks;not null" json:"unique_clicks"`
	Opens           int32  `gorm:"column:opens;not null" json:"opens"`
	UniqueOpens     int32  `gorm:"column:unique_opens;not null" json:"unique_opens"`
	SpamReports     int32  `gorm:"column:spamReports;not null" json:"spamReports"`
	Blocked         int32  `gorm:"column:blocked;not null" json:"blocked"`
	Unsubscribed    int32  `gorm:"column:unsubscribed;not null" json:"unsubscribed"`
	Invalid         int32  `gorm:"column:invalid;not null" json:"invalid"`
}

// TableName Organization's table name
func (*EmailCampaignAnalytics) TableName() string {
	return TableNameEmailCampaignAnalytics
}

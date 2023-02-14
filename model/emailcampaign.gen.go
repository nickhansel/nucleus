// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import "github.com/lib/pq"

const TableNameEmailCampaign = "EmailCampaign"

// EmailCampaign mapped from table <EmailCampaign>
type EmailCampaign struct {
	ID           int64  `gorm:"column:id;primaryKey" json:"id"`
	TargetEmails  pq.StringArray `gorm:"type:varchar(255)[]"  json:"target_emails"`
	From         string `gorm:"column:from;not null" json:"from"`
	SendTime     string `gorm:"column:send_time;not null" json:"send_time"`
	Subject      string `gorm:"column:subject;not null" json:"subject"`
	Text         string `gorm:"column:text;not null" json:"text"`
	HTML         string `gorm:"column:html;not null" json:"html"`
	CampaignID   int64  `gorm:"column:campaignId;not null" json:"campaignId"`
	EmailCampaignAnalytics []EmailCampaignAnalytics `gorm:"column:email_campaign_analytics;foreignKey:emailCampaignId;references:ID" json:"email_campaign_analytics"`
}

// TableName EmailCampaign's table name
func (*EmailCampaign) TableName() string {
	return TableNameEmailCampaign
}

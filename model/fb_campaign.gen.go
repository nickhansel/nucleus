// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

const TableNameFbCampaign = "fb_campaign"

// FbCampaign mapped from table <fb_campaign>
type FbCampaign struct {
	ID         int32  `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Name       string `gorm:"column:name;not null" json:"name"`
	Objective  string `gorm:"column:objective;not null" json:"objective"`
	Status     string `gorm:"column:status;not null" json:"status"`
	CampaignID int32 `gorm:"column:campaignId;not null" json:"campaignId"`
}

// TableName FbCampaign's table name
func (*FbCampaign) TableName() string {
	return TableNameFbCampaign
}

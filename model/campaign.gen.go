// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

const TableNameCampaign = "Campaign"

// Campaign mapped from table <Campaign>
type Campaign struct {
	ID             int32   `gorm:"column:id;primaryKey" json:"id"`
	CampaignId    string  `gorm:"column:campaign_id;" json:"campaign_id"`
	CreatedAt      string  `gorm:"column:created_at;not null" json:"created_at"`
	Type           string  `gorm:"column:type;not null" json:"type"`
	Budget         float64 `gorm:"column:budget;not null" json:"budget"`
	OrganizationID int32   `gorm:"column:organizationId;not null" json:"organizationId"`
	FbCampaign    FbCampaign `gorm:"foreignKey:campaignId;references:ID"`
}

// TableName Campaign's table name
func (*Campaign) TableName() string {
	return TableNameCampaign
}

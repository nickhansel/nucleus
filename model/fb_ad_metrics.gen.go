package model

//model fb_ad_metrics {
//id         BigInt   @id @default(autoincrement())
//start_date String
//end_date   String
//impressions Int? @default(0)
//reach      Int? @default(0)
//clicks     Int? @default(0)
//cpm        Float? @default(0.0)
//cpc        Float? @default(0.0)
//ctr        Float? @default(0.0)
//unique_clicks Int? @default(0)
//frequency  Float? @default(0.0)
//spend      Float
//account_id String
//ad_id      String
//campaign_id String
//adset_id   String
//level      String
//organizationId BigInt
//organization Organization @relation(fields: [organizationId], references: [id])
//
//@@index([id])
//}

const TableNameFb_Ad_Metrics = "fb_ad_metrics"

type FbAdMetrics struct {
	ID             int64   `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	StartDate      string  `gorm:"column:start_date;not null" json:"start_date"`
	EndDate        string  `gorm:"column:end_date;not null" json:"end_date"`
	Impressions    int     `gorm:"column:impressions;not null" json:"impressions"`
	Reach          int     `gorm:"column:reach;not null" json:"reach"`
	Clicks         int     `gorm:"column:clicks;not null" json:"clicks"`
	Cpm            float64 `gorm:"column:cpm;not null" json:"cpm"`
	Cpc            float64 `gorm:"column:cpc;not null" json:"cpc"`
	Ctr            float64 `gorm:"column:ctr;not null" json:"ctr"`
	UniqueClicks   int     `gorm:"column:unique_clicks;not null" json:"unique_clicks"`
	Frequency      int     `gorm:"column:frequency;not null" json:"frequency"`
	Spend          int     `gorm:"column:spend;not null" json:"spend"`
	AccountID      string  `gorm:"column:account_id;not null" json:"account_id"`
	AdID           string  `gorm:"column:ad_id;not null" json:"ad_id"`
	CampaignID     string  `gorm:"column:campaign_id;not null" json:"campaign_id"`
	AdsetID        string  `gorm:"column:adset_id;not null" json:"adset_id"`
	Level          string  `gorm:"column:level;not null" json:"level"`
	OrganizationID int64   `gorm:"column:organizationId" json:"organizationId"`
}

// TableName FbAd's table name
func (*FbAdMetrics) TableName() string {
	return TableNameFb_Ad_Metrics
}

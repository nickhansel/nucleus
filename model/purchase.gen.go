// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

const TableNamePurchase = "Purchase"

// Purchase mapped from table <Purchase>
type Purchase struct {
	ID          int64   `gorm:"column:id;primaryKey" json:"id"`
	CreatedAt   string  `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt   string  `gorm:"column:updated_at;not null" json:"updated_at"`
	AmountMoney float64 `gorm:"column:amount_money;not null" json:"amount_money"`
	Currency    string  `gorm:"column:currency;not null" json:"currency"`
	Status      string  `gorm:"column:status;not null" json:"status"`
	SourceType  string  `gorm:"column:source_type;not null" json:"source_type"`
	PurchasedItems []PurchasedItem `gorm:"foreignKey:PurchaseID;references:PurchaseID"`
	CustomerID  int64   `gorm:"column:customerId" json:"customerId"`
	PurchaseID  string  `gorm:"column:purchase_id;not null" json:"purchase_id"`
	ProductType string  `gorm:"column:product_type" json:"product_type"`
	LocationID  int64   `gorm:"column:locationId;not null" json:"locationId"`
	Location_ID string `gorm:"column:location_id;not null" json:"location_id"`
	ItemsID     int64   `gorm:"column:itemsId" json:"itemsId"`
	Location   StoreLocation `gorm:"foreignKey:LocationID;references:ID"`
	AttributedCampaign Campaign `gorm:"foreignKey:AttributedCampaignID;references:ID"`
	AttributedCampaignID int64 `gorm:"column:attributedCampaignId;not null" json:"attributedCampaignId"`
	Customer  Customer `gorm:"foreignKey:CustomerID;references:ID"`
	OrganizationID int64 `gorm:"column:organizationId;not null" json:"organizationId"`
}

// TableName Purchase's table name
func (*Purchase) TableName() string {
	return TableNamePurchase
}

// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

const TableNamePurchasedItem = "purchased_item"

// PurchasedItem mapped from table <purchased_item>
type PurchasedItem struct {
	ID          int64   `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	ItemID      int64   `gorm:"column:itemId" json:"itemId"`
	Name        string  `gorm:"column:name;not null" json:"name"`
	Quantity    int32   `gorm:"column:quantity;not null" json:"quantity"`
	Cost        float64 `gorm:"column:cost;not null" json:"cost"`
	PurchaseID  string  `gorm:"column:purchaseId;not null" json:"purchaseId"`
	IsVaration  bool    `gorm:"column:is_varation;not null" json:"is_varation"`
	VariationID int64   `gorm:"column:variationId" json:"variationId"`
	// variations is a relation to table <variations>
	Variation Variation `gorm:"foreignKey:VariationID;references:ID" json:"variation"`
}

// TableName PurchasedItem's table name
func (*PurchasedItem) TableName() string {
	return TableNamePurchasedItem
}

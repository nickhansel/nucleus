// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

const TableNameItemsToitemCategory = "_CustomerTocustomer_groups"

// ItemsToitemCategory mapped from table <_ItemsToitem_category>
type CustomersToCustomerGroups struct {
	// A is a foreign key to Customer table that has primary key ID
	A int32 `gorm:"column:A" json:"A"`

	// B is a foreign key to CustomerGroup table that has primary key ID
	B int32 `gorm:"column:B" json:"B"`

	// many to one relationship of Customer to CustomerGroup
	Customer Customer `gorm:"foreignKey:A;references:ID"`

	// many to one relationship of CustomerGroup to Customer
	CustomerGroup CustomerGroup `gorm:"foreignKey:B;references:ID"`
}

// TableName ItemsToitemCategory's table name
func (*CustomersToCustomerGroups) TableName() string {
	return TableNameItemsToitemCategory
}
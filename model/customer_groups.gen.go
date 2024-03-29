// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model
import "time"

const TableNameCustomerGroup = "customer_groups"

// CustomerGroup mapped from table <customer_groups>
type CustomerGroup struct {
	ID             int64  `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Name           string `gorm:"column:name;not null" json:"name"`
	// many to one relationship of Customers to CustomerGroup
	OrganizationID int64  `gorm:"column:organizationId;not null" json:"organizationId"`
	FbCustomAudienceID string `gorm:"column:fb_custom_audience_id" json:"fb_custom_audience_id"`
	CreatedAt	  time.Time `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt	  time.Time `gorm:"column:updated_at;not null" json:"updated_at"`
	// Customers is a reference to _CustomersTo_customer_groups table that has foreign key B
	Customers []Customer `gorm:"foreignKey:ID;references:ID" json:"customers"`
}

// TableName CustomerGroup's table name
func (*CustomerGroup) TableName() string {
	return TableNameCustomerGroup
}

// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

const TableNameOrganizationToUser = "_OrganizationToUser"

// OrganizationToUser mapped from table <_OrganizationToUser>
type OrganizationToUser struct {
	A int32 `gorm:"column:A;not null" json:"A"`
	B int32 `gorm:"column:B;not null" json:"B"`
}

// TableName OrganizationToUser's table name
func (*OrganizationToUser) TableName() string {
	return TableNameOrganizationToUser
}

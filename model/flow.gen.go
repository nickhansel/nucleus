package model

const TabelFlow = "flow"

type Flow struct {
	ID              int32     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Name            string    `gorm:"column:name;not null" json:"name"`
	Trigger         []Trigger `gorm:"foreignKey:FlowID;references:ID" json:"trigger"`
	OrganizationID  int32     `gorm:"column:organizationId;not null" json:"organizationId"`
	CustomerGroupID int32     `gorm:"column:customer_groupId;not null" json:"customer_groupId"`
}

// TableName FbTarget's table name
func (*Flow) TableName() string {
	return TabelFlow
}

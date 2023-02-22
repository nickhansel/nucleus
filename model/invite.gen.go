package model

const TableNameInvite = "invite"

// Invite mapped from table <geo_locations>
type Invite struct {
	ID             int64  `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Email          string `gorm:"column:email;not null" json:"email"`
	IsAccepted     bool   `gorm:"column:isAccepted;not null" json:"isAccepted"`
	OrganizationID int64  `gorm:"column:organizationId;not null" json:"organizationId"`
}

// TableName GeoLocation's table name
func (*Invite) TableName() string {
	return TableNameInvite
}

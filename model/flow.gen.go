package model

const TabelFlow = "flow"

type Flow struct {
	ID                        int64     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Name                      string    `gorm:"column:name;not null" json:"name"`
	Trigger                   []Trigger `gorm:"foreignKey:FlowID;references:ID" json:"trigger"`
	OrganizationID            int64     `gorm:"column:organizationId;not null" json:"organizationId"`
	CustomerGroupID           int64     `gorm:"column:customer_groupId;not null" json:"customer_groupId"`
	Status                    string    `gorm:"column:status;not null" json:"status"`
	CreatedAt                 string    `gorm:"column:created_at;not null" json:"created_at"`
	SmartSending              bool      `gorm:"column:smart_sending;not null" json:"smart_sending"`
	MessageType               string    `gorm:"column:message_type;not null" json:"message_type"`
	AttributedRevenue         float64   `gorm:"column:attributed_revenue;not null" json:"attributed_revenue"`
	AttributedConversionCount int64     `gorm:"column:attributed_conversion_count;not null" json:"attributed_conversion_count"`
}

// TableName FbTarget's table name
func (*Flow) TableName() string {
	return TabelFlow
}

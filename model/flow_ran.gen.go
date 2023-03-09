package model

const TableFlow_Ran = "flow_ran"

type FlowRan struct {
	ID          int64  `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	FlowID      int64  `gorm:"column:flowId;not null" json:"flowId"`
	CustomerID  int64  `gorm:"column:customerId;not null" json:"customerId"`
	Date        string `gorm:"column:date;not null" json:"date"`
	Trigger     string `gorm:"column:trigger;not null" json:"trigger"`
	MessageType string `gorm:"column:message_type;not null" json:"message_type"`
}

// TableName FbTarget's table name
func (*FlowRan) TableName() string {
	return TableFlow_Ran
}

package model

const TabelTrigger = "trigger"

type Trigger struct {
	ID     int64    `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Event  string   `gorm:"column:event;not null" json:"event"`
	FlowID int64    `gorm:"column:flowId;not null" json:"flowId"`
	Action []Action `gorm:"foreignKey:TriggerID;references:ID" json:"action"`
}

// TableName FbTarget's table name
func (*Trigger) TableName() string {
	return TabelTrigger
}

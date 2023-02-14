package model

const TabelAction = "action"

type Action struct {
	ID        int64  `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Type      string `gorm:"column:type;not null" json:"type"`
	TriggerID int64  `gorm:"column:triggerId;not null" json:"triggerId"`
	WaitTime  string `gorm:"column:wait_time;not null" json:"waitTime"`
}

// TableName FbTarget's table name
func (*Action) TableName() string {
	return TabelAction
}

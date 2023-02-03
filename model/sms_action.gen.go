package model

const TabelSms_action = "sms_action"

type SmsAction struct {
	ID       int32  `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Body     string `gorm:"column:body;not null" json:"body"`
	ActionID int32  `gorm:"column:actionId;not null" json:"actionId"`
}

// TableName FbTarget's table name
func (*SmsAction) TableName() string {
	return TabelSms_action
}

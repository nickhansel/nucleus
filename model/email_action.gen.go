package model

const TabelEmail_action = "email_action"

type EmailAction struct {
	ID       int64  `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Subject  string `gorm:"column:subject;not null" json:"subject"`
	Body     string `gorm:"column:body;not null" json:"body"`
	ActionID int64  `gorm:"column:actionId;not null" json:"actionId"`
}

// TableName FbTarget's table name
func (*EmailAction) TableName() string {
	return TabelEmail_action
}

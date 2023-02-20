package model

const Tabelemail_template = "email_template"

type EmailTemplate struct {
	ID             int64  `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Html           string `gorm:"column:html" json:"html"`
	Name           string `gorm:"column:name" json:"name"`
	OrganizationID int64  `gorm:"column:organizationId" json:"organizationId"`
}

// TableName FbTarget's table name
func (*EmailTemplate) TableName() string {
	return Tabelemail_template
}

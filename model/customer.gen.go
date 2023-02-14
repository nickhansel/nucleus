// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model
import "github.com/lib/pq"

const TableNameCustomer = "Customer"

// Customer mapped from table <Customer>
type Customer struct {
	ID                           int64   `gorm:"column:id;primaryKey" json:"id"`
	CreatedAt                    string  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt                    string  `gorm:"column:updated_at" json:"updated_at"`
	GivenName                    string  `gorm:"column:given_name;not null" json:"given_name"`
	FamilyName                   string  `gorm:"column:family_name;not null" json:"family_name"`
	EmailAddress                 string  `gorm:"column:email_address" json:"email_address"`
	PhoneNumber                  string  `gorm:"column:phone_number" json:"phone_number"`
	ReferenceID                  string  `gorm:"column:reference_id" json:"reference_id"`
	Note                         string  `gorm:"column:note" json:"note"`
	EmailUnsubscribed            bool    `gorm:"column:email_unsubscribed;not null" json:"email_unsubscribed"`
	CreationSource               string  `gorm:"column:creation_source;not null" json:"creation_source"`
	OrganizationID               int64   `gorm:"column:organizationId;not null" json:"organizationId"`
	PosID                        string  `gorm:"column:pos_id;not null" json:"pos_id"`
	Birthday                     string  `gorm:"column:birthday" json:"birthday"`
	PosName                      string  `gorm:"column:pos_name;not null" json:"pos_name"`
	TotalPurhcases               int32   `gorm:"column:total_purhcases" json:"total_purhcases"`
	TotalSpent                   float64 `gorm:"column:total_spent" json:"total_spent"`
	AddressLine1                 string  `gorm:"column:address_line_1" json:"address_line_1"`
	AddressLine2                 string  `gorm:"column:address_line_2" json:"address_line_2"`
	AdministrativeDistrictLevel1 string  `gorm:"column:administrative_district_level_1" json:"administrative_district_level_1"`
	Country                      string  `gorm:"column:country" json:"country"`
	Locality                     string  `gorm:"column:locality" json:"locality"`
	PostalCode                   string  `gorm:"column:postal_code" json:"postal_code"`
	CustomerGroup []CustomerGroup `gorm:"foreignKey:ID;references:ID" json:"customers"`
	IsEmailDeliverable bool `gorm:"column:is_email_deliverable" json:"is_email_deliverable"`
	IsSMSDeliverable bool `gorm:"column:is_sms_deliverable" json:"is_sms_deliverable"`
	DatesReceivedEmail  pq.StringArray `gorm:"type:varchar(255)[]" json:"dates_received_email"`
	DatesReceivedSMS  pq.StringArray `gorm:"type:varchar(255)[]" json:"dates_received_sms"`
	SmsUnsubscribed bool `gorm:"column:sms_unsubscribed" json:"sms_unsubscribed"`
}

// TableName Customer's table name
func (*Customer) TableName() string {
	return TableNameCustomer
}

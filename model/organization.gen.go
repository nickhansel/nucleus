// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
	// import pq
)

const TableNameOrganization = "Organization"

type MultiString []string

// Organization mapped from table <Organization>
type Organization struct {
	ID                int64     `gorm:"column:id;primaryKey" json:"id"`
	FbAdAccountID	 string     `gorm:"column:fb_ad_account_id;" json:"fb_ad_account_id"`
	FbAccessToken	 string    `gorm:"column:fb_access_token;" json:"fb_access_token"`
	FbPageID		 string    `gorm:"column:fb_page_id;" json:"fb_page_id"`
	FbPageFanCount	 int32    `gorm:"column:fb_page_fan_count;" json:"fb_page_fan_count"`
	FbPageName 		 string    `gorm:"column:fb_page_name;" json:"fb_page_name"`
	FbPageImgURL 	 string    `gorm:"column:fb_page_img_url;" json:"fb_page_img_url"`
	SendgridID 	  int32    `gorm:"column:sendgrid_id;" json:"sendgrid_id"`
	Members 		 []User     `gorm:"foreignKey:OrganizationID;references:ID"` 
	Name              string    `gorm:"column:name;not null" json:"name"`
	CreatedAt         time.Time `gorm:"column:createdAt;not null;default:CURRENT_TIMESTAMP" json:"createdAt"`
	FbCustomAudiences MultiString `gorm:"type:text[]" json:"fbCustomAudiences"`
	IsSendgridAuthed  bool      `gorm:"column:is_sendgrid_authed;not null" json:"is_sendgrid_authed"`
	SendinblueEmail   string    `gorm:"column:sendinblue_email;not null" json:"sendinblue_email"`
	IsSendinblueAuthed  bool      `gorm:"column:is_sendinblue_authed;not null" json:"is_sendinblue_authed"`
	IsTwilioAuthed    bool      `gorm:"column:is_twilio_authed;not null" json:"is_twilio_authed"`
	Plan              int32     `gorm:"column:plan;not null" json:"plan"`
	PosID             int64     `gorm:"column:posId;not null" json:"posId"`
	SendgridEmail     string    `gorm:"column:sendgrid_email;not null" json:"sendgrid_email"`
	TwilioNumber      string    `gorm:"column:twilio_number;not null" json:"twilio_number"`
	Type              string    `gorm:"column:type;not null" json:"type"`
	UpdatedAt         time.Time `gorm:"column:updatedAt;not null;default:CURRENT_TIMESTAMP" json:"updatedAt"`
	EmailVerificationCode string `gorm:"column:email_verification_code;not null" json:"email_verification_code"`
	ShopifyUrl		string    `gorm:"column:shopify_url;not null" json:"shopify_url"`
	EmailCount int32 `gorm:"column:email_count;not null" json:"email_count"`
	EmailLimit int32 `gorm:"column:email_limit;not null" json:"email_limit"`
	SmsCount int32 `gorm:"column:sms_count;not null" json:"sms_count"`
	SmsLimit int32 `gorm:"column:sms_limit;not null" json:"sms_limit"`
}

// TableName Organization's table name
func (*Organization) TableName() string {
	return TableNameOrganization
}

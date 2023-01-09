// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

const TableNameFbTarget = "fb_target"

// FbTarget mapped from table <fb_target>
type FbTarget struct {
	ID              int32  `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	AgeMin          int32  `gorm:"column:age_min;not null" json:"age_min"`
	AgeMax          int32  `gorm:"column:age_max;not null" json:"age_max"`
	CustomAudiences string `gorm:"column:custom_audiences;not null" json:"custom_audiences"`
	FbAdsetID       int32  `gorm:"column:fb_adsetId;not null" json:"fb_adsetId"`
	GeoLocationsID  int32  `gorm:"column:geo_locationsId;not null" json:"geo_locationsId"`
}

// TableName FbTarget's table name
func (*FbTarget) TableName() string {
	return TableNameFbTarget
}

// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import "github.com/lib/pq"

const TableNameGeoLocation = "geo_locations"

// GeoLocation mapped from table <geo_locations>
type GeoLocation struct {
	ID         int64  `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Countries  pq.StringArray `gorm:"type:varchar(255)[]" json:"countries"`
	Regions    pq.StringArray `gorm:"type:varchar(255)[]" json:"regions"`
	Cities     pq.StringArray `gorm:"type:varchar(255)[]" json:"cities"`
	ZipCodes   pq.StringArray `gorm:"type:varchar(255)[]" json:"zip_codes"`
	FbTargetID int64  `gorm:"column:fb_targetId;not null" json:"fb_targetId"`
}

// TableName GeoLocation's table name
func (*GeoLocation) TableName() string {
	return TableNameGeoLocation
}

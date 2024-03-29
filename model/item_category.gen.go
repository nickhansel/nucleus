// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

const TableNameItemCategory = "item_category"

// ItemCategory mapped from table <item_category>
type ItemCategory struct {
	ID        int64  `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	CreatedAt string `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt string `gorm:"column:updated_at;not null" json:"updated_at"`
	IsDeleted bool   `gorm:"column:isDeleted;not null" json:"isDeleted"`
	Name      string `gorm:"column:name;not null" json:"name"`
}

// TableName ItemCategory's table name
func (*ItemCategory) TableName() string {
	return TableNameItemCategory
}

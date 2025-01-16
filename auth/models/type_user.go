package models

type TypeUser struct {
	ID   int    `gorm:"type:int(11);primaryKey;autoIncrement" json:"id"`
	Type string `gorm:"type:varchar(30);not null" json:"type"`
}

func (TypeUser) TableName() string {
	return "type_user"
}
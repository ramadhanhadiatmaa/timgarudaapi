package models

import "time"

type User struct {
	Email     string    `gorm:"type:varchar(100);primaryKey;" json:"email"`
	Password  string    `gorm:"type:varchar(255);" json:"password"`
	FullName  string    `gorm:"column:first_name;type:varchar(255);null" json:"full_name"`
	Phone     string    `gorm:"type:varchar(20);null" json:"phone"`
	Type      int       `gorm:"type:int(11)" json:"type"`
	TypeInfo  TypeUser  `gorm:"foreignKey:Type;references:ID" json:"type_info"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (User) TableName() string {
	return "user"
}

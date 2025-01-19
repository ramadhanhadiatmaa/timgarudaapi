package models

import "time"

type User struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Email     string    `gorm:"type:varchar(100);unique;not null" json:"email"`
	Password  string    `gorm:"type:varchar(255);not null" json:"password"`
	FullName  string    `gorm:"type:varchar(255)" json:"full_name"`
	Phone     string    `gorm:"type:varchar(20)" json:"phone"`
	Type      int       `gorm:"type:int(11)" json:"type"`
	TypeName  string    `gorm:"column:type_name" json:"type_name"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type TypeUser struct {
	ID   int    `gorm:"type:int(11);primaryKey;autoIncrement" json:"id"`
	Type string `gorm:"type:varchar(30);not null" json:"type"`
}

func (User) TableName() string {
	return "user"
}

func (TypeUser) TableName() string {
	return "type_user"
}

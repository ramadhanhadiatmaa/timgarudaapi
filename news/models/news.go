package models

import "time"

type News struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Content   string    `gorm:"type:varchar(1000);not null" json:"content"`
	Image     string    `gorm:"type:varchar(250);not null" json:"image"`
	Category  int       `gorm:"type:int(11)" json:"category"`
	CatInfo   string    `gorm:"column:category_name" json:"category_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (News) TableName() string {
	return "news"
}

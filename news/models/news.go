package models

import "time"

type News struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Content   string    `gorm:"type:varchar(250);not null" json:"content"`
	Image     string    `gorm:"type:varchar(250);not null" json:"image"`
	Category  int       `gorm:"type:int(11)" json:"category"`
	CatInfo   Category  `gorm:"foreignKey:Category;references:ID" json:"cat_info"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (News) TableName() string {
	return "news"
}

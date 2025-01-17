package models

import "time"

type NewsComment struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`
	PostID    int       `gorm:"not null" json:"post_id"`
	PostInfo  News      `gorm:"foreignKey:ID;references:ID" json:"post_info"`
	UserID    int       `gorm:"not null" json:"user_id"`
	UserInfo  User      `gorm:"foreignKey:ID;references:ID" json:"email_info"`
	Comment   string    `gorm:"type:varchar(250);not null" json:"comment"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (NewsComment) TableName() string {
	return "news_comment"
}

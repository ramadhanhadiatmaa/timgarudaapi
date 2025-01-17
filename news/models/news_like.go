package models

import "time"

type NewsLike struct {
	ID        int       `gorm:"type:int(11);primaryKey;autoIncrement" json:"id"`
	PostID    int       `gorm:"not null" json:"post_id"`
	PostInfo  News      `gorm:"foreignKey:ID;references:ID" json:"post_info"`
	UserID    int       `gorm:"not null" json:"user_id"`
	UserInfo  User      `gorm:"foreignKey:ID;references:ID" json:"email_info"`
	CreatedAt time.Time `json:"created_at"`
}

func (NewsLike) TableName() string {
	return "news_like"
}

package models

import "time"

type NewsLike struct {
	ID        int       `gorm:"type:int(11);primaryKey;autoIncrement" json:"id"`
	PostID    int       `gorm:"type:int(11);not null" json:"post_id"` // Kunci asing ke News.ID
	Post      News      `gorm:"foreignKey:PostID;references:ID" json:"post"` // Relasi ke News
	UserID    int       `gorm:"type:int(11);not null" json:"user_id"` // Kunci asing ke User.ID
	User      User      `gorm:"foreignKey:UserID;references:ID" json:"user"` // Relasi ke User
	CreatedAt time.Time `json:"created_at"`
}

func (NewsLike) TableName() string {
	return "news_like"
}

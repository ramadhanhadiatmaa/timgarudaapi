package models

import "time"

type NewsComment struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`
	PostID    int       `gorm:"type:int(11);not null" json:"post_id"` // Kunci asing ke News.ID
	Post      News      `gorm:"foreignKey:PostID;references:ID" json:"post"` // Relasi ke News
	UserID    int       `gorm:"type:int(11);not null" json:"user_id"` // Kunci asing ke User.ID
	User      User      `gorm:"foreignKey:UserID;references:ID" json:"user"` // Relasi ke User
	Comment   string    `gorm:"type:varchar(250);not null" json:"comment"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (NewsComment) TableName() string {
	return "news_comment"
}

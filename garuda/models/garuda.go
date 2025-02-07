package models

import "time"

type Garuda struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    int       `gorm:"not null" json:"user_id"`
	Value     int       `gorm:"type:int(11);not null" json:"value"`
	CreatedAt time.Time `json:"created_at"`
}

func (Garuda) TableName() string {
	return "garuda"
}

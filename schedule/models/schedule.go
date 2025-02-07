package models

import "time"

type Schedule struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`
	HomeID    int       `gorm:"not null" json:"home_id"`
	AwayID    int       `gorm:"not null" json:"away_id"`
	Home      Team      `gorm:"foreignKey:HomeID;references:ID" json:"home"`
	Away      Team      `gorm:"foreignKey:AwayID;references:ID" json:"away"`
	Venue     string    `gorm:"type:varchar(250);not null" json:"venue"`
	Date      string    `gorm:"type:varchar(250);not null" json:"date"`
	CreatedAt time.Time `json:"created_at"`
}

func (Schedule) TableName() string {
	return "schedule"
}

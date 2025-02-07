package models

type Team struct {
	ID   int    `gorm:"primaryKey;autoIncrement" json:"id"`
	Team string `gorm:"type:varchar(250);not null" json:"team"`
	Flag string `gorm:"type:varchar(250);not null" json:"flag"`
}

func (Team) TableName() string {
	return "team"
}

package models

type User struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`
}

func (User) TableName() string {
	return "user"
}

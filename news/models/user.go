package models

type User struct {
	Email     string    `gorm:"type:varchar(100);primaryKey;" json:"email"`
}

func (User) TableName() string {
	return "user"
}
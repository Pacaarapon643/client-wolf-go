package models

type User struct {
	BaseModel
	Email    string `gorm:"uniqueIndex;not null" json:"email"`
	Password string `gorm:"not null" json:"password"`
	UserName string `gorm:"not null" json:"user_name"`
	IsOnline bool   `gorm:"default:false" json:"is_online"`
}

func (User) TableName() string {
	return "users"
}

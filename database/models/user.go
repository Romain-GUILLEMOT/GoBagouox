package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	DiscordID string `gorm:"type:varchar(20);"`
	Username  string `gorm:"type:string;"`
	Avatar    string `gorm:"type:string;"`
	Email     string `gorm:"type:string;"`
}

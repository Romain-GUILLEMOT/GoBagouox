package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	DiscordID string `gorm:"type:varchar(20);"`
	Email     string `gorm:"type:varchar(100);"`
}

package models

import "gorm.io/gorm"

type SshAccess struct {
	gorm.Model
	IP       string `gorm:"type:varchar(15);"`
	Port     int    `gorm:"type:int;"`
	Username string `gorm:"type:varchar(50);"`
	OwnerID  uint
	Owner    User `gorm:"foreignKey:OwnerID"`
}

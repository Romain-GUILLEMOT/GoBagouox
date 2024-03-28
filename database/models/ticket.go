package models

import "gorm.io/gorm"

type Ticket struct {
	gorm.Model
	Nom       string `gorm:"type:varchar(100);not null;"`
	Status    string `gorm:"type:varchar(20);not null;"`
	License   string `gorm:"type:varchar(100);not null;default:''"`
	Logs      string `gorm:"type:varchar(50);not null;default:''"`
	ChannelId string `gorm:"type:varchar(100);not null;"`
	OwnerID   uint   `gorm:"not null"`
	Owner     User   `gorm:"foreignKey:OwnerID"`
}

type TicketMessage struct {
	gorm.Model
	Content   string `gorm:"type:text;not null"`
	MessageID string `gorm:"type:string;not null;"`
	TicketID  uint   `gorm:"not null"`
	Ticket    Ticket `gorm:"foreignKey:TicketID"`
	OwnerID   uint   `gorm:"not null"`
	Owner     User   `gorm:"foreignKey:OwnerID"`
}

type TicketAttachments struct {
	gorm.Model
	TicketMessageID uint          `gorm:"not null"`
	TicketMessage   TicketMessage `gorm:"foreignKey:TicketMessageID"`
	Uuid            string        `gorm:"type:char(36);primaryKey;"`
	Type            string        `gorm:"type:string;not null;"`
	Size            int64         `gorm:"type:bigint;not null;"`
	Name            string        `gorm:"type:text;not null;"`
}

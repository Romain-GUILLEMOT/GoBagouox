package web

import (
	"GoBagouox/database"
	"GoBagouox/database/models"
	"GoBagouox/utils"
	"errors"
)

func GetTicket(email string, ticketID string) (models.Ticket, error) {
	db := database.GetDB()
	var ticket models.Ticket
	err := db.Preload("TicketMessages").Preload("Owner").Preload("TicketMessages.Owner").Preload("TicketMessages.TicketAttachments").First(&ticket, ticketID).Error
	if err != nil {
		utils.Error("Failed to get ticket transcript.", err, 0)
		return models.Ticket{}, err
	}

	if ticket.Owner.Email != email {
		var user models.User
		err := db.Where("email = ?", email).First(&user, user).Error
		if err != nil {
			utils.Error("Failed to get user.", err, 0)
			return models.Ticket{}, err
		}
		if !user.Admin {
			return models.Ticket{}, errors.New("Unauthorized access.")
		}
	}
	return ticket, nil
}

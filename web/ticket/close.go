package ticket

import (
	"GoBagouox/database"
	"GoBagouox/database/models"
	"GoBagouox/discord"
	"GoBagouox/discord/service/tickets"
	"GoBagouox/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Close(c *gin.Context) {
	utils.Debug("Execution of /ticket/:id with DELETE", 0)

	ticketID := c.Param("id")
	email := c.Query("email")
	if email == "" {
		utils.Debug("Bad request email missing on Gettranscript", 0)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Some data are missing."})
		return
	}
	db := database.GetDB()
	var ticket models.Ticket
	err := db.Preload("TicketMessages").Preload("Owner").Preload("TicketMessages.Owner").Preload("TicketMessages.TicketAttachments").First(&ticket, ticketID).Error
	if err != nil {
		utils.Error("Failed to get ticket transcript.", err, 0)
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to get ticket transcript."})
		return
	}

	if ticket.Owner.Email != email {
		var user models.User
		err := db.Where("email = ?", email).First(&user, user).Error
		if err != nil {
			utils.Error("Failed to get user.", err, 0)
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to get user."})
			return
		}
		if !user.Admin {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Unauthorized access."})
			return
		}
	}
	s := discord.GetSession()
	if s == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Can't get Discord session."})
		return
	}
	err = tickets.Close(ticket, s)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to get close the ticket."})
	}

}

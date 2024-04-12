package ticket

import (
	"GoBagouox/discord"
	"GoBagouox/discord/service/tickets"
	"GoBagouox/utils"
	"GoBagouox/utils/web"
	"github.com/gin-gonic/gin"
	"net/http"
)

type closeProps struct {
	Email string `json:"email" binding:"required"`
}

func Close(c *gin.Context) {
	utils.Debug("Execution of /ticket/:id with DELETE", 0)

	ticketID := c.Param("id")
	var msg closeProps
	if err := c.BindJSON(&msg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Some data are missing."})
		return
	}

	ticket, err := web.GetTicket(msg.Email, ticketID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err})
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

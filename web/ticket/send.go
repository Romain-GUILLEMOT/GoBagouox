package ticket

import (
	"GoBagouox/database"
	"GoBagouox/database/models"
	"GoBagouox/discord"
	"GoBagouox/discord/service/tickets"
	"GoBagouox/utils"
	"GoBagouox/utils/security"
	"GoBagouox/utils/web"
	"github.com/gin-gonic/gin"
	"net/http"
)

type sendProps struct {
	Email   string `json:"email" binding:"required"`
	Content string `json:"content" binding:"required"`
	Salt    string `json:"salt" binding:"required"`
}

func SendMessage(c *gin.Context) {
	utils.Debug("Execution of /ticket/:id with POST", 0)
	db := database.GetDB()
	ticketID := c.Param("id")
	var msg sendProps
	if err := c.BindJSON(&msg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Some data are missing."})
		return
	}

	if msg.Email == "" || ticketID == "" || msg.Content == "" || msg.Salt == "" {
		utils.Debug("Bad request query parameters missing on Gettranscript", 0)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Some data are missing."})
		return
	}

	var user models.User
	err := db.Where("email = ?", msg.Email).First(&user).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to get user."})
		return
	}
	var ticket models.Ticket
	ticket, err = web.GetTicket(msg.Email, ticketID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err})
		return
	}

	s := discord.GetSession()
	if s == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Can't get Discord session."})
		return
	}
	decryptedContent, err := security.Decrypt(msg.Content, msg.Salt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to decrypt content."})
		return
	}
	err = tickets.SendMessageToTicket(ticket, s, decryptedContent, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to get close the ticket."})
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Message sent successfully."})
}

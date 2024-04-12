package ticket

import (
	"GoBagouox/utils"
	"GoBagouox/utils/security"
	"GoBagouox/utils/web"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Gettranscript(c *gin.Context) {
	utils.Debug("Execution of /ticket/:id with GET", 0)

	ticketID := c.Param("id")
	email := c.Query("email")
	if email == "" {
		utils.Debug("Bad request email missing on Gettranscript", 0)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Some data are missing."})
		return
	}
	ticket, err := web.GetTicket(email, ticketID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err})
	}
	messages := make([]gin.H, len(ticket.TicketMessages))
	for i, message := range ticket.TicketMessages {
		attachments := make([]gin.H, len(message.TicketAttachments))
		for j, attachment := range message.TicketAttachments {
			attachments[j] = gin.H{
				"uuid": attachment.Uuid,
				"type": attachment.Type,
				"size": attachment.Size,
				"name": attachment.Name,
			}
		}
		messages[i] = gin.H{
			"content":      message.Content,
			"id":           message.ID,
			"owner":        message.Owner.Username,
			"owner_avatar": message.Owner.Avatar,
			"created_at":   message.CreatedAt,
			"updated_at":   message.UpdatedAt,
			"admin":        message.Owner.Admin,
			"attachments":  attachments,
		}
	}
	data := gin.H{
		"id":         ticket.ID,
		"name":       ticket.Name,
		"status":     ticket.Status,
		"created_at": ticket.CreatedAt,
		"updated_at": ticket.UpdatedAt,
		"license":    ticket.License,
		"logs":       ticket.Logs,
		"messages":   messages,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		utils.Error("Error during JSON conversion.", err, 0)
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Unknown error occurred."})
		return
	}
	encryptedText, saltBase64, err := security.Encrypt(string(jsonData))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Unknown error occurred."})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "salt": saltBase64, "data": encryptedText})

}

package ticket

import (
	"GoBagouox/database"
	"GoBagouox/database/models"
	"GoBagouox/utils"
	"GoBagouox/utils/security"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Gettranscript(c *gin.Context) {
	ticketID := c.Param("id")
	db := database.GetDB()
	var ticket models.Ticket
	err := db.Preload("TicketMessages").Preload("TicketMessages.Owner").Preload("TicketMessages.TicketAttachments").First(&ticket, ticketID).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to get ticket transcript."})
		utils.Error("Failed to get ticket transcript.", err, 0)
		return
	}
	//Retouner un json avec les messages et attachement (Tout en incluant l'id discord)
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
			"attachments":  attachments,
		}
	}
	data := gin.H{
		"id":       ticket.ID,
		"name":     ticket.Name,
		"status":   ticket.Status,
		"license":  ticket.License,
		"logs":     ticket.Logs,
		"messages": messages,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Unknown error occurred."})
		utils.Error("Error during JSON conversion.", err, 0)
		return
	}
	salt, err := security.CreateSalt()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Unknown error occurred."})
		utils.Error("Error during salt creation.", err, 0)
		return
	}
	encryptedKey, err := security.PBKDF2Encode(salt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Unknown error occurred."})
		utils.Error("Error during PBKDF2 creation.", err, 0)
		return
	}
	encryptedText, err := security.EncryptAES(string(jsonData), encryptedKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Unknown error occurred."})
		utils.Error("Error during AES encoding.", err, 0)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "salt": salt, "data": encryptedText})

}

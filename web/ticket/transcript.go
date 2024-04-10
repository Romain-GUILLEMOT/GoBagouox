package ticket

import (
	"GoBagouox/database"
	"GoBagouox/database/models"
	"GoBagouox/utils"
	"GoBagouox/utils/security"
	"encoding/base64"
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
	salt, err := security.CreateSalt()
	if err != nil {
		utils.Error("Error during salt creation.", err, 0)
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Unknown error occurred."})
		return
	}
	encryptedKey, err := security.PBKDF2Encode(salt)
	if err != nil {
		utils.Error("Error during PBKDF2 creation.", err, 0)
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Unknown error occurred."})
		return
	}
	encryptedText, err := security.EncryptXChaCha(string(jsonData), encryptedKey)
	if err != nil {
		utils.Error("Error during XCha encoding.", err, 0)
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Unknown error occurred."})
		return
	}
	saltBase64 := base64.StdEncoding.EncodeToString(salt)

	c.JSON(http.StatusOK, gin.H{"status": "success", "salt": saltBase64, "data": encryptedText})

}

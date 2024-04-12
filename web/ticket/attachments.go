package ticket

import (
	"GoBagouox/database"
	"GoBagouox/database/models"
	"GoBagouox/utils"
	"GoBagouox/utils/web"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path/filepath"
)

func Downloadattachment(c *gin.Context) {
	utils.Debug("Execution of /ticket/:id/attachment/:uuid with GET", 0)
	email := c.Query("email")
	ticketID := c.Param("id")
	attachmentUUID := c.Param("uuid")
	if email == "" {
		utils.Debug("Bad request email missing on Gettranscript", 0)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Some data are missing."})
		return
	}
	db := database.GetDB()
	var user models.User
	result := db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		utils.Error("Error while querying user", result.Error, 2)
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "User not found."})
		return
	}
	ticket, err := web.GetTicket(email, ticketID)
	if err != nil {
		utils.Error("Error while querying ticket", result.Error, 2)

		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err})
	}
	var attachment models.TicketAttachments
	err = db.Where("uuid = ?", attachmentUUID).Preload("TicketMessage").First(&attachment).Error
	if err != nil {
		utils.Error("Error while querying attachment", result.Error, 2)

		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err})
	}
	if attachment.TicketMessage.TicketID != ticket.ID {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Unauthorized access."})
		return
	}
	cleanUploadedFilesPath := os.Getenv("WEBSERVER_UPLOADED_FILES")

	file := filepath.Join(cleanUploadedFilesPath, attachmentUUID+"."+attachment.Type)
	if _, err := os.Stat(file); err == nil {
		c.File(file)
		return
	} else {
		utils.Error("Error while querying attachment file", errors.New(file), 2)

		c.JSON(http.StatusBadRequest, gin.H{"error": "File not found"})
		return
	}

}

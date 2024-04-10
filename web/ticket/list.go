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

func Getticketlist(c *gin.Context) {
	utils.Debug("Execution of /tickets with GET", 0)
	email := c.Query("email")
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
	var tickets []models.Ticket
	result = db.Where("owner_id = ?", user.ID).Find(&tickets)
	if result.Error != nil {
		utils.Error("Error while querying tickets", result.Error, 2)
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Internal server error."})
		return
	}
	ticketsJSON, err := json.Marshal(tickets)
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
	encryptedText, err := security.EncryptXChaCha(string(ticketsJSON), encryptedKey)
	if err != nil {
		utils.Error("Error during XCha encoding.", err, 0)
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Unknown error occurred."})
		return
	}
	saltBase64 := base64.StdEncoding.EncodeToString(salt)

	c.JSON(http.StatusOK, gin.H{"status": "success", "salt": saltBase64, "data": encryptedText})
}

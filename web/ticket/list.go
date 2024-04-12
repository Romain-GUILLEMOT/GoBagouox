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
	if user.Admin {
		result = db.Find(&tickets)
		if result.Error != nil {
			utils.Error("Error while querying tickets", result.Error, 2)
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Internal server error."})
			return
		}
	} else {
		result = db.Where("owner_id = ?", user.ID).Find(&tickets)
		if result.Error != nil {
			utils.Error("Error while querying tickets", result.Error, 2)
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Internal server error."})
			return
		}
	}

	ticketsJSON, err := json.Marshal(tickets)
	if err != nil {
		utils.Error("Error during JSON conversion.", err, 0)
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Unknown error occurred."})
		return
	}

	encryptedText, saltBase64, err := security.Encrypt(string(ticketsJSON))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Unknown error occurred."})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "salt": saltBase64, "data": encryptedText})
}

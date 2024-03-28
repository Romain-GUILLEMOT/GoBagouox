package tickets

import (
	"GoBagouox/database"
	"GoBagouox/database/models"
	"GoBagouox/utils"
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func NewTicketAttachment(s *discordgo.Session, m *discordgo.MessageCreate) {
	for _, attachment := range m.Message.Attachments {
		urlwithoutQuerry := strings.Split(attachment.URL, "?")[0]
		file := strings.Split(urlwithoutQuerry, "/")[6]
		filename := strings.Split(file, ".")[0]
		filetype := strings.Split(file, ".")[1]
		utils.Info("Download a new file: "+file, 1)
		uuidv4 := uuid.New().String()
		size := downloadFile(attachment.URL, uuidv4, filetype)
		if size == 0 {
			utils.Error("Error during file download.", errors.New("DowloadFile function error"), 1)
			return
		}
		db := database.GetDB()
		var message models.TicketMessage
		result := db.Where("message_id = ?", m.ID).First(&message)
		if result.Error != nil {
			utils.Error("Can't retrieve a message", result.Error, 1)
			return
		}
		ticketattachement := models.TicketAttachments{
			Model:         gorm.Model{},
			TicketMessage: message,
			Uuid:          uuidv4,
			Type:          filetype,
			Size:          size,
			Name:          filename,
		}
		db.Create(&ticketattachement)
		utils.Info("File downloaded!", 1)
	}

}

func downloadFile(url string, uuid string, filetype string) int64 {
	resp, err := http.Get(url)
	if err != nil {
		utils.Error("Error while downloading"+url, err, 1)
		return 0
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			utils.Error("Error while downloading (Can't close the session)"+url, err, 1)
		}
	}(resp.Body)
	cleanUploadedFilesPath := os.Getenv("WEBSERVER_UPLOADED_FILES")

	out, err := os.Create(filepath.Join(cleanUploadedFilesPath, fmt.Sprintf("%s.%s", uuid, filetype)))
	if err != nil {
		utils.Error("Error creating a file.", err, 1)
		return 0
	}
	defer func(out *os.File) {
		err := out.Close()
		if err != nil {
			utils.Error("Error closing a file.", err, 1)
		}
	}(out)
	n, err := io.Copy(out, resp.Body)
	if err != nil {
		utils.Error("Error while copying data to file.", err, 1)
		return 0
	}
	return n
}

package business

import "github.com/ONEST-Network/Whatsapp-Chatbot/bpp/backend/pkg/database/mongodb/business"

type AddBusinessRequest struct {
	Name           string            `json:"name"`
	Phone          string            `json:"phone"`
	Email          string            `json:"email"`
	PictureURLs    []string          `json:"pictureUrls"`
	Description    string            `json:"description"`
	GSTIndexNumber string            `json:"gstIndexNumber"`
	Location       business.Location `json:"location"`
	Industry       business.Industry `json:"industry"`
}

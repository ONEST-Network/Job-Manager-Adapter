package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ONEST-Network/Whatsapp-Chatbot/bpp/backend/internal/business"
	"github.com/ONEST-Network/Whatsapp-Chatbot/bpp/backend/pkg/clients"
	businessPayload "github.com/ONEST-Network/Whatsapp-Chatbot/bpp/backend/pkg/types/payload/business"
	"github.com/gin-gonic/gin"
)

// @Summary	Add business
// @Description	Add a business
// @Tags Business
// @Accept		json
// @Produce		json
// @Param request body businessPayload.AddBusinessRequest true "request body"
// @Success 200 {string} string
// @Failure 500 {object} string
// @Router	/business/add	[post]
func AddBusiness(clients *clients.Clients) gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload businessPayload.AddBusinessRequest
		if err := json.NewDecoder(c.Request.Body).Decode(&payload); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
			return
		}

		businessID, err := business.NewBusiness(clients).AddBusiness(&payload)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, businessID)
	}
}

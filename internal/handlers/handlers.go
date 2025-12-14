package handlers

import (
	"lunch_menu/internal/database"
	"lunch_menu/internal/models"
	"lunch_menu/internal/utils"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// GetAPIInfo godoc
// @Summary      Get API information
// @Description  Returns basic API information and version
// @Tags         info
// @Produce      json
// @Success      200  {object}  models.StandardResponse
// @Router       /api [get]
func GetAPIInfo(c *gin.Context) {
	version := os.Getenv("VERSION")
	if version == "" {
		version = "development"
	}

	response := models.APIResponse{
		Message: "Restaurant Management API",
		Version: version,
	}
	utils.Respond(c, http.StatusOK, "API info fetched successfully", response, nil)
}

// GetBusinessStatistics godoc
// @Summary      Get business statistics
// @Description  Returns business analytics and statistics
// @Tags         statistics
// @Produce      json
// @Success      200  {object}  models.StandardResponse
// @Failure      500  {object}  models.StandardResponse
// @Router       /api/statistics [get]
func GetBusinessStatistics(c *gin.Context) {
	stats, err := database.GetBusinessStatistics()
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, "Failed to retrieve business statistics", nil, &models.ErrorResponse{
			Error:   "DATABASE_ERROR",
			Message: "Failed to retrieve business statistics",
		})
		return
	}

	utils.Respond(c, http.StatusOK, "Business statistics fetched successfully", stats, nil)
}

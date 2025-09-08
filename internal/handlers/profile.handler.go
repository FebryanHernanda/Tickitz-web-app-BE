package handlers

import (
	"net/http"
	"strconv"

	"github.com/FebryanHernanda/Tickitz-web-app-BE/internal/repositories"
	"github.com/gin-gonic/gin"
)

type ProfileHandler struct {
	repo *repositories.ProfileRepository
}

func NewProfileHandler(repo *repositories.ProfileRepository) *ProfileHandler {
	return &ProfileHandler{
		repo: repo,
	}
}

// GetProfile godoc
// @Summary      Get user profile by ID
// @Description  Retrieve user profile by user ID
// @Tags         Profile
// @Produce      json
// @Param        userID path int true "User ID"
// @Success      200    {object} models.SuccessResponse
// @Failure      400    {object} models.ErrorResponse
// @Failure      404    {object} models.ErrorResponse
// @Router       /profile/{userID} [get]
func (h *ProfileHandler) GetProfile(ctx *gin.Context) {
	userIDParams := ctx.Param("userID")
	userID, err := strconv.Atoi(userIDParams)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid user ID",
		})
		return
	}

	userProfile, err := h.repo.GetProfile(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "User not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    userProfile,
	})
}

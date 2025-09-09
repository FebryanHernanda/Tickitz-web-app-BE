package handlers

import (
	"net/http"

	"github.com/FebryanHernanda/Tickitz-web-app-BE/internal/models"
	"github.com/FebryanHernanda/Tickitz-web-app-BE/internal/repositories"
	"github.com/FebryanHernanda/Tickitz-web-app-BE/internal/utils"
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

// /* // CreateProfile godoc
// // @Summary      Create a new profile
// // @Description  Create a profile for a user. Points are initialized to 0 and default image is set.
// // @Tags         Profile
// // @Accept       json
// // @Produce      json
// // @Security BearerAuth
// // @Param        body  body      models.CreateProfile  true  "Profile data"
// // @Success      200   {object}  map[string]interface{}  "Profile created successfully"
// // @Failure      400   {object}  map[string]interface{}  "Bad Request"
// // @Failure      500   {object}  map[string]interface{}  "Internal Server Error"
// // @Router       /profile [post]
// func (h *ProfileHandler) CreateProfile(ctx *gin.Context) {
// 	var profile models.CreateProfile

// 	if err := ctx.ShouldBind(&profile); err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{
// 			"status": false,
// 			"error":  err.Error(),
// 		})
// 		return
// 	}

// 	profile.Points = 0
// 	profile.ImagePath = "public/profile/default.png"
// 	err := h.repo.CreateProfile(ctx, &profile)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{
// 			"status": false,
// 			"error":  err.Error(),
// 		})
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, gin.H{
// 		"status":  true,
// 		"message": "Profile created",
// 	})
// } */

// UpdateProfile godoc
// @Summary      Update user's profile
// @Description  Update the profile of the logged-in user. Supports optional file upload for image.
// @Tags         Profile
// @Accept       multipart/form-data
// @Produce      json
// @Param        first_name   formData  string  false  "First Name"
// @Param        last_name    formData  string  false  "Last Name"
// @Param        phone_number formData  string  false  "Phone Number"
// @Param        image        formData  file    false  "Profile Image"
// @Success      200 {object} map[string]interface{} "Profile updated successfully"
// @Failure      400 {object} map[string]interface{} "Bad Request"
// @Failure      401 {object} map[string]interface{} "Unauthorized"
// @Failure      500 {object} map[string]interface{} "Internal Server Error"
// @Security     BearerAuth
// @Router       /profile [patch]
func (h *ProfileHandler) UpdateProfile(ctx *gin.Context) {
	rawClaims, _ := ctx.Get("claims")
	claims := rawClaims.(*utils.Claims)

	userID := claims.UserID

	var profileUpdate models.ProfileUpdate
	if err := ctx.ShouldBind(&profileUpdate); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
	}

	// file upload
	savePath, err := utils.UploadFile(ctx, "image", "public/profile", "profile")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": false,
			"error":  "Failed to upload file",
		})
		return
	}

	if savePath != "" {
		profileUpdate.ImagePath = &savePath
	}

	if err := h.repo.UpdateProfile(ctx, userID, &profileUpdate); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": false,
			"error":  err.Error(),
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Profile updated successfully",
	})
}

// GetProfile godoc
// @Summary      Get user profile by ID
// @Description  Retrieve user profile by user ID
// @Tags         Profile
// @Security     BearerAuth
// @Produce      json
// @Success      200    {object} models.SuccessResponse
// @Failure      400    {object} models.ErrorResponse
// @Failure      401    {object} models.ErrorResponse   "Unauthorized or invalid token"
// @Failure      404    {object} models.ErrorResponse	"User not found"
// @Router       /profile [get]
func (h *ProfileHandler) GetProfile(ctx *gin.Context) {
	rawClaims, _ := ctx.Get("claims")
	claims := rawClaims.(*utils.Claims)

	userID := claims.UserID

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

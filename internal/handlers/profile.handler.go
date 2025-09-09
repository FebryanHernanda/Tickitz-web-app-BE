package handlers

import (
	"log"
	"net/http"

	"github.com/FebryanHernanda/Tickitz-web-app-BE/internal/models"
	"github.com/FebryanHernanda/Tickitz-web-app-BE/internal/repositories"
	"github.com/FebryanHernanda/Tickitz-web-app-BE/internal/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type ProfileHandler struct {
	repo *repositories.ProfileRepository
}

func NewProfileHandler(repo *repositories.ProfileRepository) *ProfileHandler {
	return &ProfileHandler{
		repo: repo,
	}
}

// UpdateProfile godoc
// @Summary      Update user's profile
// @Description  Update the profile of the logged-in user. Supports optional file upload for image.
// @Tags         Profile
// @Accept       multipart/form-data
// @Produce      json
// @Param        email   formData  string  false  "Email"
// @Param        password  formData  string  false  "Password"
// @Param        first_name   formData  string  false  "First Name"
// @Param        last_name    formData  string  false  "Last Name"
// @Param        phone_number formData  string  false  "Phone Number"
// @Param        image        formData  file    false  "Profile Image"
// @Success      200 {object} map[string]interface{} "Profile updated successfully"
// @Failure      400 {object} map[string]interface{} "Bad Request"
// @Failure      401 {object} map[string]interface{} "Unauthorized"
// @Failure      500 {object} map[string]interface{} "Internal Server Error"
// @Security     BearerAuth
// @Router       /profile/edit [patch]
func (h *ProfileHandler) UpdateProfile(ctx *gin.Context) {
	rawClaims, _ := ctx.Get("claims")
	claims := rawClaims.(*utils.Claims)

	userID := claims.UserID

	// var userUpdate models.UserUpdate
	// var profileUpdate models.ProfileUpdate
	var update models.UserUpdateRequest

	if err := ctx.ShouldBind(&update); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
	}

	if update.User.Email != nil {
		if err := utils.IsValidEmail(*update.User.Email); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}
	}

	if update.User.Password != nil {
		if err := utils.IsValidPassword(*update.User.Password); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(*update.User.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "failed to hash password",
		})
		return
	}
	hashPassStr := string(hashedPass)
	update.User.Password = &hashPassStr

	if err := h.repo.UpdateProfile(ctx, userID, &update); err != nil {
		log.Printf("%s", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": false,
			"error":  err.Error(),
		})
		return
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
		update.Profile.ImagePath = &savePath
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

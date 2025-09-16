package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/FebryanHernanda/Tickitz-web-app-BE/internal/models"
	"github.com/FebryanHernanda/Tickitz-web-app-BE/internal/repositories"
	"github.com/FebryanHernanda/Tickitz-web-app-BE/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

type ProfileHandler struct {
	repo *repositories.ProfileRepository
	rdb  *redis.Client
}

func NewProfileHandler(repo *repositories.ProfileRepository, rdb *redis.Client) *ProfileHandler {
	return &ProfileHandler{
		repo: repo,
		rdb:  rdb,
	}
}

// UpdateProfile godoc
// @Summary      Update user's profile
// @Description  Update the profile of the logged-in user. Supports optional file upload for image.
// @Tags         Profile
// @Accept       multipart/form-data
// @Produce      json
// @Param        email   formData  string  false  "Email"
// @Param        old_password  formData  string  false  "Old Password"
// @Param        new_password  formData  string  false  "New Password"
// @Param        first_name   formData  string  false  "First Name"
// @Param        last_name    formData  string  false  "Last Name"
// @Param        phone_number formData  string  false  "Phone Number"
// @Param        image        formData  file    false  "Profile Image"
// @Success      200 {object} models.SuccessResponse "Profile updated successfully"
// @Failure      400 {object} models.ErrorResponse "Bad Request"
// @Failure      401 {object} models.ErrorResponse "Unauthorized"
// @Failure      500 {object} models.ErrorResponse "Internal Server Error"
// @Security     BearerAuth
// @Router       /profile/edit [patch]
func (h *ProfileHandler) UpdateProfile(ctx *gin.Context) {
	rawClaims, _ := ctx.Get("claims")
	claims := rawClaims.(*utils.Claims)

	userID := claims.UserID

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

	if update.User.NewPassword != nil {
		if update.User.OldPassword == nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "old password is required",
			})
			return
		}

		storedUser, err := h.repo.GetProfile(ctx, userID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "failed to get user data",
			})
			return
		}

		if bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(*update.User.OldPassword)) != nil {
			log.Println("bcrypt error:", err)
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "old password is incorrect",
			})
			return
		}

		if err := utils.IsValidPassword(*update.User.NewPassword); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}

		hashedPass, _ := bcrypt.GenerateFromPassword([]byte(*update.User.NewPassword), bcrypt.DefaultCost)
		hashPassStr := string(hashedPass)
		update.User.NewPassword = &hashPassStr
	}

	// file upload
	savePath, err := utils.UploadFile(ctx, "image", "public/profile", "profile", "profile")
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

	if err := h.repo.UpdateProfile(ctx, userID, &update); err != nil {
		log.Printf("%s", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": false,
			"error":  err.Error(),
		})
		return
	}

	if err := utils.InvalidateCache(ctx, h.rdb, []string{"users:"}); err != nil {
		log.Println("Redis delete cache error:", err)
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

	redisKey := fmt.Sprintf("users:profile=%d", userID)
	var cached models.Profile

	if h.rdb != nil {
		err := utils.GetCache(ctx, h.rdb, redisKey, &cached)
		if err != nil {
			log.Panicln("Redis error, back to DB : ", err)
		}
		if cached.ID != 0 {
			ctx.JSON(http.StatusOK, gin.H{
				"success": true,
				"data":    cached,
				"message": "data from cache",
			})
			return
		}
	}

	userProfile, err := h.repo.GetProfile(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "User not found",
		})
		return
	}

	if h.rdb != nil {
		err := utils.SetCache(ctx, h.rdb, redisKey, userProfile, 10*time.Minute)
		if err != nil {
			log.Println("Redis set cache error:", err)
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "data from database",
		"data":    userProfile,
	})
}

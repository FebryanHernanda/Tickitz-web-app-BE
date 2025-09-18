package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/FebryanHernanda/Tickitz-web-app-BE/internal/models"
	"github.com/FebryanHernanda/Tickitz-web-app-BE/internal/repositories"
	"github.com/FebryanHernanda/Tickitz-web-app-BE/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	repo       *repositories.UserRepository
	JWTManager *utils.JWTManager
	rdb        *redis.Client
}

func NewAuthHandler(repo *repositories.UserRepository, jwtManager *utils.JWTManager, rdb *redis.Client) *AuthHandler {
	return &AuthHandler{
		repo:       repo,
		JWTManager: jwtManager,
		rdb:        rdb,
	}
}

type LoginResponse struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"Login successful"`
	Token   string `json:"token,omitempty" example:"your token..."`
}

// Register godoc
// @Summary      Register a new user
// @Description  Create a new user account with email, password, and role.
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        user body models.RegisterUser true "Register User"
// @Success      200  {object}  models.SuccessResponse
// @Failure      400  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /auth/register [post]
func (u *AuthHandler) Register(ctx *gin.Context) {
	var req models.RegisterUser
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if req.Role == "" {
		req.Role = "user"
	}

	if err := utils.IsValidEmail(req.Email); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if err := utils.IsValidPassword(req.Password); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	va, err := utils.GenerateVirtualAccount()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "failed to generate virtual account",
		})
		return
	}

	user := models.User{
		Email:          req.Email,
		Password:       req.Password,
		Role:           req.Role,
		VirtualAccount: va,
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "failed to hash password",
		})
		return
	}
	user.Password = string(hashedPass)

	if err := u.repo.RegisterUser(ctx, &user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "registration has failed",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    user.Email,
	})
}

// Login godoc
// @Summary      User login
// @Description  Authenticate user and return JWT token
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        user body models.LoginUser true "User login credentials"
// @Success      200  {object}  LoginResponse
// @Failure      401  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /auth/login [post]
func (u *AuthHandler) Login(ctx *gin.Context) {
	var user models.LoginUser

	if err := ctx.ShouldBind(&user); err != nil {
		log.Printf("%s", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	dbUser, err := u.repo.GetUserByEmail(ctx, user.Email)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "invalid email or password",
		})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password)); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "invalid email or password",
		})
		return
	}

	token, err := u.JWTManager.GenerateToken(dbUser)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "failed to generate token",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Login successful",
		"token":   token,
		"role":    dbUser.Role,
		"email":   dbUser.Email,
	})

}

// Logout godoc
// @Summary Logout user and invalidate JWT token
// @Description Invalidate the JWT token by adding it to Redis blacklist so it cannot be used again
// @Tags Authentication
// @Accept json
// @Produce json
// @Security     BearerAuth
// @Success 200 {object} models.SuccessResponse "Logout successful"
// @Failure 401 {object} models.ErrorResponse "Token is required or unauthorized"
// @Failure 500 {object} models.ErrorResponse "Internal server error during logout"
// @Router /auth/logout [post]
func (u *AuthHandler) Logout(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	tokenString = strings.TrimSpace(tokenString)

	redisKey := fmt.Sprintf("blacklist:%s", tokenString)

	if tokenString == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Token is required",
		})
		return
	}

	err := utils.SetCache(ctx, u.rdb, redisKey, true, 30*time.Minute)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "failed to logout",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Logout successful",
	})
}

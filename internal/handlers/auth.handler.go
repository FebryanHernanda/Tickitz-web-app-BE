package handlers

import (
	"net/http"

	"github.com/FebryanHernanda/Tickitz-web-app-BE/internal/models"
	"github.com/FebryanHernanda/Tickitz-web-app-BE/internal/repositories"
	"github.com/FebryanHernanda/Tickitz-web-app-BE/internal/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	repo       *repositories.UserRepository
	JWTManager *utils.JWTManager
}

func NewAuthHandler(repo *repositories.UserRepository, jwtManager *utils.JWTManager) *AuthHandler {
	return &AuthHandler{
		repo:       repo,
		JWTManager: jwtManager,
	}
}

type LoginResponse struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"Login successful"`
	Token   string `json:"token,omitempty" example:"your token..."`
}

// Register godoc
// @Summary      Register a new user
// @Description  Register user with validation, password hashing, and virtual account generation
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        user  body      models.RegisterUser  true  "User registration data"
// @Success      200   {object}  models.SuccessResponse
// @Failure      400   {object}  models.ErrorResponse
// @Failure      500   {object}  models.ErrorResponse
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

	if err := utils.IsValidEmail(req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if err := utils.IsValidPassword(req); err != nil {
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
		"data":    user,
	})
}

// Login godoc
// @Summary      User login
// @Description  Authenticate user and return JWT token
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        user body models.LoginUser true "User login credentials"
// @Success      200  {object} LoginResponse
// @Failure      401  {object} models.ErrorResponse
// @Failure      500  {object} models.ErrorResponse
// @Router       /auth/login [post]
func (u *AuthHandler) Login(ctx *gin.Context) {
	var user models.LoginUser

	if err := ctx.ShouldBind(&user); err != nil {
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
	})

}

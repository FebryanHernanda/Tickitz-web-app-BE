package routers

import (
	"github.com/FebryanHernanda/Tickitz-web-app-BE/internal/handlers"
	"github.com/FebryanHernanda/Tickitz-web-app-BE/internal/middlewares"
	"github.com/FebryanHernanda/Tickitz-web-app-BE/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func AuthRouter(r *gin.Engine, jwtManager *utils.JWTManager, rdb *redis.Client, authHandler *handlers.AuthHandler) {
	authRoutes := r.Group("/auth")
	authRoutes.POST("/register", authHandler.Register)
	authRoutes.POST("/login", authHandler.Login)
	authRoutes.POST("/logout", middlewares.VerifyToken(jwtManager, rdb), middlewares.AuthMiddleware("user", "admin"), authHandler.Logout)
}

package routers

import (
	"github.com/FebryanHernanda/Tickitz-web-app-BE/internal/handlers"
	"github.com/FebryanHernanda/Tickitz-web-app-BE/internal/middlewares"
	"github.com/FebryanHernanda/Tickitz-web-app-BE/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func ProfileRouter(r *gin.Engine, profileHandler *handlers.ProfileHandler, jwtManager *utils.JWTManager, rdb *redis.Client) {
	profileRoutes := r.Group("/profile")
	profileRoutes.Use(middlewares.VerifyToken(jwtManager, rdb))
	profileRoutes.Use(middlewares.AuthMiddleware("user"))

	profileRoutes.GET("", profileHandler.GetProfile)
	profileRoutes.PATCH("/edit", profileHandler.UpdateProfile)
	profileRoutes.PATCH("/editpassword", profileHandler.UpdatePassword)
}

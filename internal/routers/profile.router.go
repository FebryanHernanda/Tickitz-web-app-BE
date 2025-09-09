package routers

import (
	"github.com/FebryanHernanda/Tickitz-web-app-BE/internal/handlers"
	"github.com/FebryanHernanda/Tickitz-web-app-BE/internal/middlewares"
	"github.com/FebryanHernanda/Tickitz-web-app-BE/internal/utils"
	"github.com/gin-gonic/gin"
)

func ProfileRouter(r *gin.Engine, profileHandler *handlers.ProfileHandler, jwtManager *utils.JWTManager) {
	profileRoutes := r.Group("/profile")
	profileRoutes.Use(middlewares.VerifyToken(jwtManager))
	profileRoutes.Use(middlewares.AuthMiddleware("user"))

	profileRoutes.GET("", profileHandler.GetProfile)
	profileRoutes.PATCH("/edit", profileHandler.UpdateProfile)
}

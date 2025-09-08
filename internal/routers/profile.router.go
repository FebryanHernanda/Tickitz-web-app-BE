package routers

import (
	"github.com/FebryanHernanda/Tickitz-web-app-BE/internal/handlers"
	"github.com/gin-gonic/gin"
)

func ProfileRouter(r *gin.Engine, profileHandler *handlers.ProfileHandler) {
	r.GET("/profile/:userID", profileHandler.GetProfile)
}

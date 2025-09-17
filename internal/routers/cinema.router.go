package routers

import (
	"github.com/FebryanHernanda/Tickitz-web-app-BE/internal/handlers"
	"github.com/gin-gonic/gin"
)

func CinemaRouter(r *gin.Engine, cinemaHandler *handlers.CinemaHandler) {
	cinemaRoutes := r.Group("/cinemas")

	cinemaRoutes.GET("/available-seats/:cinemas_schedule_id", cinemaHandler.GetAvailableSeats)
	cinemaRoutes.GET("/:movieID", cinemaHandler.GetScheduleFilter)
}

package routers

import (
	"github.com/FebryanHernanda/Tickitz-web-app-BE/internal/handlers"
	"github.com/gin-gonic/gin"
)

func SeatsRouter(r *gin.Engine, seatHandler *handlers.SeatHandler) {
	r.GET("cinemas/available-seats/:cinemas_schedule_id", seatHandler.GetAvailableSeats)
}

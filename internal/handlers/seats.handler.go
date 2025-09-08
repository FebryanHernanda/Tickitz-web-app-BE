package handlers

import (
	"net/http"
	"strconv"

	"github.com/FebryanHernanda/Tickitz-web-app-BE/internal/repositories"
	"github.com/gin-gonic/gin"
)

type SeatHandler struct {
	Repo *repositories.SeatRepository
}

func NewSeatHandler(repo *repositories.SeatRepository) *SeatHandler {
	return &SeatHandler{
		Repo: repo,
	}
}

// GetAvailableSeats godoc
// @Summary      Get available seats for a cinema schedule
// @Description  Retrieve list of available seats by cinema schedule ID
// @Tags         Cinemas
// @Produce      json
// @Param        cinemas_schedule_id path int true "Cinema Schedule ID"
// @Success      200  {object} models.SuccessResponse
// @Failure      400  {object} models.ErrorResponse
// @Failure      500  {object} models.ErrorResponse
// @Router       /cinemas/available-seats/{cinemas_schedule_id} [get]
func (h *SeatHandler) GetAvailableSeats(ctx *gin.Context) {
	cinemaSchedulesIDStr := ctx.Param("cinemas_schedule_id")
	cinemaSchedulesID, err := strconv.Atoi(cinemaSchedulesIDStr)
	if err != nil || cinemaSchedulesID <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid cinemas_schedule_id",
		})
		return
	}

	seats, err := h.Repo.GetAvailableSeats(ctx, cinemaSchedulesID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    seats,
	})
}

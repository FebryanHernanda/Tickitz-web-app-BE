package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/FebryanHernanda/Tickitz-web-app-BE/internal/models"
	"github.com/FebryanHernanda/Tickitz-web-app-BE/internal/repositories"
	"github.com/gin-gonic/gin"
)

type CinemaHandler struct {
	repo *repositories.CinemaRepository
}

func NewCinemaHandler(repo *repositories.CinemaRepository) *CinemaHandler {
	return &CinemaHandler{
		repo: repo,
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
func (h *CinemaHandler) GetAvailableSeats(ctx *gin.Context) {
	cinemaSchedulesIDStr := ctx.Param("cinemas_schedule_id")
	cinemaSchedulesID, err := strconv.Atoi(cinemaSchedulesIDStr)
	if err != nil || cinemaSchedulesID <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid cinemas_schedule_id",
		})
		return
	}

	exist, err := h.repo.IsCinemaScheduleExists(ctx, cinemaSchedulesID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if !exist {
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Cinema schedule not found",
		})
		return
	}

	seats, err := h.repo.GetAvailableSeats(ctx, cinemaSchedulesID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if len(seats) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "All seats are available",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    seats,
	})
}

// GetScheduleFilter godoc
// @Summary      Get cinema schedules by filter
// @Description  Retrieve cinema schedules filtered by location, date, and time with pagination
// @Tags         Cinemas
// @Produce      json
// @Param        location  query  string  false  "Location Filter"
// @Param        date      query  string  false  "Date Filter (YYYY-MM-DD)"
// @Param        time      query  string  false  "Time Filter (HH:MM)"
// @Success      200  {object}  models.SuccessResponse
// @Failure      400  {object}  models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /cinemas [get]
func (h *CinemaHandler) GetScheduleFilter(ctx *gin.Context) {
	location := ctx.Query("location")
	dateStr := ctx.Query("date")
	timeStr := ctx.Query("time")

	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil || page < 1 {
		page = 1
	}
	limit := 5
	offset := (page - 1) * limit

	var filter models.GetFilterSchedules
	if location != "" {
		filter.LocationName = &location
	}

	if dateStr != "" {
		parsedDate, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "invalid date format, must be YYYY-MM-DD",
			})
			return
		}
		filter.ScheduleDate = &parsedDate
	}

	if timeStr != "" {
		filter.ScheduleTime = &timeStr
	}

	schedule, err := h.repo.GetScheduleFilter(ctx, filter.LocationName, filter.ScheduleDate, filter.ScheduleTime, limit, offset)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if len(schedule) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "No schedules found",
			"page":    page,
			"limit":   limit,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    schedule,
		"page":    page,
		"limit":   limit,
	})
}

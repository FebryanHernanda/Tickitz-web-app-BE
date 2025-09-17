package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/FebryanHernanda/Tickitz-web-app-BE/internal/models"
	"github.com/FebryanHernanda/Tickitz-web-app-BE/internal/repositories"
	"github.com/FebryanHernanda/Tickitz-web-app-BE/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type CinemaHandler struct {
	repo *repositories.CinemaRepository
	rdb  *redis.Client
}

func NewCinemaHandler(repo *repositories.CinemaRepository, rdb *redis.Client) *CinemaHandler {
	return &CinemaHandler{
		repo: repo,
		rdb:  rdb,
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

	redisKey := fmt.Sprintf("cinemas:available-seats:%d", cinemaSchedulesID)
	var cached []models.CinemaSeat

	if h.rdb != nil {
		err := utils.GetCache(ctx, h.rdb, redisKey, &cached)
		if err != nil {
			log.Println("Redis error, back to DB : ", err)
		}
		if len(cached) != 0 {
			ctx.JSON(http.StatusOK, gin.H{
				"success": true,
				"data":    cached,
				"message": "data from cache",
			})
			return
		}
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

	if h.rdb != nil {
		err := utils.SetCache(ctx, h.rdb, redisKey, seats, 2*time.Minute)
		if err != nil {
			log.Println("Redis set cache error:", err)
		}
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
// @Param 		 movieID path int true "Movie ID"
// @Param        location  query  string  false  "Location Filter"
// @Param        date      query  string  false  "Date Filter (YYYY-MM-DD)"
// @Param        time      query  string  false  "Time Filter (HH:MM)"
// @Success      200  {object}  models.SuccessResponse
// @Failure      400  {object}  models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /cinemas/{movieID} [get]
func (h *CinemaHandler) GetScheduleFilter(ctx *gin.Context) {
	movieIDStr := ctx.Param("movieID")
	movieID, err := strconv.Atoi(movieIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid movieID parameter",
		})
		return
	}
	log.Printf("movie id : %s", movieIDStr)

	location := ctx.Query("location")
	dateStr := ctx.Query("date")
	timeStr := ctx.Query("time")

	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil || page < 1 {
		page = 1
	}
	limit := 5
	offset := (page - 1) * limit

	locationCache := location
	dateCache := dateStr
	timeCache := timeStr

	if location == "" {
		locationCache = "<empty>"
	}
	if dateStr == "" {
		dateCache = "<empty>"
	}
	if timeStr == "" {
		timeCache = "<empty>"
	}

	redisKey := fmt.Sprintf("cinemas:schedule:movieid=%d:loc=%s:date=%s:time=%s:page=%d", movieID, locationCache, dateCache, timeCache, page)

	var cached []models.GetFilterSchedules
	if h.rdb != nil {
		err := utils.GetCache(ctx, h.rdb, redisKey, &cached)
		if err != nil {
			log.Println("Redis error, back to DB : ", err)
		}
		if len(cached) > 0 {
			ctx.JSON(http.StatusOK, gin.H{
				"success": true,
				"data":    cached,
				"message": "data from cache",
			})
			return
		}
	}

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

	schedule, err := h.repo.GetScheduleFilter(ctx, movieID, filter.LocationName, filter.ScheduleDate, filter.ScheduleTime, limit, offset)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if len(schedule) == 0 {
		if h.rdb != nil {
			err := utils.SetCache(ctx, h.rdb, redisKey, []models.GetFilterSchedules{}, 1*time.Minute)
			if err != nil {
				log.Println("Redis set cache error:", err)
			}
		}
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "No schedules found",
			"page":    page,
			"limit":   limit,
		})
		return
	}

	if h.rdb != nil {
		err := utils.SetCache(ctx, h.rdb, redisKey, schedule, 2*time.Minute)
		if err != nil {
			log.Println("Redis set cache error:", err)
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "data from database",
		"data":    schedule,
		"page":    page,
		"limit":   limit,
	})
}

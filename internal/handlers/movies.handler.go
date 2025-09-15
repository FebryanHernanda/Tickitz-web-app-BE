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

type MoviesHandler struct {
	repo *repositories.MoviesRepository
	rdb  *redis.Client
}

func NewMoviesHandler(repo *repositories.MoviesRepository, rdb *redis.Client) *MoviesHandler {
	return &MoviesHandler{
		repo: repo,
		rdb:  rdb,
	}
}

// GetUpcomingMovies godoc
// @Summary Get upcoming movies
// @Description Retrieve a list of upcoming movies
// @Tags Movies
// @Produce json
// @Success 200 {object} models.SuccessResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /movies/upcoming [get]
func (h *MoviesHandler) GetUpcomingMovies(ctx *gin.Context) {
	redisKey := "movies-upcoming"
	var cached []models.Movie

	if h.rdb != nil {
		err := utils.GetCache(ctx, h.rdb, redisKey, &cached)
		if err != nil {
			log.Panicln("Redis error, back to DB : ", err)
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

	movies, err := h.repo.GetUpcomingMovies(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if len(movies) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"data":    []models.Movie{},
			"message": "No upcoming movies found",
		})
		return
	}

	if h.rdb != nil {
		err := utils.SetCache(ctx, h.rdb, redisKey, movies, 5*time.Minute)
		if err != nil {
			log.Println("Redis set cache error:", err)
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    movies,
	})
}

// GetPopularMovies godoc
// @Summary Get popular movies
// @Description Retrieve a list of popular movies
// @Tags Movies
// @Produce json
// @Success 200 {object} models.SuccessResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /movies/popular [get]
func (h *MoviesHandler) GetPopularMovies(ctx *gin.Context) {
	redisKey := "movies-popular"
	var cached []models.Movie

	if h.rdb != nil {
		// Check Cache
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

	// Cache Miss
	movies, err := h.repo.GetPopularMovies(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
	}

	if len(movies) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"data":    []models.Movie{},
			"message": "No popular movies found",
		})
		return
	}

	if h.rdb != nil {
		// set cache
		err := utils.SetCache(ctx, h.rdb, redisKey, movies, 5*time.Minute)
		if err != nil {
			log.Println("Redis set cache error:", err)
		}

	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    movies,
		"message": "data from database",
	})
}

// GetMoviesByFilter godoc
// @Summary Get movies by filter
// @Description Retrieve movies filtered by search keyword and genre with pagination
// @Tags Movies
// @Produce json
// @Param search query string false "Search keyword"
// @Param genre query string false "Genre filter"
// @Param page query int false "Page number (default 1)"
// @Success 200 {object} models.SuccessResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /movies [get]
func (h *MoviesHandler) GetMoviesByFilter(ctx *gin.Context) {
	search := ctx.Query("search")
	genre := ctx.Query("genre")

	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil || page < 1 {
		page = 1
	}

	limit := 10

	offset := (page - 1) * limit

	cacheSearch := search
	cacheGenre := genre
	if search == "" {
		cacheSearch = "<empty>"
	}
	if genre == "" {
		cacheGenre = "<empty>"
	}

	redisKey := fmt.Sprintf("movies-search:title=%s-genre=%s-page=%d", cacheSearch, cacheGenre, page)
	var cached []models.Movie

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

	movies, err := h.repo.GetMoviesByFilter(ctx, search, genre, page, limit, offset)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if len(movies) == 0 {
		if h.rdb != nil {
			err := utils.SetCache(ctx, h.rdb, redisKey, []models.Movie{}, 10*time.Minute)
			if err != nil {
				log.Println("Redis set cache error:", err)
			}
		}

		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"data":    []models.Movie{},
			"message": "No movies found",
			"page":    page,
			"limit":   limit,
		})
		return
	}

	if h.rdb != nil {
		err := utils.SetCache(ctx, h.rdb, redisKey, movies, 10*time.Minute)
		if err != nil {
			log.Println("Redis set cache error:", err)
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "data from database",
		"page":    page,
		"limit":   limit,
		"count":   len(movies),
		"data":    movies,
	})
}

// GetDetailMovies godoc
// @Summary Get movie detail
// @Description Retrieve detailed information of a movie by ID
// @Tags Movies
// @Produce json
// @Param id path int true "Movie ID"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /movies/{id}/details [get]
func (h *MoviesHandler) GetDetailMovies(ctx *gin.Context) {
	idParam := ctx.Param("id")

	movieID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil || movieID < 1 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid movie ID",
		})
		return
	}

	movies, err := h.repo.GetDetailMovies(ctx, movieID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Movie not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    movies,
	})
}

// GetSchedulesMovies godoc
// @Summary Get movie schedules
// @Description Retrieve the schedule of movies
// @Tags Movies
// @Produce json
// @Success 200 {object} models.SuccessResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /movies/schedules [get]
func (h *MoviesHandler) GetSchedulesMovies(ctx *gin.Context) {
	redisKey := "movies-schedules"
	var cached []models.MovieSchedules

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

	schedules, err := h.repo.GetSchedulesMovies(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if len(schedules) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "No schedules movies found",
		})
		return
	}

	if h.rdb != nil {
		err := utils.SetCache(ctx, h.rdb, redisKey, schedules, 5*time.Minute)
		if err != nil {
			log.Println("Redis set cache error:", err)
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "data from database",
		"data":    schedules,
	})
}

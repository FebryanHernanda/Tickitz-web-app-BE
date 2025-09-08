package handlers

import (
	"net/http"
	"strconv"

	"github.com/FebryanHernanda/Tickitz-web-app-BE/internal/models"
	"github.com/FebryanHernanda/Tickitz-web-app-BE/internal/repositories"
	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	repo *repositories.AdminRepository
}

func NewAdminHandler(repo *repositories.AdminRepository) *AdminHandler {
	return &AdminHandler{
		repo: repo,
	}
}

// GetAllMovies godoc
// @Summary      Get all movies
// @Description  Retrieve list of all movies (admin access required)
// @Tags         Admin
// @Security     BearerAuth
// @Produce      json
// @Param        Authorization header string true "Bearer token" default(Bearer <your_token_here>)
// @Success      200  {object}  models.SuccessResponse
// @Failure      401  {object}  models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /admin/movies [get]
func (a *AdminHandler) GetAllMovies(ctx *gin.Context) {

	allMovies, err := a.repo.GetAllMovies(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if len(allMovies) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Movies not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    allMovies,
	})
}

// GetAllMovies godoc
// @Summary      Add Movies
// @Description  Add Movies with all the relations (genres, cast, and schedules)
// @Tags         Admin
// @Security     BearerAuth
// @Produce      json
// @Param        user  body      models.AddMovies  true  "Add Movies data"
// @Success      200  {object}  models.SuccessResponse
// @Failure      400  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /admin/movies/add [post]
func (h *AdminHandler) AddMovies(ctx *gin.Context) {
	var movie models.AddMovies
	if err := ctx.ShouldBind(&movie); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
	}

	movieData, err := h.repo.AddMovies(ctx, &movie)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "movie add successfully",
		"data":    movieData,
	})
}

// GetAllMovies godoc
// @Summary      Get movies schedules for check the id to add cinemas schedule
// @Description  Get movies schedules for check the id to add cinemas schedule
// @Tags         Admin
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  models.SuccessResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /admin/movies/schedule [get]
func (h *AdminHandler) GetMovieSchedule(ctx *gin.Context) {
	schedules, err := h.repo.GetMovieSchedule(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    schedules,
	})
}

// GetAllMovies godoc
// @Summary      Add Cinemas Schedule
// @Description  Add cinemas schedule
// @Tags         Admin
// @Security     BearerAuth
// @Produce      json
// @Param        user  body     []models.CinemaScheduleLocation  true  "Add cinemas and location schedule data"
// @Success      200  {object}  models.SuccessResponse
// @Failure      400  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /admin/movies/cinemaschedule/add [post]
func (h *AdminHandler) AddCinemaSchedule(ctx *gin.Context) {
	var CinemaSchedules []models.CinemaScheduleLocation

	if err := ctx.ShouldBind(&CinemaSchedules); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	err := h.repo.AddCinemaSchedule(ctx, CinemaSchedules)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":          true,
		"message":          "success add cinemas schedule",
		"cinemas_schedule": CinemaSchedules,
	})
}

// DeleteMovies godoc
// @Summary      Delete a movie by ID
// @Description  Delete a movie and all related data by movie ID
// @Tags         Admin
// @Param        id   path      int  true  "Movie ID"
// @Success      200  {object}  models.SuccessResponse
// @Failure      400  {object}  models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /admin/movies/delete/{id} [delete]
func (h *AdminHandler) DeleteMovies(ctx *gin.Context) {
	moveIDStr := ctx.Param("id")
	movieID, err := strconv.Atoi(moveIDStr)
	if err != nil || movieID < 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "invalid movie id",
		})
		return
	}

	err = h.repo.DeleteMovies(ctx, movieID)
	if err != nil {
		if err.Error() == "movie not found" {
			ctx.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"error":   "movie not found",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "movie deleted successfully",
	})
}

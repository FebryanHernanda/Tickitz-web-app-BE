package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/FebryanHernanda/Tickitz-web-app-BE/internal/models"
	"github.com/FebryanHernanda/Tickitz-web-app-BE/internal/repositories"
	"github.com/FebryanHernanda/Tickitz-web-app-BE/internal/utils"
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

// AddMovie godoc
// @Summary      Add New Movie
// @Description  Add Movies with all the relations (genres, cast, and schedules)
// @Tags         Admin
// @Security     BearerAuth
// @Accept       multipart/form-data
// @Produce      json
// @Param        title         formData  string  false  "Movie Title"
// @Param        synopsis      formData  string  false  "Movie Synopsis"
// @Param        release_date  formData  string  false  "Release Date (YYYY-MM-DD)"
// @Param        rating        formData  number  false  "Movie Rating"
// @Param        age_rating    formData  string  false  "Age Rating"
// @Param        duration      formData  int     false  "Duration (minutes)"
// @Param        director_id   formData  int     false  "Director ID"
// @Param        genres        formData  string   false  "Genres [IDs, comma separated]"
// @Param        casts         formData  string   false  "Casts [IDs, comma separated]"
// @Param        schedules     formData  string  false  "Schedules (JSON array: [{}])"
// @Param        poster        formData  file    false  "Poster file"
// @Param        backdrop      formData  file    false  "Backdrop file"
// @Success      200           {object}  models.SuccessResponse
// @Failure      400           {object}  models.ErrorResponse
// @Failure      401  {object}  models.ErrorResponse
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

	posterPath, err := utils.UploadFile(ctx, "poster", "public/movies/posters", "poster")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to upload poster"})
		return
	}
	if posterPath != "" {
		movie.PosterPath = posterPath
	}

	backdropPath, err := utils.UploadFile(ctx, "backdrop", "public/movies/backdrops", "backdrop")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to upload backdrop",
		})
		return
	}
	if backdropPath != "" {
		movie.BackdropPath = backdropPath
	}

	if genresStr := ctx.PostForm("genres"); genresStr != "" {
		var genres []int
		if err := json.Unmarshal([]byte(genresStr), &genres); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "invalid genres format, must be JSON array like [1,2,3]",
			})
			return
		}
		movie.Genres = genres
	}

	if castsStr := ctx.PostForm("casts"); castsStr != "" {
		var casts []int
		if err := json.Unmarshal([]byte(castsStr), &casts); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "invalid casts format, must be JSON array like [1,2,3]",
			})
			return
		}
		movie.Casts = casts
	}

	schedulesStr := ctx.PostForm("schedules")
	if schedulesStr != "" {
		var schedules []models.Schedule
		if err := json.Unmarshal([]byte(schedulesStr), &schedules); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "invalid schedules format",
			})
			return
		}
		movie.Schedules = schedules
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
// @Summary      Get movies schedules
// @Description  Get movies schedules for check the id to add cinemas schedule
// @Tags         Admin
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  models.SuccessResponse
// @Failure      401  {object}  models.ErrorResponse
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

// UpdateMovie godoc
// @Summary      Update movie with file upload
// @Description  Update movie data by ID, allow uploading poster and backdrop
// @Tags         Admin
// @Security     BearerAuth
// @Accept       multipart/form-data
// @Produce      json
// @Param        id   			path      int  true  "Movie ID"
// @Param        title         formData  string  false  "Movie Title"
// @Param        synopsis      formData  string  false  "Movie Synopsis"
// @Param        release_date  formData  string  false  "Release Date (YYYY-MM-DD)"
// @Param        rating        formData  number  false  "Movie Rating"
// @Param        age_rating    formData  string  false  "Age Rating"
// @Param        duration      formData  int     false  "Duration (minutes)"
// @Param        director_id   formData  int     false  "Director ID"
// @Param        genres        formData  string   false  "Genres [IDs, comma separated]"
// @Param        casts         formData  string   false  "Casts [IDs, comma separated]"
// @Param        schedules     formData  string  false  "Schedules (JSON array: [{}])"
// @Param        poster        formData  file    false  "Poster file"
// @Param        backdrop      formData  file    false  "Backdrop file"
// @Success      200           {object}  models.SuccessResponse
// @Failure      400           {object}  models.ErrorResponse
// @Failure      401  		   {object}  models.ErrorResponse
// @Failure      500           {object}  models.ErrorResponse
// @Router       /admin/movies/edit/{id} [patch]
func (h *AdminHandler) UpdateMovies(ctx *gin.Context) {
	MovieID, _ := strconv.Atoi(ctx.Param("id"))

	var update models.EditMovies
	if err := ctx.ShouldBind(&update); err != nil {

		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	posterPath, err := utils.UploadFile(ctx, "poster", "public/movies/posters", "poster")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to upload poster"})
		return
	}
	if posterPath != "" {
		update.PosterPath = &posterPath
	}

	backdropPath, err := utils.UploadFile(ctx, "backdrop", "public/movies/backdrops", "backdrop")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to upload backdrop",
		})
		return
	}
	if backdropPath != "" {
		update.BackdropPath = &backdropPath
	}

	if genresStr := ctx.PostForm("genres"); genresStr != "" {
		var genres []int
		if err := json.Unmarshal([]byte(genresStr), &genres); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "invalid genres format, must be JSON array like [1,2,3]",
			})
			return
		}
		update.Genres = &genres
	}

	if castsStr := ctx.PostForm("casts"); castsStr != "" {
		var casts []int
		if err := json.Unmarshal([]byte(castsStr), &casts); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "invalid casts format, must be JSON array like [1,2,3]",
			})
			return
		}
		update.Casts = &casts
	}

	schedulesStr := ctx.PostForm("schedules")
	if schedulesStr != "" {
		var schedules []models.Schedule
		if err := json.Unmarshal([]byte(schedulesStr), &schedules); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "invalid schedules format",
			})
			return
		}
		update.Schedules = &schedules
	}

	if err := h.repo.UpdateMovies(ctx, MovieID, update); err != nil {
		log.Printf("%s", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "movie updated successfully",
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
// @Failure      401  {object}  models.ErrorResponse
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
// @Security     BearerAuth
// @Param        id   path      int  true  "Movie ID"
// @Success      200  {object}  models.SuccessResponse
// @Failure      400  {object}  models.ErrorResponse
// @Failure      401  {object}  models.ErrorResponse
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

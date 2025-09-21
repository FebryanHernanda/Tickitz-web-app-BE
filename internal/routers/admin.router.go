package routers

import (
	"github.com/FebryanHernanda/Tickitz-web-app-BE/internal/handlers"
	"github.com/FebryanHernanda/Tickitz-web-app-BE/internal/middlewares"
	"github.com/FebryanHernanda/Tickitz-web-app-BE/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func AdminRouter(r *gin.Engine, adminHandler *handlers.AdminHandler, jwtManager *utils.JWTManager, rdb *redis.Client) {
	adminRoutes := r.Group("/admin")
	adminRoutes.Use(middlewares.VerifyToken(jwtManager, rdb))
	adminRoutes.Use(middlewares.AuthMiddleware("admin"))

	adminRoutes.GET("/movies", adminHandler.GetAllMovies)
	adminRoutes.POST("/movies/add", adminHandler.AddMovies)
	adminRoutes.GET("/movies/schedule", adminHandler.GetMovieSchedule)
	adminRoutes.POST("/movies/cinemaschedule/add", adminHandler.AddCinemaSchedule)
	adminRoutes.DELETE("/movies/delete/:id", adminHandler.DeleteMovies)
	adminRoutes.PATCH("/movies/edit/:id", adminHandler.UpdateMovies)
	adminRoutes.GET("/movies/:movieEditId/edit-details", adminHandler.GetMovieEditDetail)
}

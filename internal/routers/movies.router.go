package routers

import (
	"github.com/FebryanHernanda/Tickitz-web-app-BE/internal/handlers"
	"github.com/gin-gonic/gin"
)

func MoviesRouter(r *gin.Engine, moviesHandler *handlers.MoviesHandler) {
	moviesRoutes := r.Group("/movies")
	moviesRoutes.GET("/upcoming", moviesHandler.GetUpcomingMovies)
	moviesRoutes.GET("/popular", moviesHandler.GetPopularMovies)
	moviesRoutes.GET("/:id/details", moviesHandler.GetDetailMovies)
	moviesRoutes.GET("", moviesHandler.GetMoviesByFilter)
	moviesRoutes.GET("/schedules", moviesHandler.GetSchedulesMovies)
}

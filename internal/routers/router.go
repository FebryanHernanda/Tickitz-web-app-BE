package routers

import (
	"log"
	"os"

	"github.com/FebryanHernanda/Tickitz-web-app-BE/docs"
	"github.com/FebryanHernanda/Tickitz-web-app-BE/internal/handlers"
	"github.com/FebryanHernanda/Tickitz-web-app-BE/internal/middlewares"
	"github.com/FebryanHernanda/Tickitz-web-app-BE/internal/repositories"
	"github.com/FebryanHernanda/Tickitz-web-app-BE/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func MainRouter(db *pgxpool.Pool) *gin.Engine {
	r := gin.Default()
	r.Use(middlewares.CORSmiddleware)

	// JWT
	jwtSecret := os.Getenv("JWTKEY")
	if jwtSecret == "" {
		log.Fatal("JWT Key env variable not set")
	}
	jwtManager := utils.NewJWTManager(jwtSecret)

	// Movies repo & handlers
	movieRepo := repositories.NewMovieRepository(db)
	movieHandler := handlers.NewMoviesHandler(movieRepo)
	// Profile repo & handlers
	profileRepo := repositories.NewProfileRepository(db)
	profileHandler := handlers.NewProfileHandler(profileRepo)
	// Orders repo & handlers
	ordersRepo := repositories.NewOrdersRepository(db)
	ordersHandler := handlers.NewOrdersHandler(ordersRepo)
	// Admin repo & handlers
	adminRepo := repositories.NewAdminRepository(db)
	adminHandler := handlers.NewAdminHandler(adminRepo)
	// auth repo & handlers
	authRepo := repositories.NewUserRepository(db)
	authHandler := handlers.NewAuthHandler(authRepo, jwtManager)
	// seat repo & handlers
	seatRepo := repositories.NewSeatRepository(db)
	seatHandler := handlers.NewSeatHandler(seatRepo)

	// Register router
	MoviesRouter(r, movieHandler)
	ProfileRouter(r, profileHandler, jwtManager)
	OrdersRouter(r, ordersHandler, jwtManager)
	AdminRouter(r, adminHandler, jwtManager)
	AuthRouter(r, authHandler)
	SeatsRouter(r, seatHandler)

	// Register Swagger
	docs.SwaggerInfo.BasePath = "/"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return r
}

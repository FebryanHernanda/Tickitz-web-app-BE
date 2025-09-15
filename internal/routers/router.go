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
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func MainRouter(db *pgxpool.Pool, rdb *redis.Client) *gin.Engine {
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
	movieHandler := handlers.NewMoviesHandler(movieRepo, rdb)
	// Profile repo & handlers
	profileRepo := repositories.NewProfileRepository(db)
	profileHandler := handlers.NewProfileHandler(profileRepo, rdb)
	// Orders repo & handlers
	ordersRepo := repositories.NewOrdersRepository(db)
	ordersHandler := handlers.NewOrdersHandler(ordersRepo, rdb)
	// Admin repo & handlers
	adminRepo := repositories.NewAdminRepository(db)
	adminHandler := handlers.NewAdminHandler(adminRepo, rdb)
	// auth repo & handlers
	authRepo := repositories.NewUserRepository(db)
	authHandler := handlers.NewAuthHandler(authRepo, jwtManager, rdb)
	// seat repo & handlers
	cinemaRepo := repositories.NewCinemaRepository(db)
	cinemaHandler := handlers.NewCinemaHandler(cinemaRepo, rdb)

	// Register router
	MoviesRouter(r, movieHandler)
	ProfileRouter(r, profileHandler, jwtManager, rdb)
	OrdersRouter(r, ordersHandler, jwtManager, rdb)
	AdminRouter(r, adminHandler, jwtManager, rdb)
	AuthRouter(r, jwtManager, rdb, authHandler)
	CinemaRouter(r, cinemaHandler)

	// register file upload
	r.Static("/public", "./public")

	// Register Swagger
	docs.SwaggerInfo.BasePath = "/"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return r
}

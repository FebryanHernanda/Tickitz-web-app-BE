package routers

import (
	"github.com/FebryanHernanda/Tickitz-web-app-BE/internal/handlers"
	"github.com/FebryanHernanda/Tickitz-web-app-BE/internal/middlewares"
	"github.com/FebryanHernanda/Tickitz-web-app-BE/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func OrdersRouter(r *gin.Engine, ordersHandler *handlers.OrdersHandler, jwtManager *utils.JWTManager, rdb *redis.Client) {
	ordersRoutes := r.Group("/orders")
	ordersRoutes.Use(middlewares.VerifyToken(jwtManager, rdb))
	ordersRoutes.Use(middlewares.AuthMiddleware("user"))
	ordersRoutes.POST("/", ordersHandler.CreateOrder)
	ordersRoutes.GET("/history", ordersHandler.GetOrdersHistory)
}

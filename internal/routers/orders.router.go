package routers

import (
	"github.com/FebryanHernanda/Tickitz-web-app-BE/internal/handlers"
	"github.com/FebryanHernanda/Tickitz-web-app-BE/internal/middlewares"
	"github.com/FebryanHernanda/Tickitz-web-app-BE/internal/utils"
	"github.com/gin-gonic/gin"
)

func OrdersRouter(r *gin.Engine, ordersHandler *handlers.OrdersHandler, jwtManager *utils.JWTManager) {
	ordersRoutes := r.Group("/orders")
	ordersRoutes.Use(middlewares.VerifyToken(jwtManager))
	ordersRoutes.Use(middlewares.AuthMiddleware("user"))
	ordersRoutes.POST("/", ordersHandler.CreateOrder)
	ordersRoutes.GET("/history", ordersHandler.GetOrdersHistory)
}

package routers

import (
	"github.com/FebryanHernanda/Tickitz-web-app-BE/internal/handlers"
	"github.com/gin-gonic/gin"
)

func OrdersRouter(r *gin.Engine, ordersHandler *handlers.OrdersHandler) {
	ordersRoutes := r.Group("/orders")
	ordersRoutes.POST("/", ordersHandler.CreateOrder)
	ordersRoutes.GET("/history/:userID", ordersHandler.GetOrdersHistory)
}

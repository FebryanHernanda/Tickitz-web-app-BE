package handlers

import (
	"net/http"
	"strconv"

	"github.com/FebryanHernanda/Tickitz-web-app-BE/internal/models"
	"github.com/FebryanHernanda/Tickitz-web-app-BE/internal/repositories"
	"github.com/FebryanHernanda/Tickitz-web-app-BE/internal/utils"
	"github.com/gin-gonic/gin"
)

type OrdersHandler struct {
	repo *repositories.OrdersRepository
}

func NewOrdersHandler(repo *repositories.OrdersRepository) *OrdersHandler {
	return &OrdersHandler{
		repo: repo,
	}
}

// CreateOrder godoc
// @Summary Create Order
// @Description  Create a new order
// @Tags Orders
// @Accept       json
// @Produce      json
// @Param order body models.OrderRequest true "Order data"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 409 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /orders/ [post]
func (h *OrdersHandler) CreateOrder(ctx *gin.Context) {
	var req models.OrderRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	QRCode := utils.GenerateQRCode()

	order := models.Order{
		QRCode:            QRCode,
		IsPaid:            req.IsPaid,
		IsActive:          req.IsActive,
		TotalPrices:       req.TotalPrices,
		UserID:            req.UserID,
		CinemasScheduleID: req.CinemasScheduleID,
		PaymentMethodID:   req.PaymentMethodID,
		OrderSeats:        req.OrderSeats,
	}

	seatIDs := make([]int, len(order.OrderSeats))
	for i, seat := range order.OrderSeats {
		seatIDs[i] = seat.SeatID
	}

	available, err := h.repo.AreSeatsAvailable(ctx, order.CinemasScheduleID, seatIDs)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	if !available {
		ctx.JSON(http.StatusConflict, gin.H{
			"success": false,
			"error":   "One or more seats are already booked",
		})
		return
	}

	orderID, err := h.repo.CreateOrder(ctx, &order)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	err = h.repo.CreateOrderSeats(ctx, orderID, order.OrderSeats)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	order.ID = orderID
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Order created successfully",
		"order":   order,
	})
}

// GetOrdersHistory godoc
// @Summary      Get order history by user ID
// @Description  Retrieve a user's order history
// @Tags         Orders
// @Produce      json
// @Param        userID path int true "User ID"
// @Success      200    {object} models.SuccessResponse
// @Failure 	 400    {object} models.ErrorResponse
// @Failure 	 404    {object} models.ErrorResponse
// @Failure 	 500    {object} models.ErrorResponse
// @Router       /orders/history/{userID} [get]
func (h *OrdersHandler) GetOrdersHistory(ctx *gin.Context) {
	userIDParams := ctx.Param("userID")
	userID, err := strconv.Atoi(userIDParams)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	orderHistory, err := h.repo.GetOrdersHistory(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to get order history",
		})
		return
	}

	if len(orderHistory) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Order history not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    orderHistory,
	})
}

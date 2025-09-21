package middlewares

import (
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
)

func CORSmiddleware(ctx *gin.Context) {
	whitelist := []string{"http://127.0.0.1:3000", "http://localhost:3000", "http://localhost:5173", "http://127.0.0.1:5173"}
	// get Header origin
	origin := ctx.GetHeader("Origin")
	// check request same with origin or not
	if slices.Contains(whitelist, origin) {
		ctx.Header("Access-Control-Allow-Origin", origin)
	}
	// give the access of method http
	ctx.Header("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
	// give the access for header
	ctx.Header("Access-Control-Allow-Headers", "Authorization, Content-Type")

	// check the preflight, if not same stop it
	if ctx.Request.Method == http.MethodOptions {
		ctx.AbortWithStatus(http.StatusNoContent)
		return
	}
	ctx.Next()
}

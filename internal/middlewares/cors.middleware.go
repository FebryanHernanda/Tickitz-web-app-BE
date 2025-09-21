package middlewares

import (
	"net/http"
	"os"
	"slices"

	"github.com/gin-gonic/gin"
)

func CORSmiddleware(ctx *gin.Context) {
	whitelistEnv := os.Getenv("FRONTEND_URL")
	whitelist := []string{whitelistEnv}
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

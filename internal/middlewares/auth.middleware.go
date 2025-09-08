package middlewares

import (
	"net/http"
	"slices"
	"strings"

	"github.com/FebryanHernanda/Tickitz-web-app-BE/internal/utils"
	"github.com/gin-gonic/gin"
)

// get token from header and verify that
func VerifyToken(jwtManager *utils.JWTManager) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		parts := strings.Fields(authHeader)

		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" || parts[1] == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Please login",
			})
			return
		}

		claims, err := jwtManager.ValidateToken(parts[1])
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Invalid or expired token",
			})
			return
		}

		ctx.Set("claims", claims)
		ctx.Next()
	}
}

// checks the user role is included in the allowed roles.
func AuthMiddleware(allowedRoles ...string) func(*gin.Context) {
	return func(ctx *gin.Context) {
		rawClaims, exist := ctx.Get("claims")
		if !exist {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"success": false,
				"error":   "Please login again",
			})
			return
		}

		claims, ok := rawClaims.(*utils.Claims)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "Internal server error",
			})
			return
		}

		if !slices.Contains(allowedRoles, claims.Role) {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"success": false,
				"error":   "You do not have access rights to this resource",
			})
			return
		}
		ctx.Next()
	}
}

package middleware

import (
	"gateway_mrc/helpers/csrf"
	"github.com/gin-gonic/gin"
	"net/http"
)

var (
	csrfHeaderName = "X-CSRF-TOKEN"
)

// TODO: надо доработать этот шедевр
func CSRFMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		_, err := ctx.Cookie("csrf_token")
		//fmt.Println("csrf", csrf)
		if err != nil && ctx.Request.Method == http.MethodGet {
			// Generate a new CSRF token
			token, err := csrf.GenerateCSRFToken()
			if err != nil {
				ctx.AbortWithStatus(http.StatusInternalServerError)
				return
			}
			csrf.SetCSRFCookie(ctx, token)
		} else if ctx.Request.Method == http.MethodPost {
			token, err := csrf.GetCSRFToken(ctx)
			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "CSRF token not found"})
				return
			}
			requestToken := ctx.GetHeader(csrfHeaderName)
			if token != requestToken {
				ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Invalid CSRF token"})
				return
			}
		}
		ctx.Next()

	}
}

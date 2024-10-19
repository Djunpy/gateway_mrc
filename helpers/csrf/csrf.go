package csrf

import (
	"crypto/rand"

	"encoding/base64"
	"github.com/gin-gonic/gin"
)

const (
	csrfTokenLength = 32
	csrfTokenName   = "csrf_token"
)

func GenerateCSRFToken() (string, error) {
	bytes := make([]byte, csrfTokenLength)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

func SetCSRFCookie(ctx *gin.Context, token string) {
	ctx.SetCookie(csrfTokenName, token, 3600, "/", "", false, true)
}

func GetCSRFToken(ctx *gin.Context) (string, error) {
	token, err := ctx.Cookie(csrfTokenName)
	if err != nil {
		return "", err
	}
	return token, nil
}

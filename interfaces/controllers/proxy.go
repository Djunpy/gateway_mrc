package controllers

import (
	"encoding/json"
	"fmt"
	"gateway_mrc/config"
	"gateway_mrc/entities"
	"gateway_mrc/helpers/server"
	"gateway_mrc/infrastructure/jwt_token"
	"gateway_mrc/usecase"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

type ProxyController struct {
	config         config.Config
	jwtMaker       jwt_token.Maker
	sessionUsecase usecase.SessionUsecase
}

func NewProxyController(jwtMaker jwt_token.Maker, sessionUsecase usecase.SessionUsecase, config config.Config) ProxyController {
	return ProxyController{sessionUsecase: sessionUsecase, jwtMaker: jwtMaker, config: config}
}

func (c *ProxyController) ProxySignInReq(ctx *gin.Context, target string) {
	var payload entities.SignInResponse
	sessionIdStr, exists := ctx.Get("sessionId")
	if !exists {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, server.Response(nil, server.GET_COOKIE_ERR_CODE, nil))
		return
	}

	url := server.GetReqFullUrl(ctx, target)
	fmt.Print(ctx.Request.Body)
	resp, err := server.CreateAndSendRequest(
		ctx.Request.Method,
		url,
		ctx.Request.Body,
		ctx.Request.Header,
	)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, server.Response(err, server.CREATING_REQUEST_ERR_CODE, nil))
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
		return
	}

	if resp.StatusCode == http.StatusOK {
		err = json.Unmarshal(body, &payload)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, server.Response(err, server.PARSING_RESPONSE_ERR_CODE, nil))
			return
		}
		refreshToken := payload.Body.RefreshToken
		accessToken := payload.Body.AccessToken
		_, statusCode, err := c.sessionUsecase.UpdateSession(ctx, sessionIdStr.(string), refreshToken, accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, server.Response(err, statusCode, nil))
			return
		}
	}
	if resp.StatusCode == http.StatusOK {
		ctx.JSON(resp.StatusCode, server.Response(nil, server.SUCCESS_CODE, nil))
	} else {
		ctx.Data(resp.StatusCode, resp.Header.Get("Content-Type"), body)
	}
}

func (c *ProxyController) ProxyCommonReq(ctx *gin.Context, target string) {
	sessionIdStr, exists := ctx.Get("sessionId")
	if !exists {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, server.Response(nil, server.GET_COOKIE_ERR_CODE, nil))
		return
	}

	session, errCode, err := c.sessionUsecase.GetSession(ctx, sessionIdStr.(string))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, server.Response(err, errCode, nil))
		return
	}
	url := server.GetReqFullUrl(ctx, target)
	headers := ctx.Request.Header.Clone()
	headers.Set("Authorization", fmt.Sprintf("Bearer %s", session.AccessToken.String))

	resp, err := server.CreateAndSendRequest(
		ctx.Request.Method,
		url,
		ctx.Request.Body,
		headers,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, server.Response(err, server.CREATING_REQUEST_ERR_CODE, nil))
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, server.Response(err, server.CREATING_REQUEST_ERR_CODE, nil))
		return
	}

	ctx.Data(resp.StatusCode, resp.Header.Get("Content-Type"), body)
}

func (c *ProxyController) ProxyLogoutReq(ctx *gin.Context) {
	ctx.SetCookie("session_id", "", -1, "/", "localhost", false, true)
	ctx.Status(http.StatusOK)
}

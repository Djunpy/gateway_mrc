package jwt_token

import (
	"errors"
	"fmt"
	"gateway_mrc/entities"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"time"
)

var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
)

type Payload struct {
	TokenType string    `json:"token_type"`
	ID        uuid.UUID `json:"id"`
	UserId    int64     `json:"user_id"`
	Email     string    `json:"email"`
	Groups    []string  `json:"roles"`
	UserType  string    `json:"user_type"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

func NewPayload(user entities.UserResponse, tokenType string, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	payload := &Payload{
		TokenType: tokenType,
		ID:        tokenID,
		UserId:    user.UserId,
		Email:     user.Email,
		Groups:    user.Groups,
		UserType:  user.UserType,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}
	return payload, nil
}

func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}

func GetJWTPayload(ctx *gin.Context) (*Payload, bool) {
	var jwtPayload *Payload
	ctxPayload, exists := ctx.Get("jwtTokenPayload")
	if !exists {
		return jwtPayload, false
	}

	jwtPayload, ok := ctxPayload.(*Payload)
	fmt.Println("jwtPayload", jwtPayload)
	if !ok {
		return jwtPayload, false
	}

	return jwtPayload, true
}

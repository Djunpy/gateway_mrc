package jwt_token

import "gateway_mrc/entities"

// Maker is an interface for managing tokens
type Maker interface {
	CreateToken(user entities.UserResponse, tokenType string) (token string, payload *Payload, err error)
	VerifyToken(token string) (*Payload, error)
	GetErrorCode(err error) (errorCode int32)
}

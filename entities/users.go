package entities

type UserResponse struct {
	UserId   int64    `json:"user_id"`
	Email    string   `json:"email"`
	Groups   []string `json:"roles"`
	UserType string   `json:"user_type"`
}

type SignInBodyResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type SignInResponse struct {
	CommonResponse
	Body SignInBodyResponse `json:"body"`
}

type RefreshTokenBodyResponse struct {
	AccessToken string `json:"access_token"`
}

type RefreshTokenResponse struct {
	CommonResponse
	RefreshTokenBodyResponse `json:"body"`
}

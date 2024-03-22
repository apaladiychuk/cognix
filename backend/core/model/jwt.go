package model

type LoginResponse struct {
	AccessToken  string
	RefreshToken string
	User         *User
}

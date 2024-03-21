package model

type JwtClaim struct {
	AccessToken  string
	RefreshToken string
	UserName     string
	TenantID     string
}

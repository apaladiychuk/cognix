package connector

import (
	"golang.org/x/oauth2"
	"time"
)

func TestGoogeDriveTest() {
	token := &oauth2.Token{
		AccessToken:  "",
		TokenType:    "",
		RefreshToken: "",
		Expiry:       time.Time{},
	}
	GetDriver(token)
}

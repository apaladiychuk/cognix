package connector

import (
	"golang.org/x/oauth2"
	"testing"
	"time"
)

func TestGoogeDrive_Test(t *testing.T) {
	token := &oauth2.Token{
		AccessToken:  "ya29.a0AXooCgsLHpznp2JjdMkkzMP31REsg-k7je01I2hxwZW-M63K-bNTKw7zysv9oYjL7RsY4_H0Lj9MQpFwt-V9CYEtyyIEiebivOk0jDv3eGm3dvDcs7_q2ky_1phwDVLx9yya_ti5rUAM3zngoWXxZPvbfaqjWudQ3mnvaCgYKAaYSARMSFQHGX2MiQnjLAaECKySfyAp5v5yHEg0171",
		TokenType:    "Bearer",
		RefreshToken: "1//09kIdjeE8HWsHCgYIARAAGAkSNwF-L9IrklqUB2BgQ37PEq_G_CsdaA9ZgRFfDUyEhb89tE7ywvUayG0UgLiP5SjBiHpS40xjiFo",
		Expiry:       time.Now(),
	}
	GetDriver(token)
}

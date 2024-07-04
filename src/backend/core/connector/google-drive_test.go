package connector

import (
	"golang.org/x/oauth2"
	"testing"
	"time"
)

func TestGoogeDrive_Test(t *testing.T) {
	token := &oauth2.Token{
		AccessToken:  "ya29.a0AXooCgucn_I6e6NfvfljATdn9AxjUGH7SzOjN3d5k5iVQnKAJs8JwSLnhNub7QJS3qMBvgo5AEhHqHxSO668UsM1fh4zKd51wVMka7w50C8-Gj-dpX6OUAZiLcwInfGFhtg1Mt8Oiw-A292UvIT1poyG7d4878nTLPeMaCgYKAWESARMSFQHGX2MicyDe4oMBrqR9_bTIuh-yPQ0171",
		TokenType:    "Bearer",
		RefreshToken: "1//09yOzMP7QLgd1CgYIARAAGAkSNwF-L9IrXM_6MjfiXJ1eJLhW3vW0bDrq0Dhx8YtycrlV63r9CRkBdP5Z6k_WoxdZ9dd4V5kvKcw",
		Expiry:       time.Now(),
	}
	GetDriver(token)
}

package pinningcli

import (
	"os"
	"strings"
)

var baseUrl = "https://api.w3ipfs.aioz.network/api"

func init() {
	if value := strings.TrimSpace(os.Getenv("PINNING_SERVICE_URL")); value != "" {
		SetBaseURL(value)
	}
}

func SetBaseURL(value string) {
	value = strings.TrimRight(strings.TrimSpace(value), "/")
	if value != "" {
		baseUrl = value
	}
}

const ErrNotOverwritten = "aborted. Existing credentials were not overwritten"

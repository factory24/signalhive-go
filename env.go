package connectivity

import (
	"fmt"
	"os"
	"strings"
)

// env reads one setting, e.g. env("CELLULAR.API_KEY") reads CONNECTIVITY.CELLULAR.API_KEY.
// The prefix names what the setting is for, not who sells it, so a rename never
// touches a consumer's configuration. The underscored spelling exists because
// some consumers cannot inject dots.
func env(key string) string {
	underscored := strings.ReplaceAll(key, ".", "_")
	for _, name := range []string{"CONNECTIVITY." + key, "CONNECTIVITY_" + underscored} {
		if v := os.Getenv(name); v != "" {
			return v
		}
	}
	return ""
}

// NewFromEnv builds a Client from environment configuration:
//
//	CONNECTIVITY.BASE_URL       API root, e.g. https://api.example.com/v1 (required)
//	CONNECTIVITY.API_TOKEN      static bearer token, or
//	CONNECTIVITY.TOKEN_URL      OpenID token endpoint for client credentials
//	CONNECTIVITY.CLIENT_ID      with
//	CONNECTIVITY.CLIENT_SECRET
//	CONNECTIVITY.BUSINESS_ID    only for staff principals acting for a business
//
// Underscored spellings (CONNECTIVITY_BASE_URL, ...) are accepted too.
func NewFromEnv() (*Client, error) {
	baseURL := env("BASE_URL")
	if baseURL == "" {
		return nil, fmt.Errorf("connectivity: CONNECTIVITY.BASE_URL is not set")
	}

	var tokens TokenSource
	if t := env("API_TOKEN"); t != "" {
		tokens = StaticToken(t)
	} else {
		tokenURL := env("TOKEN_URL")
		clientID := env("CLIENT_ID")
		clientSecret := env("CLIENT_SECRET")
		if tokenURL == "" || clientID == "" {
			return nil, fmt.Errorf("connectivity: set CONNECTIVITY.API_TOKEN, or CONNECTIVITY.TOKEN_URL and CONNECTIVITY.CLIENT_ID/CLIENT_SECRET")
		}
		tokens = NewClientCredentialsTokenSource(tokenURL, clientID, clientSecret)
	}

	var opts []Option
	if biz := env("BUSINESS_ID"); biz != "" {
		opts = append(opts, WithBusinessID(biz))
	}
	return New(baseURL, tokens, opts...), nil
}

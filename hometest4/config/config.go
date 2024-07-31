package config

import (
	"os"
	"time"
)

const (
	JwtContextKey       string        = "userToken"
	UserContextKey      string        = "user"
	DefaultJwtSecret    string        = "it's secret"
	TokenUsernameField  string        = "username"
	TokenExpiration     time.Duration = time.Hour * 2
	TopTransactionLimit int           = 10
)

func GetJwtSecret() string {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = DefaultJwtSecret
	}
	return jwtSecret
}

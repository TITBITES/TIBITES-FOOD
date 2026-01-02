package services

import "context"

type AuthService interface {
	Register(ctx context.Context, name, email, phone, password string) (AuthTokens, error)
	Login(ctx context.Context, email, password string) (AuthTokens, error)
	Refresh(ctx context.Context, refreshToken string) (AuthTokens, error)
	AdminOAuthStart(ctx context.Context, state string) (authorizationURL string, err error)
	AdminOAuthCallback(ctx context.Context, code, state string) (AuthTokens, error)
}

type AuthTokens struct {
	AccessToken  string
	RefreshToken string
}

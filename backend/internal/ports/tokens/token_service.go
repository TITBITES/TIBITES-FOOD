package tokens

import "context"

type TokenService interface {
	GenerateAccessToken(ctx context.Context, userID, role string) (string, error)
	GenerateRefreshToken(ctx context.Context, userID string) (string, error)
	ValidateAccessToken(ctx context.Context, token string) (userID string, role string, err error)
	ValidateRefreshToken(ctx context.Context, token string) (userID string, err error)
}

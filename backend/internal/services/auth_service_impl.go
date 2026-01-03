package services

import (
	"context"
	"errors"
	"strings"

	portscrypto "local.dev/foodapp/internal/ports/crypto"
	"local.dev/foodapp/internal/ports/repositories"
	portstokens "local.dev/foodapp/internal/ports/tokens"
)

type authService struct{
	repo repositories.UserRepository
	hasher portscrypto.PasswordHasher
	tokens portstokens.TokenService
}

func NewAuthService(repo repositories.UserRepository, hasher portscrypto.PasswordHasher, tokens portstokens.TokenService) AuthService {
	return &authService{repo: repo, hasher: hasher, tokens: tokens}
}

var (
	errDuplicateEmail = errors.New("duplicate_email")
	errInvalidCredentials = errors.New("invalid_credentials")
)

func (a *authService) Register(ctx context.Context, name, email, phone, password string) (AuthTokens, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	exists, err := a.repo.ExistsByEmail(ctx, email)
	if err != nil { return AuthTokens{}, err }
	if exists { return AuthTokens{}, errDuplicateEmail }
	hash, err := a.hasher.Hash(password)
	if err != nil { return AuthTokens{}, err }
	rec := repositories.UserRecord{Email: email, PasswordHash: hash, Role: "Customer"}
	rec, err = a.repo.Create(ctx, rec)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") { return AuthTokens{}, errDuplicateEmail }
		return AuthTokens{}, err
	}
	access, err := a.tokens.GenerateAccessToken(ctx, rec.ID, rec.Role)
	if err != nil { return AuthTokens{}, err }
	refresh, err := a.tokens.GenerateRefreshToken(ctx, rec.ID)
	if err != nil { return AuthTokens{}, err }
	return AuthTokens{AccessToken: access, RefreshToken: refresh}, nil
}

func (a *authService) Login(ctx context.Context, email, password string) (AuthTokens, error) {
	rec, err := a.repo.GetByEmail(ctx, strings.TrimSpace(strings.ToLower(email)))
	if err != nil { return AuthTokens{}, errInvalidCredentials }
	if err := a.hasher.Compare(rec.PasswordHash, password); err != nil { return AuthTokens{}, errInvalidCredentials }
	access, err := a.tokens.GenerateAccessToken(ctx, rec.ID, rec.Role)
	if err != nil { return AuthTokens{}, err }
	refresh, err := a.tokens.GenerateRefreshToken(ctx, rec.ID)
	if err != nil { return AuthTokens{}, err }
	return AuthTokens{AccessToken: access, RefreshToken: refresh}, nil
}

func (a *authService) Refresh(ctx context.Context, refreshToken string) (AuthTokens, error) {
	userID, err := a.tokens.ValidateRefreshToken(ctx, refreshToken)
	if err != nil { return AuthTokens{}, errInvalidCredentials }
	// Issue new pair with default Customer role (no repo dependency here)
	access, err := a.tokens.GenerateAccessToken(ctx, userID, "Customer")
	if err != nil { return AuthTokens{}, err }
	refresh, err := a.tokens.GenerateRefreshToken(ctx, userID)
	if err != nil { return AuthTokens{}, err }
	return AuthTokens{AccessToken: access, RefreshToken: refresh}, nil
}

func (a *authService) AdminOAuthStart(ctx context.Context, state string) (authorizationURL string, err error) {
	return "", errors.New("not_implemented")
}

func (a *authService) AdminOAuthCallback(ctx context.Context, code, state string) (AuthTokens, error) {
	return AuthTokens{}, errors.New("not_implemented")
}

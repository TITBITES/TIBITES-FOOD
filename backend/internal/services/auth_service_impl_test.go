package services

import (
	"context"
	"testing"
	"time"

	adaptercrypto "local.dev/foodapp/internal/adapters/crypto"
	adapterrepo "local.dev/foodapp/internal/adapters/repositories/memory"
	adaptertokens "local.dev/foodapp/internal/adapters/tokens"
)

func newAuth() (AuthService, *adapterrepo.UserRepoMemory) {
	repo := adapterrepo.NewUserRepoMemory()
	hasher := adaptercrypto.NewBcryptHasher(10)
	tokens := adaptertokens.NewJWTService("acc", "ref", time.Minute, 7*24*time.Hour)
	return NewAuthService(repo, hasher, tokens), repo
}

func TestRegisterSuccess(t *testing.T) {
	auth, _ := newAuth()
	ctx := context.Background()
	tok, err := auth.Register(ctx, "Alice", "alice@example.com", "+10000000000", "password")
	if err != nil { t.Fatalf("register error: %v", err) }
	if tok.AccessToken == "" || tok.RefreshToken == "" { t.Fatalf("expected tokens") }
}

func TestRegisterDuplicateEmail(t *testing.T) {
	auth, _ := newAuth()
	ctx := context.Background()
	_, _ = auth.Register(ctx, "Alice", "dup@example.com", "+1", "password")
	_, err := auth.Register(ctx, "Bob", "dup@example.com", "+2", "password")
	if err == nil { t.Fatalf("expected duplicate error") }
}

func TestLoginAndRefresh(t *testing.T) {
	auth, _ := newAuth()
	ctx := context.Background()
	_, _ = auth.Register(ctx, "Alice", "login@example.com", "+1", "password")
	tok, err := auth.Login(ctx, "login@example.com", "password")
	if err != nil { t.Fatalf("login error: %v", err) }
	if tok.AccessToken == "" || tok.RefreshToken == "" { t.Fatalf("expected tokens") }
	newTok, err := auth.Refresh(ctx, tok.RefreshToken)
	if err != nil { t.Fatalf("refresh error: %v", err) }
	if newTok.AccessToken == "" || newTok.RefreshToken == "" { t.Fatalf("expected new tokens") }
}

func TestLoginInvalidPassword(t *testing.T) {
	auth, _ := newAuth()
	ctx := context.Background()
	_, _ = auth.Register(ctx, "Alice", "wrongpw@example.com", "+1", "password")
	_, err := auth.Login(ctx, "wrongpw@example.com", "bad")
	if err == nil { t.Fatalf("expected invalid credentials error") }
}

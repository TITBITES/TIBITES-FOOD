package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	adaptercrypto "local.dev/foodapp/internal/adapters/crypto"
	adaptertokens "local.dev/foodapp/internal/adapters/tokens"
	stub "local.dev/foodapp/internal/adapters/payments/stub"
	"local.dev/foodapp/internal/adapters/postgres"
	"local.dev/foodapp/internal/httpapi"
	"local.dev/foodapp/internal/services"
)

func main() {
	// Config from environment (with defaults for local dev)
	accessSecret := getenv("JWT_ACCESS_SECRET", "dev-access-secret")
	refreshSecret := getenv("JWT_REFRESH_SECRET", "dev-refresh-secret")
	bcryptCost := 12

	// Adapters
	hasher := adaptercrypto.NewBcryptHasher(bcryptCost)
	tokenSvc := adaptertokens.NewJWTService(accessSecret, refreshSecret, 15*time.Minute, 7*24*time.Hour)

	// PostgreSQL persistence
	pgURL := postgres.MustEnvDatabaseURL()
	pgdb, err := postgres.Connect(context.Background(), pgURL)
	if err != nil { log.Fatalf("db connect error: %v", err) }
	defer pgdb.Close()
	if err := postgres.EnsureExtensions(context.Background(), pgdb); err != nil { log.Fatalf("db ensure ext: %v", err) }

	userRepo := postgres.NewUserRepositoryPG(pgdb)
	catRepo := postgres.NewCategoryRepositoryPG(pgdb)
	itemRepo := postgres.NewMenuItemRepositoryPG(pgdb)
	orderRepo := postgres.NewOrderRepositoryPG(pgdb)
	orderSecretRepo := postgres.NewOrderSecretRepositoryPG(pgdb)
	piRepo := postgres.NewPaymentIntentRepositoryPG(pgdb)
	provider := stub.NewProviderStub()
	_ = stub.NewWebhookVerifierStub()

	// Services
	_ = services.NewAuthService(userRepo, hasher, tokenSvc)
	_ = services.NewMenuService(catRepo, itemRepo)
	_ = services.NewOrderService(orderRepo, orderSecretRepo)
	_ = services.NewPaymentIntentService(orderRepo, orderSecretRepo, piRepo, provider)

	// HTTP
	r := httpapi.NewRouter(tokenSvc)
	addr := ":8080"
	log.Printf("starting http server on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" { return v }
	return def
}

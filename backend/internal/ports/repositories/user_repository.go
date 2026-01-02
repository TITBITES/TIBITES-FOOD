package repositories

import "context"

type UserRepository interface {
	GetByEmail(ctx context.Context, email string) (UserRecord, error)
	Create(ctx context.Context, rec UserRecord) (UserRecord, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}

type UserRecord struct{
	ID string
	Email string
	PasswordHash string
	Role string
}

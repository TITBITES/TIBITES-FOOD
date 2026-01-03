package crypto

import (
	"golang.org/x/crypto/bcrypt"
)

type BcryptHasher struct{
	Cost int
}

func NewBcryptHasher(cost int) *BcryptHasher { return &BcryptHasher{Cost: cost} }

func (b *BcryptHasher) Hash(plain string) (string, error) {
	h, err := bcrypt.GenerateFromPassword([]byte(plain), b.Cost)
	if err != nil { return "", err }
	return string(h), nil
}

func (b *BcryptHasher) Compare(hash, plain string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain))
}

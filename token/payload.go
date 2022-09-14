package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ERR_TOKEN_EXPIRED = errors.New("Token has expired")
	ERR_INVALID_TOKEN = errors.New("Invalid token")
)

type Payload struct {
	ID        uuid.UUID `json:"uid"`
	Name      string    `json:"name"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

func NewPayload(username string, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	payload := &Payload{
		ID:        tokenID,
		Name:      username,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}
	return payload, nil
}

func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return ERR_TOKEN_EXPIRED
	}
	return nil
}

package token

import (
	"fmt"
	"time"

	"github.com/o1egl/paseto"
	"golang.org/x/crypto/chacha20poly1305"
)

type PasetoMaker struct {
	paseto *paseto.V2
	secret []byte
}

func NewPasetoMaker(secret string) (Maker, error) {
	if len(secret) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size: required %d characters", chacha20poly1305.KeySize)
	}
	maker := &PasetoMaker{
		paseto: paseto.NewV2(),
		secret: []byte(secret),
	}
	return maker, nil
}

func (maker *PasetoMaker) CreateToken(userID int64, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(userID, duration)
	if err != nil {
		return "", nil, err
	}
	token, err := maker.paseto.Encrypt(maker.secret, payload, nil)
	if err != nil {
		return "", nil, err
	}
	return token, payload, nil
}

func (maker *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}
	err := maker.paseto.Decrypt(token, maker.secret, payload, nil)
	if err != nil {
		return nil, ERR_INVALID_TOKEN
	}

	err = payload.Valid()
	if err != nil {
		return nil, err
	}

	return payload, nil
}

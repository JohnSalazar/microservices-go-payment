package models

import (
	"crypto/rsa"
	"time"
)

type RSAPrivateKey struct {
	PrivateKey *rsa.PrivateKey `json:"privateKey"`
	Kid        string          `json:"kid"`
	ExpiresAt  time.Time       `json:"expires_at"`
}

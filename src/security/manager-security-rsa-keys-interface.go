package security

import (
	"crypto/rsa"
	"payment/src/models"

	common_models "github.com/oceano-dev/microservices-go-common/models"
)

type ManagerSecurityRSAKeys interface {
	GetAllRSAPublicKeys() []*common_models.RSAPublicKey
	GetAllRSAPrivateKeys() []*models.RSAPrivateKey
	GetNewestRSAPrivateKey() *models.RSAPrivateKey
	GetKeyFromKid(kid string) *rsa.PrivateKey
}

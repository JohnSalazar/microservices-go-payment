package security

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"payment/src/models"
	"time"

	"github.com/JohnSalazar/microservices-go-common/config"
	common_models "github.com/JohnSalazar/microservices-go-common/models"
	"github.com/google/uuid"
)

var (
	rsaKeys               []*models.RSAPrivateKey
	newestRSAPrivateKey   *models.RSAPrivateKey
	refreshRSAPrivateKeys = time.Now()
)

type managerSecurityRSAKeys struct {
	config *config.Config
}

func NewManagerSecurityRSAKeys(
	config *config.Config,
) *managerSecurityRSAKeys {
	return &managerSecurityRSAKeys{
		config: config,
	}
}

func (m *managerSecurityRSAKeys) GetAllRSAPublicKeys() []*common_models.RSAPublicKey {
	modelsRSAPrivateKeys := m.GetAllRSAPrivateKeys()

	var rsaPublicKeys []*common_models.RSAPublicKey

	for _, model := range modelsRSAPrivateKeys {
		modelRSAPublicKey := &common_models.RSAPublicKey{
			Key:       &model.PrivateKey.PublicKey,
			Kid:       model.Kid,
			ExpiresAt: model.ExpiresAt,
		}

		rsaPublicKeys = append(rsaPublicKeys, modelRSAPublicKey)
	}

	return rsaPublicKeys
}

func (m *managerSecurityRSAKeys) GetAllRSAPrivateKeys() []*models.RSAPrivateKey {
	if rsaKeys == nil {
		newestRSAPrivateKey = m.GetNewestRSAPrivateKey()
	}

	return rsaKeys
}

func (m *managerSecurityRSAKeys) GetNewestRSAPrivateKey() *models.RSAPrivateKey {
	var err error
	if newestRSAPrivateKey == nil {
		newestRSAPrivateKey = m.getNewestRSAPrivateKeys()
		m.refreshRSAPrivateKeys()
	}

	if newestRSAPrivateKey == nil {
		newestRSAPrivateKey, err = m.generateRSAPrivateKey()
		if err != nil {
			return nil
		}

		m.refreshRSAPrivateKeys()

		return newestRSAPrivateKey
	}

	rsaPrivateKeysRefresh := refreshRSAPrivateKeys.Before(time.Now().UTC())
	if rsaPrivateKeysRefresh {
		newestRSAPrivateKey = m.getNewestRSAPrivateKeys()

		m.refreshRSAPrivateKeys()
		fmt.Println("refresh rsa private keys")
	}

	rsaPrivateKeyExpires := newestRSAPrivateKey.ExpiresAt.Before(time.Now().UTC())
	if rsaPrivateKeyExpires {
		newestRSAPrivateKey, err = m.generateRSAPrivateKey()
		if err != nil {
			return nil
		}

		m.refreshRSAPrivateKeys()
	}

	return newestRSAPrivateKey
}

func (m *managerSecurityRSAKeys) GetKeyFromKid(kid string) *rsa.PrivateKey {
	for _, key := range rsaKeys {
		if key.Kid == kid {
			return key.PrivateKey
		}
	}

	return nil
}

func (m *managerSecurityRSAKeys) getNewestRSAPrivateKeys() *models.RSAPrivateKey {
	if rsaKeys == nil {
		return nil
	}

	return rsaKeys[0]
}

func (m *managerSecurityRSAKeys) generateRSAPrivateKey() (*models.RSAPrivateKey, error) {
	rsaPrivateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	modelRSAPrivateKey := &models.RSAPrivateKey{
		PrivateKey: rsaPrivateKey,
		Kid:        uuid.New().String(),
		ExpiresAt:  time.Now().UTC().Add(time.Duration(24*m.config.SecurityRSAKeys.DaysToExpireRSAKeys) * time.Hour),
	}

	if rsaKeys != nil {
		rsaKeys = append(rsaKeys, nil)
		copy(rsaKeys[1:], rsaKeys)
		rsaKeys[0] = modelRSAPrivateKey
	} else {
		rsaKeys = append(rsaKeys, modelRSAPrivateKey)
	}

	return modelRSAPrivateKey, nil
}

func (m *managerSecurityRSAKeys) refreshRSAPrivateKeys() {
	refreshRSAPrivateKeys = time.Now().UTC().Add(time.Minute * time.Duration(m.config.SecurityRSAKeys.MinutesToRefreshRSAPrivateKeys))
}

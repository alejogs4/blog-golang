package token

import (
	"crypto/rsa"
	"io/ioutil"
	"sync"

	"github.com/dgrijalva/jwt-go"
)

var (
	signKey            *rsa.PrivateKey
	verifyKey          *rsa.PublicKey
	certificatesLoader sync.Once
)

func LoadCertificates(privateKeyPath, publicKeyPath string) error {
	var error error
	certificatesLoader.Do(func() {
		error = loadCerticates(privateKeyPath, publicKeyPath)
	})

	return error
}

func loadCerticates(privateKeyPath, publicKeyPath string) error {
	privateKey, err := ioutil.ReadFile(privateKeyPath)
	if err != nil {
		return err
	}

	publicKey, err := ioutil.ReadFile(publicKeyPath)
	if err != nil {
		return err
	}

	return parseRSAKeys(privateKey, publicKey)
}

func parseRSAKeys(privateKey, publicKey []byte) error {
	var err error

	signKey, err = jwt.ParseRSAPrivateKeyFromPEM(privateKey)
	if err != nil {
		return err
	}

	verifyKey, err = jwt.ParseRSAPublicKeyFromPEM(publicKey)
	if err != nil {
		return err
	}

	return nil
}

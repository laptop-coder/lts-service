// see https://pkg.go.dev/crypto/ed25519
package utils

import (
	. "backend/config"
	"crypto/ed25519"
	"encoding/pem"
	"errors"
	"os"
)

func GenKeysIfNotExist() error {
	publicKeyExists, err := CheckFileExistence(&Cfg.ED25519.PathToPublicKey)
	if err != nil {
		return err
	}
	privateKeyExists, err := CheckFileExistence(&Cfg.ED25519.PathToPrivateKey)
	if err != nil {
		return err
	}

	if !*publicKeyExists || !*privateKeyExists {
		// nil means using the default crypto/rand.Reader
		publicKey, privateKey, err := ed25519.GenerateKey(nil)
		if err != nil {
			return err
		}

		publicPEM := pem.EncodeToMemory(&pem.Block{
			Type:  "ED25519 PUBLIC KEY",
			Bytes: publicKey,
		})
		privatePEM := pem.EncodeToMemory(&pem.Block{
			Type:  "ED25519 PRIVATE KEY",
			Bytes: privateKey,
		})

		if err := os.WriteFile(Cfg.ED25519.PathToPublicKey, publicPEM, 0400); err != nil {
			return err
		}
		if err := os.WriteFile(Cfg.ED25519.PathToPrivateKey, privatePEM, 0400); err != nil {
			return err
		}
	}

	return nil
}

func GetPublicKey() (ed25519.PublicKey, error) {
	publicPEM, err := os.ReadFile(Cfg.ED25519.PathToPublicKey)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(publicPEM)
	if block == nil {
		return nil, errors.New("error decoding public key PEM")
	}

	if block.Type != "ED25519 PUBLIC KEY" {
		return nil, errors.New("wrong public key type")
	}

	return block.Bytes, nil
}

func GetPrivateKey() (ed25519.PrivateKey, error) {
	privatePEM, err := os.ReadFile(Cfg.ED25519.PathToPrivateKey)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(privatePEM)
	if block == nil {
		return nil, errors.New("error decoding private key PEM")
	}

	if block.Type != "ED25519 PRIVATE KEY" {
		return nil, errors.New("wrong private key type")
	}

	return block.Bytes, nil
}

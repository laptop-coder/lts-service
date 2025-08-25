package utils

import (
	. "backend/config"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
)

func GenKeysIfNotExist() error {
	privateKeyExists, err := CheckFileExistence(&Cfg.RSA.PathToPrivateKey)
	if err != nil {
		return err
	}
	publicKeyExists, err := CheckFileExistence(&Cfg.RSA.PathToPublicKey)
	if err != nil {
		return err
	}

	if !*privateKeyExists && !*publicKeyExists {
		privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
		if err != nil {
			return err
		}
		privatePem := pem.EncodeToMemory(&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
		})
		publicPem := pem.EncodeToMemory(&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: x509.MarshalPKCS1PublicKey(&privateKey.PublicKey),
		})
		if err := os.WriteFile(Cfg.RSA.PathToPrivateKey, privatePem, 0444); err != nil {
			return err
		}
		if err := os.WriteFile(Cfg.RSA.PathToPublicKey, publicPem, 0444); err != nil {
			return err
		}
	}

	return nil
}

func GetPrivateKey() (*rsa.PrivateKey, error) {
	pemData, err := os.ReadFile(Cfg.RSA.PathToPrivateKey)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(pemData)
	if block == nil {
		return nil, errors.New("error decoding private key PEM")
	}

	if block.Type != "RSA PRIVATE KEY" {
		return nil, errors.New("wrong private key type")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

func GetPublicKey() (*rsa.PublicKey, error) {
	pemData, err := os.ReadFile(Cfg.RSA.PathToPublicKey)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(pemData)
	if block == nil {
		return nil, errors.New("error decoding public key PEM")
	}

	if block.Type != "RSA PUBLIC KEY" {
		return nil, errors.New("wrong public key type")
	}

	publicKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return publicKey, nil
}

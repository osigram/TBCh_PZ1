package keystorage

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
)

type KeyFile struct {
	PrivateKey string `json:"private_key"`
}

func MarshalECDSAPrivateKeyToJSON(privateKey *ecdsa.PrivateKey) ([]byte, error) {
	if privateKey == nil {
		return nil, errors.New("private key is nil")
	}

	privateKeyBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return nil, err
	}

	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privateKeyBytes,
	})

	keyJSON := KeyFile{
		PrivateKey: string(privateKeyPEM),
	}
	
	return json.Marshal(keyJSON)
}

func UnmarshalECDSAPrivateKeyFromJSON(jsonData []byte) (*ecdsa.PrivateKey, error) {
	var keyJSON KeyFile
	if err := json.Unmarshal(jsonData, &keyJSON); err != nil {
		return nil, err
	}

	block, _ := pem.Decode([]byte(keyJSON.PrivateKey))
	if block == nil {
		return nil, errors.New("failed to decode PEM block")
	}

	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	ecdsaKey, ok := privateKey.(*ecdsa.PrivateKey)
	if !ok {
		return nil, errors.New("not an ECDSA private key")
	}

	return ecdsaKey, nil
}

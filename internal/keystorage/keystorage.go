package keystorage

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"io"
	"os"
)

type KeyStorage struct {
	key *ecdsa.PrivateKey
}

func (k KeyStorage) Key() ecdsa.PrivateKey {
	return *k.key
}

func MustNewKeyStorage(filepath string) *KeyStorage {
	file, err := os.Open(filepath)
	if err != nil {
		return mustNewKeyStorage(filepath)
	}

	data, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}

	key, err := UnmarshalECDSAPrivateKeyFromJSON(data)
	if err != nil {
		panic(err)
	}

	return &KeyStorage{key: key}
}

func mustNewKeyStorage(filename string) *KeyStorage {
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}

	data, err := MarshalECDSAPrivateKeyToJSON(key)
	if err != nil {
		panic(err)
	}

	_, err = file.Write(data)
	if err != nil {
		panic(err)
	}

	return &KeyStorage{key: key}
}

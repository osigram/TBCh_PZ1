package domain

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
)

type Client ecdsa.PublicKey

type ECDSAPublicKeyJSON struct {
	Curve string `json:"curve"`
	X     string `json:"x"`
	Y     string `json:"y"`
}

func (pub Client) MarshalJSON() ([]byte, error) {
	xBytes := pub.X.Bytes()
	yBytes := pub.Y.Bytes()

	jsonKey := ECDSAPublicKeyJSON{
		Curve: pub.Curve.Params().Name,
		X:     base64.StdEncoding.EncodeToString(xBytes),
		Y:     base64.StdEncoding.EncodeToString(yBytes),
	}

	data, err := json.Marshal(jsonKey)
	if err != nil {
		return nil, err
	}

	newData := base64.URLEncoding.EncodeToString(data)

	return []byte("\"" + newData + "\""), nil
}

func (pub *Client) UnmarshalJSON(data []byte) error {
	data = data[1:]
	data = data[:len(data)-1]

	buff := make([]byte, base64.URLEncoding.DecodedLen(len(data)))
	n, err := base64.URLEncoding.Decode(buff, data)
	if err != nil {
		return err
	}
	buff = buff[:n]

	var jsonKey ECDSAPublicKeyJSON
	if err := json.Unmarshal(buff, &jsonKey); err != nil {
		return err
	}

	xBytes, err := base64.StdEncoding.DecodeString(jsonKey.X)
	if err != nil {
		return err
	}

	yBytes, err := base64.StdEncoding.DecodeString(jsonKey.Y)
	if err != nil {
		return err
	}

	switch jsonKey.Curve {
	case "P-224":
		pub.Curve = elliptic.P224()
	case "P-256":
		pub.Curve = elliptic.P256()
	case "P-384":
		pub.Curve = elliptic.P384()
	case "P-521":
		pub.Curve = elliptic.P521()
	default:
		return fmt.Errorf("unsupported curve: %s", jsonKey.Curve)
	}

	pub.X = new(big.Int).SetBytes(xBytes)
	pub.Y = new(big.Int).SetBytes(yBytes)

	return nil
}

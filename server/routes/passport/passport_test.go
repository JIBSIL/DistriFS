// Rather than having file-specific tests, this file contains a complete Passport test/example.

package passport

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"

	// "crypto/x509"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"testing"
)

// mostly auto-generated from Insomnia routes & responses cred: https://mholt.github.io/json-to-go/
type AuthenticateResponse struct {
	Data struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	} `json:"data"`
	Success bool `json:"success"`
}

type GetKeyRequest struct {
	Pubkey        string `json:"pubkey"`
	Signedmessage string `json:"signedmessage"`
	Messagekey    string `json:"messagekey"`
}

type GetKeyResponse struct {
	Data    string `json:"data"`
	Success bool   `json:"success"`
}

type VerifyKeyRequest struct {
	Key string `json:"key"`
}

type VerifyKeyResponse struct {
	Data    string `json:"data"`
	Success bool   `json:"success"`
}

func GetSigningKeys(passportSigningKey string) (string, []byte) {
	sigkey, _ := NewSigningKey()
	sig, _ := Sign([]byte(passportSigningKey), sigkey)
	msg := hex.EncodeToString(sig)
	fmt.Printf("Using %s hex-signed message\n", msg)
	// isVerified := Verify([]byte(passportSigningKey), sig, &sigkey.PublicKey)
	// fmt.Printf("SigVerify output: %s\n", strconv.FormatBool(isVerified))

	//privkey, _ := x509.MarshalECPrivateKey(sigkey)
	// fmt.Printf("Using privkey %s\n", hex.EncodeToString(privkey))

	// DECODING
	//derBytes, _ := hex.DecodeString(hexKey)
	//privateKey, _ := x509.ParseECPrivateKey(derBytes)

	publicKey := sigkey.Public().(*ecdsa.PublicKey)
	publicKeyBytes := elliptic.Marshal(publicKey.Curve, publicKey.X, publicKey.Y)
	// fmt.Printf("Using pubkey %s", hex.EncodeToString(publicKeyBytes))

	// DECODING
	//fpubKeyBytes, _ := hex.DecodeString(hexKey)
	//x, y := elliptic.Unmarshal(elliptic.P256(), pubKeyBytes)
	//publicKey := &ecdsa.PublicKey{Curve: elliptic.P256(), X: x, Y: y}

	return msg, publicKeyBytes
}

// NewSigningKey generates a random P-256 ECDSA private key.
func NewSigningKey() (*ecdsa.PrivateKey, error) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	return key, err
}

// Sign signs arbitrary data using ECDSA.
func Sign(data []byte, privkey *ecdsa.PrivateKey) ([]byte, error) {
	// hash message
	digest := sha256.Sum256(data)

	// sign the hash
	r, s, err := ecdsa.Sign(rand.Reader, privkey, digest[:])
	if err != nil {
		return nil, err
	}

	// encode the signature {R, S}
	// big.Int.Bytes() will need padding in the case of leading zero bytes
	params := privkey.Curve.Params()
	curveOrderByteSize := params.P.BitLen() / 8
	rBytes, sBytes := r.Bytes(), s.Bytes()
	signature := make([]byte, curveOrderByteSize*2)
	copy(signature[curveOrderByteSize-len(rBytes):], rBytes)
	copy(signature[curveOrderByteSize*2-len(sBytes):], sBytes)

	return signature, nil
}

func Verify(data, signature []byte, pubkey *ecdsa.PublicKey) bool {
	// hash message
	digest := sha256.Sum256(data)

	curveOrderByteSize := pubkey.Curve.Params().P.BitLen() / 8

	r, s := new(big.Int), new(big.Int)
	r.SetBytes(signature[:curveOrderByteSize])
	s.SetBytes(signature[curveOrderByteSize:])

	return ecdsa.Verify(pubkey, digest[:], r, s)
}

func SendGetRequest(t *testing.T, url string) []byte {
	res, err := http.Get(url)
	if err != nil {
		t.Errorf("Failed to get send request: %s", err)
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Failed to read response: %s", err)
	}

	return body
}

func SendPostRequest(t *testing.T, url string, reqBody interface{}) []byte {
	json_data, err := json.Marshal(reqBody)
	if err != nil {
		t.Errorf("Failed to marshal POST request: %s", err)
	}

	res, err := http.Post(url, "application/json", bytes.NewBuffer(json_data))
	if err != nil {
		t.Errorf("Failed to get send request: %s", err)
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Failed to read response: %s", err)
	}

	return body
}

func TestPassport(t *testing.T) {
	// Create a new Passport
	response := SendGetRequest(t, "http://localhost:8000/passport/authenticate")
	var result AuthenticateResponse
	if err := json.Unmarshal(response, &result); err != nil {
		t.Errorf("Failed to unmarshal JSON (req 1): %s", err)
	}

	msg, publicKeyBytes := GetSigningKeys(result.Data.Value)
	publicKey := hex.EncodeToString(publicKeyBytes)
	t.Logf("Using %s hex-signed message\n", msg)

	request := GetKeyRequest{
		Pubkey:        publicKey,
		Signedmessage: msg,
		Messagekey:    result.Data.Key,
	}
	// get the key from the server
	response = SendPostRequest(t, "http://localhost:8000/passport/getKey", request)
	var result2 GetKeyResponse
	if err := json.Unmarshal(response, &result2); err != nil {
		t.Errorf("Failed to unmarshal JSON (req 2): %s", err)
	}

	// verify key with server
	key := result2.Data
	request2 := VerifyKeyRequest{Key: key}
	response = SendPostRequest(t, "http://localhost:8000/passport/verify", request2)

	var result3 VerifyKeyResponse
	if err := json.Unmarshal(response, &result3); err != nil {
		t.Errorf("Failed to unmarshal JSON (req 3): %s", err)
	}

	t.Logf("Key returned by server is %s", result3.Data)

	verified := result3.Data == publicKey
	if !verified {
		t.Errorf("Key verification failed: public keys do not match (server thinks we're somebody else)!")
	} else {
		t.Logf("Public keys match!")
	}
}

package handoff

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io"
	"io/ioutil"

	"golang.org/x/crypto/nacl/box"
)

const (
	PEMType         = "HANDOFF MESSAGE"
	HeaderReference = "Reference"
	KeySize         = 32
	MaxPlaintext    = 4096
)

var b64Enc = base64.RawStdEncoding

type Keys struct {
	Public  string
	Private string
}

func Generate(reference string) (Keys, error) {
	var keys Keys
	pub, priv, err := box.GenerateKey(rand.Reader)
	if err != nil {
		return keys, err
	}

	keys.Public = b64Enc.EncodeToString(append(pub[:], []byte(reference)...))
	keys.Private = b64Enc.EncodeToString(priv[:])

	return keys, nil
}

func ParsePubKey(keyB64 string) ([KeySize]byte, string, error) {
	var key [KeySize]byte
	keyBytes, err := b64Enc.DecodeString(keyB64)
	if err != nil {
		return key, "", err
	}

	copy(key[:], keyBytes[:32])
	reference := string(keyBytes[32:])

	return key, reference, nil
}

func ParsePrivKey(keyB64 string) ([KeySize]byte, error) {
	var key [KeySize]byte
	keyBytes, err := b64Enc.DecodeString(keyB64)
	if err != nil {
		return key, err
	}

	copy(key[:], keyBytes[:32])
	return key, nil
}

func Encrypt(pub [KeySize]byte, reference string, plaintext io.Reader) ([]byte, error) {
	limited := &io.LimitedReader{
		R: plaintext,
		N: MaxPlaintext + 1, // detect truncated message
	}

	b, err := ioutil.ReadAll(limited)
	if err != nil {
		return nil, fmt.Errorf("reading plaintext: %w", err)
	}

	if len(b) > MaxPlaintext {
		return nil, fmt.Errorf("plaintext larger than %d", MaxPlaintext)
	}

	out, err := box.SealAnonymous(nil, b, &pub, rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("encrypting plaintext: %w", err)
	}

	block := pem.Block{
		Type: PEMType,
		Headers: map[string]string{
			HeaderReference: reference,
		},
		Bytes: out,
	}
	b = pem.EncodeToMemory(&block)
	if err != nil {
		return nil, fmt.Errorf("PEM-encoding ciphertext: %w", err)
	}
	return b, nil
}

func Decrypt(ciphertextPEM []byte, getKeys func(reference string) (Keys, error)) ([]byte, error) {
	block, _ := pem.Decode(ciphertextPEM)
	if block.Type != PEMType {
		return nil, fmt.Errorf("not a handoff message")
	}

	reference, present := block.Headers["Reference"]
	if !present || reference == "" {
		return nil, fmt.Errorf("no reference in message")
	}

	keys, err := getKeys(reference)
	if err != nil {
		return nil, fmt.Errorf("getting keys for %s: %w", reference, err)
	}

	pub, _, err := ParsePubKey(keys.Public)
	if err != nil {
		return nil, fmt.Errorf("parsing retrieved public key for %s: %w", reference, err)
	}

	priv, err := ParsePrivKey(keys.Private)
	if err != nil {
		return nil, fmt.Errorf("parsing retrieved private key for %s: %w", reference, err)
	}

	message, ok := box.OpenAnonymous(nil, block.Bytes, &pub, &priv)
	if !ok {
		return nil, fmt.Errorf("decrypting message failed")
	}

	return message, nil
}

package main

import (
	"bytes"
	"fmt"
	"syscall/js"

	"github.com/jasonmf/handoff"
)

func encrypt(pubkey, secret string) (string, error) {
	pub, reference, err := handoff.ParsePubKey(pubkey)
	if err != nil {
		return "", fmt.Errorf("parsing public key: %w", err)
	}
	secretBuf := bytes.NewReader([]byte(secret))
	encrypted, err := handoff.Encrypt(pub, reference, secretBuf)
	if err != nil {
		return "", fmt.Errorf("encrypting secret: %w", err)
	}
	return string(encrypted), nil
}

func wrapEncrypt() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) != 2 {
			return "requires two arguments"
		}
		inputKey := args[0].String()
		inputSecret := args[1].String()
		encrypted, err := encrypt(inputKey, inputSecret)
		if err != nil {
			return err.Error()
		}
		return encrypted
	})
}

func generate() (handoff.Keys, error) {
	return handoff.Generate("handoff")
}

func wrapGenerate() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) != 0 {
			return "expects no arguments"
		}
		keys, err := generate()
		if err != nil {
			return err.Error()
		}
		return map[string]interface{}{
			"private": keys.Private,
			"public":  keys.Public,
		}
	})
}

func decrypt(cipherText, pubKey, privKey string) (string, error) {
	keyFn := func(_ string) (handoff.Keys, error) {
		return handoff.Keys{
			Private: privKey,
			Public:  pubKey,
		}, nil
	}

	b, err := handoff.Decrypt([]byte(cipherText), keyFn)
	return string(b), err
}

func wrapDecrypt() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) != 3 {
			return "requires three arguments"
		}
		cipherText := args[0].String()
		pubKey := args[1].String()
		privKey := args[2].String()
		decrypted, err := decrypt(cipherText, pubKey, privKey)
		if err != nil {
			return err.Error()
		}
		return decrypted
	})
}

func main() {
	fmt.Println("Go Web Assembly")
	js.Global().Set("handoffDecrypt", wrapDecrypt())
	js.Global().Set("handoffEncrypt", wrapEncrypt())
	js.Global().Set("handoffGenerate", wrapGenerate())
	<-make(chan bool)
}

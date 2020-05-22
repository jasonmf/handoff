package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/AgentZombie/handoff"
)

func fatalIfError(err error, msg string) {
	if err != nil {
		log.Fatal("error ", msg, ": ", err)
	}
}

var (
	fGenerate = flag.String("generate", "", "generate a new key with this reference")
	fEncrypt  = flag.String("encrypt", "", "encrypt STDIN using this key")
	fDecrypt  = flag.String("decrypt", "", "decrypt from STDIN to this file")
)

func main() {
	flag.Parse()

	modeCnt := 0
	for _, f := range []*string{fGenerate, fEncrypt, fDecrypt} {
		if *f != "" {
			modeCnt++
		}
	}
	if modeCnt != 1 {
		fmt.Fprintln(os.Stderr, "must supply one of -generate, -encrypt, or -decrypt")
		flag.Usage()
		os.Exit(-1)
	}

	switch {
	case *fGenerate != "":
		Generate()
	case *fEncrypt != "":
		Encrypt()
	case *fDecrypt != "":
		Decrypt()
	}
}

func Generate() {
	outName := "handoff-" + *fGenerate + ".json"
	if _, err := os.Stat(outName); err == nil {
		log.Fatal("already exists: " + outName)
	}

	keys, err := handoff.Generate(*fGenerate)

	b, err := json.Marshal(keys)
	fatalIfError(err, "serializing keys")

	fatalIfError(ioutil.WriteFile(outName, b, 0400), "writing key file "+outName)

	fmt.Println("send to user: " + keys.Public)
}

func Encrypt() {
	key, reference, err := handoff.ParsePubKey(*fEncrypt)
	fatalIfError(err, "parsing public key")

	fmt.Printf("reading for %s from STDIN\n", reference)
	b, err := handoff.Encrypt(key, reference, os.Stdin)
	fatalIfError(err, "encrypting from STDIN")

	fmt.Println("send in response:")
	fmt.Println(string(b))
}

func Decrypt() {
	fmt.Println("reading message from STDIN, writing to " + *fDecrypt)
	b, err := ioutil.ReadAll(os.Stdin)
	fatalIfError(err, "reading from STDIN")

	b, err = handoff.Decrypt(b, getKeys)
	fatalIfError(err, "decrypting")

	fatalIfError(ioutil.WriteFile(*fDecrypt, b, 0400), "writing output to "+*fDecrypt)
}

func getKeys(reference string) (handoff.Keys, error) {
	var keys handoff.Keys

	b, err := ioutil.ReadFile("handoff-" + reference + ".json")
	if err != nil {
		return keys, fmt.Errorf("reading key file: %w", err)
	}

	if err := json.Unmarshal(b, &keys); err != nil {
		return keys, fmt.Errorf("parsing key file: %w", err)
	}
	return keys, nil
}

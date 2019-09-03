package main

import (
	"encoding/json"
	"github.com/btcsuite/btcutil/base58"
	"log"
	"lucastetreault/did-tangaroa/pkg/did"
)

func main() {
	log.SetFlags(0)

	ddoc, priv, err := did.NewDocument()
	if err != nil {
		panic(err.Error())
	}

	b, err := json.Marshal(ddoc)
	if err != nil {
		panic(err.Error())
	}
	log.Println("DID Document:")
	log.Println(string(b) + "\n")

	log.Println("Private Key:")
	log.Println(base58.Encode(priv))
}

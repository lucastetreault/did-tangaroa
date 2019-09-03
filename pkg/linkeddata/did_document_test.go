package linkeddata

import (
	"encoding/base64"
	"encoding/json"
	"github.com/btcsuite/btcutil/base58"
	"log"
	"testing"
)

func TestNewDocument(t *testing.T) {
	d, priv, err := NewDocument()
	if err != nil {
		t.Errorf(err.Error())
	}

	v := Verify(d, base58.Decode(d.PublicKey[0].PublicKeyBase58))
	if !v {
		t.Errorf("error verifying signature")
	}

	b, _ := json.Marshal(d)
	log.Println("DID DidDocument:")
	log.Println(string(b))

	log.Println("Private Key:")
	log.Println(base64.StdEncoding.EncodeToString(priv))
}

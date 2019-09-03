package linkeddata

import (
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"golang.org/x/crypto/ed25519"
	"time"
)

const Ed25519Signature2018 = "Ed25519Signature2018"

type Proof struct {
	Capability     *ObjectCapability `json:"capability,omitempty"`
	Created        string            `json:"created"`
	Creator        string            `json:"creator"`
	ProofPurpose   string            `json:"proofPurpose,omitempty"`
	SignatureValue string            `json:"signatureValue,omitempty"`
	Type           string            `json:"type,omitempty"`
}

type Signable interface {
	Clone() Signable
	GetProof() *Proof
	SetProof(*Proof)
}

func Sign(obj Signable, pk ed25519.PrivateKey) error {
	created, verifyHash := createVerifyHash(obj)

	sig := ed25519.Sign(pk, verifyHash)
	originalProof := obj.GetProof()
	originalProof.SignatureValue = base64.StdEncoding.EncodeToString(sig)
	originalProof.Created = created
	obj.SetProof(originalProof)
	return nil
}

func Verify(obj Signable, pub ed25519.PublicKey) bool {
	_, verifyHash := createVerifyHash(obj)

	sig, err := base64.StdEncoding.DecodeString(obj.GetProof().SignatureValue)
	if err != nil {
		return false
	}

	return ed25519.Verify(pub, verifyHash, sig)
}

func createVerifyHash(obj Signable) (string, []byte) {
	clone := obj.Clone()
	cloneProof := clone.GetProof()
	cloneProof.Type = ""
	cloneProof.SignatureValue = ""
	if cloneProof.Created == "" {
		cloneProof.Created = time.Now().UTC().Format("2006-01-02T03:04:05MST")
	}
	clone.SetProof(nil)
	proofHash, err := canonicalizeAndHash(cloneProof)
	if err != nil {
		//return err
	}
	tbsHash, err := canonicalizeAndHash(clone)
	if err != nil {
		//return err
	}
	verifyHash := append(proofHash, tbsHash...)
	return cloneProof.Created, verifyHash
}

func canonicalizeAndHash(tbc interface{}) ([]byte, error) {
	b, err := json.Marshal(tbc)
	if err != nil {
		return nil, err
	}

	var m map[string]interface{}
	err = json.Unmarshal(b, m)
	if err != nil {
		return nil, err
	}

	b, err = json.Marshal(tbc)
	if err != nil {
		return nil, err
	}

	h := sha512.New()
	h.Write(b)
	return h.Sum(nil), nil
}

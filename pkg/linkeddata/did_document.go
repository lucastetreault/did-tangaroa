package linkeddata

import (
	"encoding/json"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	"github.com/google/uuid"
	"golang.org/x/crypto/ed25519"
	"time"
)

type PublicKey struct {
	Id              string `json:"id"`
	Type            string `json:"type"`
	Controller      string `json:"controller"`
	PublicKeyBase58 string `json:"publicKeyBase58"`
}

type Service struct {
	Id              string `json:"id"`
	Type            string `json:"type"`
	ServiceEndpoint string `json:"serviceEndpoint"`
}

type DidDocument struct {
	Context        string      `json:"@context"`
	Id             string      `json:"id"`
	PublicKey      []PublicKey `json:"public_key,omitempty"`
	Authentication []string    `json:"authentication,omitempty"`
	Service        []Service   `json:"service,omitempty"`
	Created        string      `json:"created"`
	Updated        string      `json:"updated,omitempty"`
	Proof          *Proof      `json:"proof,omitempty"`
}

func new() string {
	u, _ := uuid.New().MarshalBinary()
	b := base58.Encode(u)
	return fmt.Sprintf("did:dtang:%s", b)
}

func NewDocument() (*DidDocument, ed25519.PrivateKey, error) {
	id := new()
	return GenerateDidDocument(id)
}

func GenerateDidDocument(id string) (*DidDocument, ed25519.PrivateKey, error) {
	pub, pk, err := ed25519.GenerateKey(nil)
	if err != nil {
		return nil, nil, err
	}

	dpub := PublicKey{
		Id:              fmt.Sprintf("%s#keys-1", id),
		Type:            "Ed25519VerificationKey2018",
		Controller:      id,
		PublicKeyBase58: base58.Encode(pub),
	}

	d := DidDocument{
		Context:        "https://www.w3.org/2019/did/v1",
		Id:             id,
		PublicKey:      []PublicKey{dpub},
		Authentication: []string{dpub.Id},
		Created:        time.Now().UTC().Format("2006-01-02T03:04:05MST"),
		Proof: &Proof{
			Created: time.Now().UTC().Format("2006-01-02T03:04:05MST"),
			Creator: dpub.Id,
			Type:    Ed25519Signature2018,
		},
	}

	err = Sign(&d, pk)
	if err != nil {
		return nil, nil, err
	}

	return &d, pk, nil
}

func (d *DidDocument) Clone() Signable {
	b, _ := json.Marshal(d)
	var clone DidDocument
	json.Unmarshal(b, &clone)
	return &clone
}

func (d *DidDocument) GetProof() *Proof {
	return d.Proof
}

func (d *DidDocument) SetProof(p *Proof) {
	d.Proof = p
}

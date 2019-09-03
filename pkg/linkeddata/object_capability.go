package linkeddata

import (
	"encoding/json"
	"github.com/google/uuid"
	"golang.org/x/crypto/ed25519"
)

type ObjectCapability struct {
	Id               uuid.UUID `json:"id"`
	ParentCapability string    `json:"parentCapability"`
	Invoker          string    `json:"invoker"`
	Proof            *Proof    `json:"proof,omitempty"`
}

func (d *ObjectCapability) Invoke(action string, pk ed25519.PrivateKey) (*ObjectCapabilityInvocation, error) {
	invocation := ObjectCapabilityInvocation{
		Id:     uuid.New(),
		Action: action,
		Proof: &Proof{
			ProofPurpose: "capabilityInvocation",
			Creator:      d.Invoker,
			Type:         Ed25519Signature2018,
			Capability:   d,
		},
	}

	err := Sign(&invocation, pk)
	if err != nil {
		return nil, err
	}
	return &invocation, nil
}

func (d *ObjectCapability) Clone() Signable {
	b, _ := json.Marshal(d)
	var clone DidDocument
	json.Unmarshal(b, &clone)
	return &clone
}

func (d *ObjectCapability) GetProof() *Proof {
	return d.Proof
}

func (d *ObjectCapability) SetProof(p *Proof) {
	d.Proof = p
}

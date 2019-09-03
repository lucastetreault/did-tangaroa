package linkeddata

import (
	"encoding/json"
	"github.com/google/uuid"
)

type ObjectCapabilityInvocation struct {
	Id     uuid.UUID `json:"id"`
	Action string    `json:"action"`
	Proof  *Proof    `json:"proof,omitempty"`
}

func (d *ObjectCapabilityInvocation) Clone() Signable {
	b, _ := json.Marshal(d)
	var clone DidDocument
	json.Unmarshal(b, &clone)
	return &clone
}

func (d *ObjectCapabilityInvocation) GetProof() *Proof {
	return d.Proof
}

func (d *ObjectCapabilityInvocation) SetProof(p *Proof) {
	d.Proof = p
}

// Copyright 2015 The etcd Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	"github.com/google/uuid"
	"golang.org/x/crypto/ed25519"
	"io/ioutil"
	"lucastetreault/did-tangaroa/pkg/linkeddata"
	"lucastetreault/did-tangaroa/raft/raftpb"
	"os"
)

func main() {
	genesis := flag.Bool("genesis", false, "set to true if this is the first node of the cluster being started for the first time")
	id := flag.String("id", "", "node ID")
	address := flag.String("address", "", "address for the genesis node")
	kvport := flag.Int("port", 9121, "key-value server port")
	join := flag.Bool("join", false, "join an existing cluster")
	flag.Parse()

	if id == nil || *id == "" {
		panic("id is required")
	}

	var pk ed25519.PrivateKey
	var ddoc *linkeddata.DidDocument
	if notExists(*id) {
		var pk ed25519.PrivateKey
		var err error
		ddoc, pk, err = linkeddata.GenerateDidDocument(*id)
		if err != nil {
			panic(err.Error())
		}
		storeDidDocument(".", ddoc.Id, ddoc, pk)
	} else {
		pk = readNodePrivateKey(*id, pk)
	}

	var clusterDdoc *linkeddata.DidDocument
	if isGenesisNode(genesis) && notExists("cluster") {
		clusterDdoc = bootstrapCluster(*id, *address)
	}

	proposeC := make(chan string)
	defer close(proposeC)
	confChangeC := make(chan raftpb.ConfChange)
	defer close(confChangeC)

	// raft provides a commit stream for the proposals from the http api
	var kvs *kvstore
	getSnapshot := func() ([]byte, error) { return kvs.getSnapshot() }
	commitC, errorC, snapshotterReady := newRaftNode(*id, pk, *address, *join, getSnapshot, proposeC, confChangeC)

	kvs = newKVStore(<-snapshotterReady, proposeC, commitC, errorC)

	if clusterDdoc != nil {
		saveDidDocument(clusterDdoc, kvs)
	}
	if ddoc != nil {
		saveDidDocument(ddoc, kvs)
	}

	// the key-value http handler will propose updates to raft
	serveHttpKVAPI(kvs, *kvport, confChangeC, errorC)
}

func readNodePrivateKey(id string, pk ed25519.PrivateKey) ed25519.PrivateKey {
	b, err := ioutil.ReadFile(fmt.Sprintf("./%s/%s.key", id, id))
	if err != nil {
		panic(err.Error())
	}
	pk = b
	return pk
}

func notExists(id string) bool {
	return !exists(fmt.Sprintf("./%s.json", id)) || !exists(fmt.Sprintf("/%s.key", id))
}

func saveDidDocument(ddoc *linkeddata.DidDocument, kvs *kvstore) {
	b, _ := json.MarshalIndent(ddoc, "", "    ")
	kvs.Propose(ddoc.Id, string(b))
}

func bootstrapCluster(id string, nodeUrl string) *linkeddata.DidDocument {
	genesisDdoc, pk, err := linkeddata.NewDocument()
	if err != nil {
		panic(err.Error())
	}
	genesisDdoc.Service = make([]linkeddata.Service, 0)
	genesisDdoc.Service = append(genesisDdoc.Service, linkeddata.Service{Id: id, Type: "DtangClusterNode", ServiceEndpoint: nodeUrl})
	err = linkeddata.Sign(genesisDdoc, pk)
	if err != nil {
		panic(err.Error())
	}
	storeDidDocument("./", "cluster", genesisDdoc, pk)

	ocap := linkeddata.ObjectCapability{
		Id:               uuid.New(),
		Invoker:          fmt.Sprintf("%s#keys-1", id),
		ParentCapability: genesisDdoc.Id,
		Proof: &linkeddata.Proof{
			Type:         "Ed25519Signature2018",
			Creator:      genesisDdoc.PublicKey[0].Id,
			ProofPurpose: "capabilityDelegation",
		},
	}
	err = linkeddata.Sign(&ocap, pk)
	if err != nil {
		panic(err.Error())
	}
	b, err := json.MarshalIndent(ocap, "", "    ")
	if err != nil {
		panic(err.Error())
	}
	ocapf, err := os.Create(fmt.Sprintf("./%s.ocap", id))
	if err != nil {
		panic(err.Error())
	}
	_, err = ocapf.Write(b)
	if err != nil {
		panic(err.Error())
	}
	err = ocapf.Close()
	if err != nil {
		panic(err.Error())
	}

	return genesisDdoc
}

func isGenesisNode(genesis *bool) bool {
	return genesis != nil && *genesis
}

func storeDidDocument(path string, name string, ddoc *linkeddata.DidDocument, pk ed25519.PrivateKey) {
	os.Mkdir(path, 0755)

	b, err := json.MarshalIndent(ddoc, "", "    ")
	if err != nil {
		panic(err.Error())
	}
	ddocf, err := os.Create(fmt.Sprintf("%s/%s.json", path, name))
	if err != nil {
		panic(err.Error())
	}
	_, err = ddocf.Write(b)
	if err != nil {
		panic(err.Error())
	}
	err = ddocf.Close()
	if err != nil {
		panic(err.Error())
	}
	pkf, err := os.Create(fmt.Sprintf("%s/%s.key", path, name))
	if err != nil {
		panic(err.Error())
	}
	_, err = pkf.Write([]byte(base58.Encode(pk)))
	if err != nil {
		panic(err.Error())
	}
	err = pkf.Close()
	if err != nil {
		panic(err.Error())
	}
	return
}

func exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

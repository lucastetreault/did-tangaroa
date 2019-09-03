package main

import (
	"encoding/json"
	"io/ioutil"
	"lucastetreault/did-tangaroa/pkg/linkeddata"
)

var clusterDdoc *linkeddata.DidDocument

func loadClusterDdoc() {
	if clusterDdoc == nil {
		b, err := ioutil.ReadFile("./cluster.json")
		if err != nil {
			panic(err.Error())
		}

		var ddoc linkeddata.DidDocument
		err = json.Unmarshal(b, &ddoc)
		if err != nil {
			panic(err.Error())
		}
		clusterDdoc = &ddoc
	}
}

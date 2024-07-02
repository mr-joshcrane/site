package store_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/mr-joshcrane/site/store"
)

func TestStoreSyft(t *testing.T) {
	s, err := store.NewSQLiteStore("test.db")
	if err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile("testdata/syft.json")
	if err != nil {
		t.Fatal(err)
	}
	syft := store.SyftDocument{}
	err = json.Unmarshal(data, &syft)
	if err != nil {
		t.Fatal(err)
	}
	err = s.StoreSyft(syft, "1234567890")
	if err != nil {
		t.Fatal(err)
	}

}

// {
// 	"id": "e185fe1c57dc5e3e",
// 	"location": {
// 		"path": "/.github/workflows/backstage-validator.yml"
// 	}
// },
// {
// 	"id": "0b6593b85bf207fb",
// 	"location": {
// 		"path": "/.github/workflows/integration-test.yml"
// 	}
// },

// type SyftDocument struct {
// 	Commit                string         `json:"commit"`
// 	Organisation          string         `json:"organisation"`
// 	Repository            string         `json:"repository"`
// 	Artifacts             []Package      `json:"artifacts"`
// 	ArtifactRelationships []Relationship `json:"artifactRelationships"`
// 	Files                 []File         `json:"files"`
// 	Source                Source         `json:"source"`
// 	Descriptor            Descriptor     `json:"descriptor"`
// 	Schema                Schema         `json:"schema"`
// 	PURL                  string         `json:"purl"`
// 	MetadataType          string         `json:"metadataType"`
// }

// {
// 	"id": "18253d8fcb32a252",
// 	"name": "@aashutoshrathi/word-wrap",
// 	"version": "1.2.6",
// 	"type": "npm",
// 	"foundBy": "javascript-lock-cataloger",
// 	"locations": [
// 		{
// 			"path": "/ops/yarn.lock",
// 			"accessPath": "/ops/yarn.lock",
// 			"annotations": {
// 				"evidence": "primary"
// 			}
// 		}
// 	],
// 	"licenses": [],
// 	"language": "javascript",
// 	"cpes": [
// 		{
// 			"cpe": "cpe:2.3:a:\\@aashutoshrathi\\/word-wrap:\\@aashutoshrathi\\/word-wrap:1.2.6:*:*:*:*:*:*:*",
// 			"source": "syft-generated"
// 		},
// 		{
// 			"cpe": "cpe:2.3:a:\\@aashutoshrathi\\/word-wrap:\\@aashutoshrathi\\/word_wrap:1.2.6:*:*:*:*:*:*:*",
// 			"source": "syft-generated"
// 		},
// 		{
// 			"cpe": "cpe:2.3:a:\\@aashutoshrathi\\/word_wrap:\\@aashutoshrathi\\/word-wrap:1.2.6:*:*:*:*:*:*:*",
// 			"source": "syft-generated"
// 		},
// 		{
// 			"cpe": "cpe:2.3:a:\\@aashutoshrathi\\/word_wrap:\\@aashutoshrathi\\/word_wrap:1.2.6:*:*:*:*:*:*:*",
// 			"source": "syft-generated"
// 		},
// 		{
// 			"cpe": "cpe:2.3:a:\\@aashutoshrathi\\/word:\\@aashutoshrathi\\/word-wrap:1.2.6:*:*:*:*:*:*:*",
// 			"source": "syft-generated"
// 		},
// 		{
// 			"cpe": "cpe:2.3:a:\\@aashutoshrathi\\/word:\\@aashutoshrathi\\/word_wrap:1.2.6:*:*:*:*:*:*:*",
// 			"source": "syft-generated"
// 		}
// 	],
// 	"purl": "pkg:npm/%40aashutoshrathi/word-wrap@1.2.6",
// 	"metadataType": "javascript-yarn-lock-entry",
// 	"metadata": {
// 		"resolved": "https://registry.yarnpkg.com/@aashutoshrathi/word-wrap/-/word-wrap-1.2.6.tgz#bd9154aec9983f77b3a034ecaa015c2e4201f6cf",
// 		"integrity": "sha512-1Yjs2SvM8TflER/OD3cOjhWWOZb58A2t7wpE2S9XfBYTiIl+XFhQG2bjy4Pu1I+EAlCNUzRDYDdFwFYUKvXcIA=="
// 	}
// },

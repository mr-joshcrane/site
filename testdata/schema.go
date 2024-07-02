// "Document": {
//   "properties": {
//     "artifacts": {
//       "items": {
//         "$ref": "#/$defs/Package"
//       },
//       "type": "array"
//     },
//     "artifactRelationships": {
//       "items": {
//         "$ref": "#/$defs/Relationship"
//       },
//       "type": "array"
//     },
//     "files": {
//       "items": {
//         "$ref": "#/$defs/File"
//       },
//       "type": "array"
//     },
//     "source": {
//       "$ref": "#/$defs/Source"
//     },
//     "distro": {
//       "$ref": "#/$defs/LinuxRelease"
//     },
//     "descriptor": {
//       "$ref": "#/$defs/Descriptor"
//     },
//     "schema": {
//       "$ref": "#/$defs/Schema"
//     }
//   },
//   "type": "object",
//   "required": [
//     "artifacts",
//     "artifactRelationships",
//     "source",
//     "distro",
//     "descriptor",
//     "schema"
//   ]
// }

// {
//   "id": "18253d8fcb32a252",
//   "name": "@aashutoshrathi/word-wrap",
//   "version": "1.2.6",
//   "type": "npm",
//   "foundBy": "javascript-lock-cataloger",
//   "locations": [
//       {
//           "path": "/ops/yarn.lock",
//           "accessPath": "/ops/yarn.lock",
//           "annotations": {
//               "evidence": "primary"
//           }
//       }
//   ],
//   "licenses": [],
//   "language": "javascript",
//   "cpes": [
//       {
//           "cpe": "cpe:2.3:a:\\@aashutoshrathi\\/word-wrap:\\@aashutoshrathi\\/word-wrap:1.2.6:*:*:*:*:*:*:*",
//           "source": "syft-generated"
//       },
//       {
//           "cpe": "cpe:2.3:a:\\@aashutoshrathi\\/word-wrap:\\@aashutoshrathi\\/word_wrap:1.2.6:*:*:*:*:*:*:*",
//           "source": "syft-generated"
//       },
//       {
//           "cpe": "cpe:2.3:a:\\@aashutoshrathi\\/word_wrap:\\@aashutoshrathi\\/word-wrap:1.2.6:*:*:*:*:*:*:*",
//           "source": "syft-generated"
//       },
//       {
//           "cpe": "cpe:2.3:a:\\@aashutoshrathi\\/word_wrap:\\@aashutoshrathi\\/word_wrap:1.2.6:*:*:*:*:*:*:*",
//           "source": "syft-generated"
//       },
//       {
//           "cpe": "cpe:2.3:a:\\@aashutoshrathi\\/word:\\@aashutoshrathi\\/word-wrap:1.2.6:*:*:*:*:*:*:*",
//           "source": "syft-generated"
//       },
//       {
//           "cpe": "cpe:2.3:a:\\@aashutoshrathi\\/word:\\@aashutoshrathi\\/word_wrap:1.2.6:*:*:*:*:*:*:*",
//           "source": "syft-generated"
//       }
//   ],
//   "purl": "pkg:npm/%40aashutoshrathi/word-wrap@1.2.6",
//   "metadataType": "javascript-yarn-lock-entry",
//   "metadata": {
//       "resolved": "https://registry.yarnpkg.com/@aashutoshrathi/word-wrap/-/word-wrap-1.2.6.tgz#bd9154aec9983f77b3a034ecaa015c2e4201f6cf",
//       "integrity": "sha512-1Yjs2SvM8TflER/OD3cOjhWWOZb58A2t7wpE2S9XfBYTiIl+XFhQG2bjy4Pu1I+EAlCNUzRDYDdFwFYUKvXcIA=="
//   }
// },

package sitem

type SyftDocument struct {
	Artifacts             []Package      `json:"artifacts"`
	ArtifactRelationships []Relationship `json:"artifactRelationships"`
	Files                 []File         `json:"files"`
	Source                Source         `json:"source"`
	Descriptor            Descriptor     `json:"descriptor"`
	Schema                Schema         `json:"schema"`
	PURL                  string         `json:"purl"`
	MetadataType          string         `json:"metadataType"`
}

type Package struct {
	Name      string     `json:"name"`
	Version   string     `json:"version"`
	Type      string     `json:"type"`
	FoundBy   string     `json:"foundBy"`
	Locations []Location `json:"locations"`
	Language  string     `json:"language"`
}

type Relationship struct {
	Source  string `json:"source"`
	Target  string `json:"target"`
	Type    string `json:"type"`
	Ref     string `json:"ref"`
	RefType string `json:"refType"`
}

type Location struct {
	Path        string `json:"path"`
	AccessPath  string `json:"accessPath"`
	Annotations struct {
		Evidence string `json:"evidence"`
	} `json:"annotations"`
}

type File struct {
	Path string `json:"path"`
}

type Source struct {
	Target string `json:"target"`
}

type Descriptor struct {
	LayerInfo struct {
		LayerIndex int `json:"layerIndex"`
	} `json:"layerInfo"`
}

type Schema struct {
	SchemaVersion string `json:"schemaVersion"`
}

// CREATE TABLE syft_document (
//     id INT PRIMARY KEY AUTO_INCREMENT,
//     purl VARCHAR(255) NOT NULL,
//     metadata_type VARCHAR(50) NOT NULL,
//     schema_version VARCHAR(50) NOT NULL
// );

// CREATE TABLE package (
//     id INT PRIMARY KEY AUTO_INCREMENT,
//     document_id INT,
//     name VARCHAR(100) NOT NULL,
//     version VARCHAR(50) NOT NULL,
//     type VARCHAR(50) NOT NULL,
//     found_by VARCHAR(100) NOT NULL,
//     language VARCHAR(50) NOT NULL,
//     FOREIGN KEY (document_id) REFERENCES syft_document(id) ON DELETE CASCADE
// );

// CREATE TABLE location (
//     id INT PRIMARY KEY AUTO_INCREMENT,
//     package_id INT,
//     path VARCHAR(255) NOT NULL,
//     access_path VARCHAR(255) NOT NULL,
//     evidence VARCHAR(50) NOT NULL,
//     FOREIGN KEY (package_id) REFERENCES package(id) ON DELETE CASCADE
// );

// CREATE TABLE relationship (
//     id INT PRIMARY KEY AUTO_INCREMENT,
//     document_id INT,
//     source VARCHAR(100) NOT NULL,
//     target VARCHAR(100) NOT NULL,
//     type VARCHAR(50) NOT NULL,
//     ref VARCHAR(100),
//     ref_type VARCHAR(50),
//     FOREIGN KEY (document_id) REFERENCES syft_document(id) ON DELETE CASCADE
// );

// CREATE TABLE file (
//     id INT PRIMARY KEY AUTO_INCREMENT,
//     document_id INT,
//     path VARCHAR(255) NOT NULL,
//     FOREIGN KEY (document_id) REFERENCES syft_document(id) ON DELETE CASCADE
// );
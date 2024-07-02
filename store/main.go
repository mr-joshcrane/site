package store

import (
	"database/sql"
	"encoding/json"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Store interface {
	ListRepos() (RepositoriesModel, error)
	AddRepo(org string, repo string) error
	AddFullRepo(repo RepositoryModel) error
	GetRepo(org string, repo string) (RepositoryModel, error)
	StoreSyft(document SyftDocument) error
	GetSyft(commit string) (SyftDocument, error)
	// StoreGrype(document GrypeDocument) error
	// GetGrype(commit string) (GrypeDocument, error)
}

// type InMemoryStore struct {
// 	repos map[string]RepositoryModel
// }

// func NewMemoryStore() InMemoryStore {
// 	return InMemoryStore{
// 		repos: make(map[string]RepositoryModel),
// 	}
// }

// func (s *InMemoryStore) ListRepos() (RepositoriesModel, error) {
// 	var repos RepositoriesModel
// 	for _, r := range s.repos {
// 		repos = append(repos, r)
// 	}
// 	return repos, nil
// }

// func (s *InMemoryStore) AddRepo(org string, repo string) error {
// 	s.repos[org+"/"+repo] = RepositoryModel{
// 		Org:         org,
// 		Name:        repo,
// 		LastScanned: time.Now(),
// 	}
// 	return nil
// }

// func (s *InMemoryStore) AddFullRepo(repo RepositoryModel) error {
// 	s.repos[repo.Org+"/"+repo.Name] = repo
// 	return nil
// }

// func (s *InMemoryStore) GetRepo(org string, repo string) (RepositoryModel, error) {
// 	r, ok := s.repos[org+"/"+repo]
// 	if !ok {
// 		return RepositoryModel{}, nil
// 	}
// 	return r, nil
// }

// func (s *InMemoryStore) AddDependency(document SyftDocument) error {
// 	return nil
// }

// func (s *InMemoryStore) GetDependency(purl string) (SyftDocument, error) {
// 	return SyftDocument{}, nil
// }

type SQLiteStore struct {
	db *sql.DB
}

func NewSQLiteStore(filepath string) (SQLiteStore, error) {
	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		return SQLiteStore{}, err
	}
	err = createTables(db)
	if err != nil {
		return SQLiteStore{}, err
	}
	if err != nil {
		return SQLiteStore{}, err
	}
	return SQLiteStore{
		db: db,
	}, nil
}

func createTables(db *sql.DB) error {
	createReposTable := `
		CREATE TABLE IF NOT EXISTS repos (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			org VARCHAR(255) NOT NULL,
			repo VARCHAR(255) NOT NULL,
			created_on TIMESTAMP NOT NULL,
			updated_on TIMESTAMP NOT NULL,
			head VARCHAR(255) NOT NULL,
			last_scanned TIMESTAMP NOT NULL,
			archived BOOLEAN NOT NULL
		);`
	createSyftDocumentTable := `
		CREATE TABLE IF NOT EXISTS syft_document (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			repo FOREIGN KEY REFERENCES repos(id),
			commit VARCHAR(255) NOT NULL,
			data BLOB NOT NULL,
		);`
	createGrypeDocumentTable := `
		CREATE TABLE IF NOT EXISTS grype_document (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			repo FOREIGN KEY REFERENCES repos(id),
			commit VARCHAR(255) NOT NULL,
			data BLOB NOT NULL,
		);`
	_, err := db.Exec(createReposTable)
	if err != nil {
		return err
	}
	_, err = db.Exec(createSyftDocumentTable)
	if err != nil {
		return err
	}
	_, err = db.Exec(createGrypeDocumentTable)
	if err != nil {
		return err
	}
	return nil
}

type RepositoryModel struct {
	Org         string    `json:"org"`
	Name        string    `json:"repo"`
	CreatedOn   time.Time `json:"created_on"`
	UpdatedOn   time.Time `json:"updated_on"`
	HEAD        string    `json:"head"`
	LastScanned time.Time `json:"last_scanned"`
	Archived    bool      `json:"archived"`
}

type RepositoriesModel []RepositoryModel

func (r RepositoriesModel) Len() int { return len(r) }
func (r RepositoriesModel) Less(i, j int) bool {
	UpdatedOnI := r[i].UpdatedOn.Unix()
	UpdatedOnJ := r[j].UpdatedOn.Unix()
	return UpdatedOnI > UpdatedOnJ
}
func (r RepositoriesModel) Swap(i, j int) { r[i], r[j] = r[j], r[i] }

func (r RepositoriesModel) Names() []string {
	var names []string
	for _, repo := range r {
		names = append(names, repo.Org+"/"+repo.Name)
	}
	return names
}

func (s *SQLiteStore) ListRepos() (RepositoriesModel, error) {
	rows, err := s.db.Query("SELECT * FROM repos")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var repos RepositoriesModel
	for rows.Next() {
		var r RepositoryModel
		err = rows.Scan(&r.Org, &r.Name, &r.CreatedOn, &r.UpdatedOn, &r.HEAD, &r.LastScanned, &r.Archived)
		if err != nil {
			return nil, err
		}
		repos = append(repos, r)
	}
	return repos, nil
}

func (s *SQLiteStore) AddRepo(org string, repo string) error {
	_, err := s.db.Exec("INSERT INTO repos (org, repo) VALUES (?, ?)", org, repo)
	return err
}

func (s *SQLiteStore) AddFullRepo(repo RepositoryModel) error {
	_, err := s.db.Exec(
		"INSERT INTO repos (org, repo, created_on, updated_on, head, last_scanned, archived) VALUES (?, ?, ?, ?, ?, ?, ?)",
		repo.Org, repo.Name, repo.CreatedOn, repo.UpdatedOn, repo.HEAD, repo.LastScanned, repo.Archived,
	)
	return err
}

func (s *SQLiteStore) GetRepo(org string, repo string) (RepositoryModel, error) {
	var r RepositoryModel
	err := s.db.QueryRow("SELECT * FROM repos WHERE org = ? AND repo = ?", org, repo).Scan(
		&r.Org, &r.Name, &r.CreatedOn, &r.UpdatedOn, &r.HEAD, &r.LastScanned, &r.Archived,
	)
	if err != nil {
		return RepositoryModel{}, err
	}
	return r, nil
}

func (s *SQLiteStore) StoreSyft(document SyftDocument) error {
	data, err := json.Marshal(document)
	if err != nil {
		return err
	}
	_, err = s.db.Exec("INSERT INTO syft_document (repo, commit, data) VALUES (?, ?, ?)", document.PURL, document.MetadataType, document.Schema.SchemaVersion, data)
	return err
}

func (s *SQLiteStore) GetSyft(commit string) (SyftDocument, error) {
	row := s.db.QueryRow("SELECT * FROM syft_document WHERE id = ?", commit)
	var document SyftDocument
	//PURL doesnt exist
	err := row.Scan()
	if err != nil {
		return SyftDocument{}, err
	}
	return document, nil
}

type SyftDocument struct {
	Commit                string         `json:"commit"`
	Organisation          string         `json:"organisation"`
	Repository            string         `json:"repository"`
	Artifacts             []Package      `json:"artifacts"`
	ArtifactRelationships []Relationship `json:"artifactRelationships"`
	Files                 []File         `json:"files"`
	Source                Source         `json:"source"`
	Descriptor            Descriptor     `json:"descriptor"`
	Schema                Schema         `json:"schema"`
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
	Path       string `json:"path"`
	AccessPath string `json:"accessPath"`
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

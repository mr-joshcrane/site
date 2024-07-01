package store

import "time"

type InMemoryStore struct {
	repos map[string]RepositoryModel
}

func NewMemoryStore() InMemoryStore {
	return InMemoryStore{
		repos: make(map[string]RepositoryModel),
	}
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

type Store interface {
	ListRepos() (RepositoriesModel, error)
	AddRepo(org string, repo string) error
	AddFullRepo(repo RepositoryModel) error
	GetRepo(org string, repo string) (RepositoryModel, error)
}

func (s *InMemoryStore) ListRepos() (RepositoriesModel, error) {
	var repos RepositoriesModel
	for _, r := range s.repos {
		repos = append(repos, r)
	}
	return repos, nil
}

func (s *InMemoryStore) AddRepo(org string, repo string) error {
	s.repos[org+"/"+repo] = RepositoryModel{
		Org:         org,
		Name:        repo,
		LastScanned: time.Now(),
	}
	return nil
}

func (s *InMemoryStore) AddFullRepo(repo RepositoryModel) error {
	s.repos[repo.Org+"/"+repo.Name] = repo
	return nil
}

func (s *InMemoryStore) GetRepo(org string, repo string) (RepositoryModel, error) {
	r, ok := s.repos[org+"/"+repo]
	if !ok {
		return RepositoryModel{}, nil
	}
	return r, nil
}

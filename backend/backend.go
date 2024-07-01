package backend

import (
	"fmt"
	"os"
	"time"

	"github.com/mr-joshcrane/site/store"
)

func Workers(s store.Store) error {
	for {
		err := ScanRepos(s)
		if err != nil {
			return err
		}
		repos, err := s.ListRepos()
		if err != nil {
			return err
		}
		for _, repo := range repos {
			fmt.Println("Processing repo:", repo)
		}
		time.Sleep(120 * time.Second)
	}

}

func ScanRepos(s store.Store) error {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return fmt.Errorf("GITHUB_TOKEN not set")
	}
	c := NewGithubClient(token)
	repos, err := GetRepositories(c)
	if err != nil {
		return err
	}
	r, err := Do(repos)
	if err != nil {
		return err
	}
	for _, repo := range r {
		r := store.RepositoryModel{
			Org:         repo.Org,
			Name:        repo.Name,
			CreatedOn:   repo.CreatedOn,
			UpdatedOn:   repo.UpdatedOn,
			HEAD:        repo.HEAD,
			LastScanned: repo.LastScanned,
			Archived:    repo.Archived,
		}
		err := s.AddFullRepo(r)
		if err != nil {
			return err
		}
	}
	return nil
}

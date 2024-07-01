package backend

import (
	"context"
	"os"
	"sync"
	"time"

	"github.com/google/go-github/v62/github"
	"github.com/mr-joshcrane/site/store"
)

func NewGithubClient(authToken string) *github.Client {
	if authToken == "" {
		authToken = os.Getenv("GITHUB_TOKEN")
	}
	return github.NewClient(nil).WithAuthToken(authToken)
}

func GetRepositories(c *github.Client) ([]*github.Repository, error) {
	var repositories []*github.Repository
	opt := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}

	for {
		repos, resp, err := c.Repositories.ListByOrg(context.Background(), "cultureamp", opt)
		if err != nil {
			return nil, err
		}
		repositories = append(repositories, repos...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return repositories, nil
}

func Do(repos []*github.Repository) ([]store.RepositoryModel, error) {
	wg := sync.WaitGroup{}
	wg.Add(len(repos))
	ch := make(chan store.RepositoryModel, len(repos))
	var now = time.Now()
	for _, repo := range repos {
		go func(ch chan store.RepositoryModel, repo *github.Repository) {
			defer wg.Done()
			r := store.RepositoryModel{
				Org:         *repo.Owner.Login,
				Name:        *repo.Name,
				CreatedOn:   repo.GetCreatedAt().Time,
				UpdatedOn:   repo.GetUpdatedAt().Time,
				LastScanned: now,
				Archived:    repo.GetArchived(),
			}
			c := NewGithubClient("")
			commit, _, err := c.Repositories.GetCommit(context.Background(), "cultureamp", *repo.Name, *repo.DefaultBranch, nil)
			if err != nil {
				r.HEAD = "error"
			} else {
				r.HEAD = *commit.SHA
			}
			ch <- r

		}(ch, repo)
	}

	wg.Wait()
	var repositories []store.RepositoryModel
	for i := 0; i < len(repos); i++ {
		repositories = append(repositories, <-ch)
	}
	return repositories, nil
}

type Repository struct {
	Org         string `json:"org"`
	Name        string `json:"repo"`
	CreatedOn   string `json:"created_on"`
	UpdatedOn   string `json:"updated_on"`
	HEAD        string `json:"head"`
	LastScanned string `json:"last_scanned"`
	Archived    bool   `json:"archived"`
}

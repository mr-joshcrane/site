package backend

import (
	"bytes"
	"context"
	"fmt"
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
	var total int
	var repositories []*github.Repository
	opt := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{PerPage: 1},
		Sort:        "updated",
	}

	for {
		total += 1
		if total > 10 {
			break
		}
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

func GetArchiveZip(c *github.Client, owner, repo, ref string) (*bytes.Buffer, error) {
	archiveURL, _, err := c.Repositories.GetArchiveLink(context.Background(), owner, repo, github.Zipball, &github.RepositoryContentGetOptions{Ref: ref}, 1)
	if err != nil {
		return nil, err
	}
	archive, err := c.Client().Get(archiveURL.String())
	if err != nil {
		return nil, err
	}
	fmt.Print(archive)
	defer archive.Body.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(archive.Body)

	return buf, nil
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

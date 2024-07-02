package backend

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/mr-joshcrane/site/store"
)

func Workers(s store.Store) error {
	for {
		// err := ScanRepos(s)
		// if err != nil {
		// 	return err
		// }
		repos, err := s.ListRepos()
		if err != nil {
			return err
		}
		for _, repo := range repos {
			fmt.Println("Processing repo:", repo)
		}
		time.Sleep(1 * time.Hour)
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
		time.Sleep(200 * time.Millisecond)
	}
	return nil
}

func GitCloneRepo(org string, repo string, commit string) (*os.File, error) {
	// git clone
	filename := fmt.Sprintf("%s.%s.%s.tar.gz", org, repo, commit)
	command := fmt.Sprintf("git archive --format=tar.gz --output=%s %s $(git ls-remote --get-url", filename, commit)
	commands := strings.Split(command, " ")
	// clone
	cmd := exec.Command(commands[0], commands[1:]...)
	err := cmd.Run()
	if err != nil {
		return nil, err
	}
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func Syft(file *os.File) (io.Reader, error) {
	// file is a tar.gz
	tempDir, err := os.MkdirTemp("", "syft")
	if err != nil {
		return nil, err
	}
	// extract the tar.gz
	cmd := exec.Command("tar", "-xzf", file.Name(), "-C", tempDir)
	err = cmd.Run()
	if err != nil {
		return nil, err
	}
	cmd = exec.Command("syft", "scan", tempDir, "-o", "json")
	data := new(bytes.Buffer)
	cmd.Stdout = data
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return nil, err
	}
	return data, nil
}

func ParseSyft(data io.Reader, origin string) (store.SyftDocument, error) {
	var document store.SyftDocument
	d, err := io.ReadAll(data)
	if err != nil {
		return store.SyftDocument{}, err
	}
	err = json.Unmarshal(d, document)
	if err != nil {
		return store.SyftDocument{}, err
	}
	commit, organization, repository, err := originSplit(origin)
	if err != nil {
		return store.SyftDocument{}, err
	}
	document.Commit = commit
	document.Organisation = organization
	document.Repository = repository
	return document, nil
}

func originSplit(origin string) (string, string, string, error) {
	s := strings.Split(origin, "/")
	if len(s) != 3 {
		return "", "", "", fmt.Errorf("origin must be in the format org/repo/commit")
	}
	return s[0], s[1], s[2], nil
}

func Grype(org string, repo string) error {
	dbDir := fmt.Sprintf("db/%s/%s/.", org, repo)
	cmd := exec.Command("grype", "db", "vuln", "--config", "db/grype-config.yaml", "--output", dbDir+"/grype.json", dbDir+"/syft.json")
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

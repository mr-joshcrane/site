package site

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/mr-joshcrane/site/store"
)

func StatusInternalServerError(err error) events.LambdaFunctionURLResponse {
	return events.LambdaFunctionURLResponse{
		StatusCode: http.StatusInternalServerError,
		Body:       "Internal Server Error: " + err.Error(),
		Headers: map[string]string{
			"Content-Type": "text/plain",
		},
	}
}

func StoreHandler(s store.Store) func(ctx context.Context, request events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	return func(ctx context.Context, request events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
		route := request.RawPath
		fmt.Println(route)
		switch route {
		case "/scan":
			fmt.Println("Scanning")
			return Scan(ctx, s, request)
		case "/addRepo":
			fmt.Println("Adding Repo")
			return AddRepo(ctx, s, request)
		case "/repos":
			fmt.Println("Getting Repos")
			return Repos(ctx, s, request)
		default:
			fmt.Println("Index")
			return indexPage("JoshNyk")
		}
	}
}

type Component interface {
	Render(ctx context.Context, w io.Writer) error
}

func componentToHTML(c Component) (string, error) {
	buf := new(bytes.Buffer)
	err := c.Render(context.Background(), buf)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func indexPage(title string) (events.LambdaFunctionURLResponse, error) {
	index := Index(title)
	page, err := componentToHTML(index)
	if err != nil {
		return StatusInternalServerError(err), nil
	}
	return events.LambdaFunctionURLResponse{
		StatusCode: 200,
		Body:       page,
		Headers: map[string]string{
			"Content-Type": "text/plain",
		},
	}, nil
}

func Repos(ctx context.Context, s store.Store, request events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	repoList, err := s.ListRepos()
	if err != nil {
		return StatusInternalServerError(err), nil
	}
	sort.Sort(repoList)
	repoComponent := RepositoryList(repoList)
	page, err := componentToHTML(repoComponent)
	if err != nil {
		return StatusInternalServerError(err), nil
	}
	return events.LambdaFunctionURLResponse{
		StatusCode: 200,
		Body:       page,
		Headers: map[string]string{
			"Content-Type": "text/plain",
		},
	}, nil
}

func AddRepo(ctx context.Context, s store.Store, request events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	org, repo, err := ParseGithubRepoURL(request.Body)
	if err != nil {
		return StatusInternalServerError(err), nil
	}
	err = s.AddRepo(org, repo)
	if err != nil {
		return StatusInternalServerError(err), nil
	}
	fmt.Println("Adding Repo", org, repo)

	partialRepository := store.RepositoryModel{
		Org:  org,
		Name: repo,
	}
	r := RepositoryComponent(partialRepository)
	repoHTML, err := componentToHTML(r)
	if err != nil {
		return StatusInternalServerError(err), nil
	}
	return events.LambdaFunctionURLResponse{
		StatusCode: 200,
		Body:       repoHTML,
		Headers: map[string]string{
			"Content-Type": "text/plain",
		},
	}, nil

}

func Scan(ctx context.Context, s store.Store, request events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	userInput := request.Headers["Hx-Trigger-Name"]
	userInput, err := url.QueryUnescape(userInput)
	if err != nil {
		return StatusInternalServerError(err), nil
	}
	userInput = strings.TrimSpace(userInput)
	org, repo, err := ParseGithubRepoURL(userInput)
	if err != nil {
		fmt.Println(err)
		return StatusInternalServerError(err), nil
	}
	err = GetVulnerabilityData(org, repo)
	if err != nil {
		fmt.Println(err)
		return StatusInternalServerError(err), nil
	}
	return events.LambdaFunctionURLResponse{
		StatusCode: 204,
		Body:       "OK",
		Headers: map[string]string{
			"Content-Type": "text/plain",
		},
	}, nil
}

func parseSyft(path string) (*SyftFinding, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var syft SyftFinding
	err = json.Unmarshal(data, &syft)
	if err != nil {
		return nil, err
	}
	return &syft, nil
}

func parseCycloneDX(path string) (*CycloneDX, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var vulns CycloneDX
	err = json.Unmarshal(data, &vulns)
	if err != nil {
		return nil, err
	}
	return &vulns, nil
}

func parseSHA(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

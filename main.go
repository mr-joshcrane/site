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
	"strings"

	"github.com/aws/aws-lambda-go/events"
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

func Handler(ctx context.Context, request events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	fmt.Println("CALLED HANDLER " + request.RawPath)
	route := request.RawPath
	switch route {
	case "/scan":
		return Scan(ctx, request)
	case "/addRepo":
		return AddRepo(ctx, request)
	case "/result":

		return Result(ctx, request)
	default:
		return indexPage("JoshNyk")
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

func Result(ctx context.Context, request events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	org, repo, err := ParseGithubRepoURL(request.Headers["Hx-Trigger-Name"])
	if err != nil {
		return StatusInternalServerError(err), nil
	}
	var findings RepositoryFindings
	findings.Organisation = org
	findings.Repository = repo
	project := fmt.Sprintf("./db/%s/%s/", org, repo)
	sha, err := parseSHA(project + "sha")
	if err != nil {
		return StatusInternalServerError(err), nil
	}
	findings.SHA = sha

	deps, err := parseSyft(project + "syft.json")
	if err != nil {
		return StatusInternalServerError(err), nil
	}
	findings.Dependencies = deps.Artifacts
	vulns, err := parseCycloneDX(project + "vuln.json")
	if err != nil {
		return StatusInternalServerError(err), nil
	}
	findings.Vulnerabilities = vulns.Vulnerabilities
	cleanedVulns := Clean(findings.Vulnerabilities)
	sorted := SortVulns(cleanedVulns)
	matched := MatchDepsToVulns(findings.Dependencies, sorted)
	findings.Vulnerabilities = matched
	repoName := fmt.Sprintf("%s/%s", org, repo)
	r := Vulnerabilities(repoName, findings)
	var resultHTML string
	if len(findings.Vulnerabilities) == 0 {
		resultHTML = "<a>No vulnerabilities found</a>"
	} else {
		resultHTML, err = componentToHTML(r)
		if err != nil {
			return StatusInternalServerError(err), nil
		}
	}

	return events.LambdaFunctionURLResponse{
		StatusCode: 200,
		Body:       resultHTML,
		Headers: map[string]string{
			"Content-Type": "text/plain",
		},
	}, nil
}

func AddRepo(ctx context.Context, request events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	org, repo, err := ParseGithubRepoURL(request.Body)
	if err != nil {
		return StatusInternalServerError(err), nil
	}
	_, err = os.Stat(fmt.Sprintf("db/%s/%s", org, repo))
	if err != nil && !os.IsNotExist(err) {
		return StatusInternalServerError(err), nil
	}
	// if !os.IsNotExist(err) {
	// 	return events.LambdaFunctionURLResponse{
	// 		StatusCode: 204,
	// 		// No need to update the list
	// 		Headers: map[string]string{
	// 			"Content-Type": "text/plain",
	// 		},
	// 	}, nil

	// }
	err = os.MkdirAll(fmt.Sprintf("db/%s/%s", org, repo), 0755)
	if err != nil {
		return StatusInternalServerError(err), nil
	}
	r := RepositoryComponent(fmt.Sprintf("%s/%s", org, repo))
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

func Scan(ctx context.Context, request events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
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

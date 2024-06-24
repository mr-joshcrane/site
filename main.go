package site

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

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
	route := request.RawPath
	switch route {
	case "/scan":
		return Scan(ctx, request)
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

func Scan(ctx context.Context, request events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	findings, err := parse()
	if err != nil {
		fmt.Println(err)
		return StatusInternalServerError(err), nil
	}
	vulns := Vulnerabilities(findings)
	page, err := componentToHTML(vulns)
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

func parse() ([]Vulnerability, error) {
	var findings []Vulnerability
	data, err := os.ReadFile("vulns.json")
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &findings)
	if err != nil {
		return nil, err
	}
	return findings, nil
}

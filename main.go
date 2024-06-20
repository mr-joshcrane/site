package site

import (
	"bytes"
	"context"
	"embed"
	"html/template"
	"io/fs"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

//go:embed content/* templates/*
var f embed.FS

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
	// Generate Table of Contents
	var htmlFile []byte
	var err error
	if request.RawPath == "/" {
		htmlFile, err = generateTOC(f)
	} else {
		htmlFile, err = fs.ReadFile(f, request.RawPath[1:])
	}
	if err != nil {
		return StatusInternalServerError(err), nil
	}
	return events.LambdaFunctionURLResponse{
		StatusCode: 200,
		Body:       string(htmlFile),
		Headers: map[string]string{
			"Content-Type": "text/html",
		},
	}, nil
}

type Entry struct {
	Title string
	URL   string
}

// generateTOC creates a list of links for the Table of Contents
func generateTOC(fsys fs.FS) ([]byte, error) {
	var toc []Entry
	err := fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !strings.Contains(path, "post") {
			return nil
		}
		toc = append(toc, Entry{
			Title: d.Name(),
			URL:   path,
		})
		return nil
	})
	baseTemplate, err := fs.ReadFile(fsys, "templates/base_template.html")
	if err != nil {
		return nil, err
	}

	tmpl, err := template.New("base").Parse(string(baseTemplate))
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, toc)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

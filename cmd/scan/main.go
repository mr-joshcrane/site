package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/mr-joshcrane/site"
)

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		fmt.Println("Usage: scan <repository>")
		os.Exit(1)
	}
	s := strings.Split(args[0], "/")
	if len(s) != 2 {
		fmt.Println("Invalid repository format. Use <org>/<repo>")
		os.Exit(1)
	}
	org := s[0]
	repo := s[1]

	err := site.GetVulnerabilityData(org, repo)
	if err != nil {
		fmt.Println(err)
	}
}

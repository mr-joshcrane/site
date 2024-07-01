package site

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func GetVulnerabilityData(org string, repo string) error {
	// get repository list

	err := gitClone(org, repo)
	if err != nil {
		return err
	}
	return syftAndGrype(org, repo)
}

func gitClone(org string, repo string) error {
	// git clone
	combined := org + "/" + repo
	location := "./db/" + org + "/" + repo + "/"
	err := os.MkdirAll(location, 0755)
	if err != nil {
		return err
	}
	// clone
	cmd := exec.Command("gh", "repo", "clone", combined, location)
	cmd.Run()

	// get commit sha
	sha := new(bytes.Buffer)
	cmd = exec.Command("git", "rev-parse", "HEAD")
	cmd.Stdout = sha
	err = cmd.Run()
	if err != nil {
		return err
	}
	return os.WriteFile(location+"sha", sha.Bytes(), 0644)
}

func syftAndGrype(org string, repo string) error {
	dbDir := fmt.Sprintf("db/%s/%s/.", org, repo)
	cmd := exec.Command("syft", "scan", dbDir, "-o", "json")
	data := new(bytes.Buffer)
	cmd.Stdout = data
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return err
	}
	err = os.WriteFile(dbDir+"/syft.json", data.Bytes(), 0644)
	if err != nil {
		return err
	}
	cmd = exec.Command("grype", "-o", "cyclonedx-json")
	cmd.Stdin = data
	vulns := new(bytes.Buffer)
	cmd.Stdout = vulns
	err = cmd.Run()
	if err != nil {
		return err
	}
	return os.WriteFile(dbDir+"/vuln.json", vulns.Bytes(), 0644)
}

func ParseGithubRepoURL(rawURL string) (string, string, error) {
	s := strings.Split(rawURL, "/")
	if len(s) == 2 {
		return s[0], s[1], nil
	}
	// repoName=cultureamp%2Fmurmur
	spl := strings.Split(rawURL, "=")
	if len(spl) != 2 {
		return "", "", fmt.Errorf("Invalid URL")
	}
	spl = strings.Split(spl[1], "%2F")
	if len(spl) != 2 {
		return "", "", fmt.Errorf("Invalid URL")
	}
	return spl[0], spl[1], nil
}

package site

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

func GetVulnerabilityData(org string, repo string) error {
	// get repository list
	currentDir, err  := os.Getwd()
	if err != nil {
		return err
	}
	fmt.Println(currentDir)

	tempDir, err := os.MkdirTemp("", "tempdir")
	if err != nil {
		return err
	}
	path, err := gitClone(org, repo, tempDir)
	if err != nil {
		return err
	}

	vulns, err := syftAndGrype(path)
	if err != nil {
		return err
	}
	fmt.Println(vulns)
	err = os.WriteFile(currentDir + "/vulns.json", vulns.Bytes(), 0644)
	if err != nil {
		return err
	}
	return nil
}

func gitClone(org string, repo string, tempDirPath string) (string, error) {
	// git clone
	err := os.Chdir(tempDirPath)
	if err != nil {
		return "", err
	}
	combined := org + "/" + repo
	location := tempDirPath + "/" + repo + "/"
	err = os.Mkdir(location, 0755)
	if err != nil {
		return "", err
	}

	cmd := exec.Command("gh", "repo", "clone", combined, location)
	cmd.Start()
	err = cmd.Wait()
	if err != nil {
		return "", err
	}

	return location, nil
}

func syftAndGrype(path string) (*bytes.Buffer, error) {

	err := os.Chdir(path)
	if err != nil {
		return nil, err
	}
	// syft scan . -o json

	cmd := exec.Command("syft", "scan", path, "-o", "json")
	data := new(bytes.Buffer)
	cmd.Stdout = data
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return nil, err
	}
	cmd = exec.Command("grype", "-o", "cyclonedx-json")
	cmd.Stdin = data
	vulns := new(bytes.Buffer)
	cmd.Stdout = vulns
	err = cmd.Run()
	if err != nil {
		return nil, err
	}
	fmt.Println(vulns.String())
	return vulns, nil
}

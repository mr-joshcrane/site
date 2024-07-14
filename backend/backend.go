package backend

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/anchore/grype/grype"
	"github.com/anchore/grype/grype/db"
	"github.com/anchore/grype/grype/pkg"
	"github.com/anchore/stereoscope/pkg/image"
	"github.com/anchore/syft/syft"
	"github.com/anchore/syft/syft/format"
	"github.com/anchore/syft/syft/format/cyclonedxjson"
	"github.com/anchore/syft/syft/sbom"
	"github.com/mr-joshcrane/site/store"
)

func Workers(s store.Store) error {

	err := ScanRepos(s)
	if err != nil {
		return err
	}
	repos, err := s.ListRepos()
	if err != nil {
		return err
	}

	for _, repo := range repos {
		fmt.Println("Processing", repo.Org, repo.Name)
		file, err := GitCloneRepo(repo.Org, repo.Name, repo.HEAD)
		if err != nil {
			return err
		}
		syftDoc, err := Syft(file, fmt.Sprintf("%s/%s/%s", repo.Org, repo.Name, repo.HEAD))
		if err != nil {
			fmt.Println("Error syfting")
			return err
		}
		err = s.StoreSyft(repo.Org, repo.Name, repo.HEAD, *syftDoc)
		if err != nil {
			fmt.Println("Error storing syft")
			return err
		}
		encoder, err := cyclonedxjson.NewFormatEncoderWithConfig(cyclonedxjson.EncoderConfig{
			Version: "1.6",
		})
		if err != nil {
			return err
		}
		d, err := format.Encode(*syftDoc, encoder)
		if err != nil {
			return err
		}
		filename := fmt.Sprintf("sbom_for%s.json", repo.Name)
		err = os.WriteFile(filename, d, 0644)
		if err != nil {
			return err
		}
		vulnDB, _, closer, err := grype.LoadVulnerabilityDB(db.Config{
			DBRootDir:           "~/.cache/grype/db",
			ListingURL:          "https://toolbox-data.anchore.io/grype/databases/listing.json",
			ValidateByHashOnGet: false,
		}, true)
		if err != nil {
			return err
		}
		defer closer.Close()

		matcher := grype.DefaultVulnerabilityMatcher(*vulnDB)
		p, c, _, err := pkg.Provide(fmt.Sprintf("file:%s", filename), pkg.ProviderConfig{
			SyftProviderConfig: pkg.SyftProviderConfig{
				SBOMOptions: syft.DefaultCreateSBOMConfig(),
				RegistryOptions: &image.RegistryOptions{
					Credentials: []image.RegistryCredentials{},
				}}})
		if err != nil {
			return err
		}
		matches, noMATCH, err := matcher.FindMatches(p, c)
		if err != nil {
			return err
		}
		fmt.Println(noMATCH)
		fmt.Println(matches.Count())

	}
	return nil
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
	}
	return nil
}

func GitCloneRepo(org string, repo string, commit string) (*bytes.Buffer, error) {
	// git clone
	zip, err := GetArchiveZip(NewGithubClient(""), org, repo, commit)
	if err != nil {
		return nil, err
	}
	return zip, nil
}

func Syft(archive *bytes.Buffer, name string) (*sbom.SBOM, error) {
	// file is a zip
	ctx := context.Background()

	tempDir, err := os.MkdirTemp("", "syft")
	if err != nil {
		return nil, err

	}
	defer os.RemoveAll(tempDir)
	outputPath := filepath.Join(tempDir, "archive.zip")
	err = os.WriteFile(outputPath, archive.Bytes(), 0644)
	if err != nil {
		return nil, err
	}
	source, err := syft.GetSource(ctx, outputPath, nil)
	if err != nil {
		return nil, err
	}
	sbom, err := syft.CreateSBOM(ctx, source, nil)
	if err != nil {
		return nil, err
	}
	return sbom, nil
}

func originSplit(origin string) (string, string, string, error) {
	s := strings.Split(origin, "/")
	if len(s) != 3 {
		return "", "", "", fmt.Errorf("origin must be in the format org/repo/commit")
	}
	return s[0], s[1], s[2], nil
}

// func Grype(syftDocument store.SyftDocument) error {
// 	fmt.Println("Gryping", syftDocument.Organisation, syftDocument.Repository, syftDocument.Commit)
// 	rawSyft, err := json.Marshal(syftDocument)
// 	if err != nil {
// 		return err
// 	}
// 	buf := new(bytes.Buffer)
// 	cmd := exec.Command("grype", "--output", "cyclonedx-json")
// 	cmd.Stderr = os.Stderr
// 	cmd.Stdout = buf
// 	cmd.Stdin = bytes.NewReader(rawSyft)
// 	err = cmd.Run()
// 	if err != nil {
// 		return err
// 	}
// 	fmt.Println(buf.String())
// 	return nil
// }

// func ParseGrype(data io.Reader) (store.GrypeDocument, error) {
// 	var document store.GrypeDocument
// 	d, err := io.ReadAll(data)
// 	if err != nil {
// 		return store.GrypeDocument{}, err
// 	}
// 	err = json.Unmarshal(d, &document)
// 	if err != nil {
// 		return store.GrypeDocument{}, err
// 	}
// 	return document, nil
// }

package site

import (
	"fmt"
	"sort"
	"strings"
)

type RepositoryFindings struct {
	SHA             string          `json:"sha"`
	Organisation    string          `json:"organisation"`
	Repository      string          `json:"repository"`
	Dependencies    []Dependency    `json:"dependencies"`
	Vulnerabilities []Vulnerability `json:"vulnerabilities"`
}

type Location struct {
	Path string `json:"path"`
}

type Dependency struct {
	Name       string     `json:"name"`
	Version    string     `json:"version"`
	Type       string     `json:"type"`
	Locations  []Location `json:"locations"`
	PackageURL string     `json:"purl"`
	Language   string     `json:"language"`
	Metadata   struct {
		Resolved string `json:"resolved"`
	} `json:"metadata"`
}

func (d Dependency) FileLocations() string {
	for _, loc := range d.Locations {
		return loc.Path
	}
	return ""
}

func MatchDepsToVulns(deps []Dependency, vuln []Vulnerability) []Vulnerability {
	var matchedVulns []Vulnerability
	for _, vuln := range vuln {
		var augmentedDep Dependency
		for _, dep := range deps {
			if len(vuln.Affects) != 1 {
				fmt.Println("VIOLATION: Affects length")
				fmt.Println(vuln.ID)
				fmt.Println(vuln.Affects)
			}
			if strings.HasPrefix(vuln.Affects[0].Ref, dep.PackageURL) {
				augmentedDep = dep
				break
			}
		}
		vuln.Dependency = augmentedDep
		matchedVulns = append(matchedVulns, vuln)
	}
	return matchedVulns
}

type SyftFinding struct {
	Artifacts []Dependency `json:"artifacts"`
}

// Vulnerability holds the data structure for the vulnerability information.
type Vulnerability struct {
	BomRef      string      `json:"bom-ref"`
	PackageURL  string      `json:"purl"`
	ID          string      `json:"id"`
	Source      Source      `json:"source"`
	References  []Reference `json:"references"`
	Ratings     []Rating    `json:"ratings"`
	Description string      `json:"description"`
	Advisories  []Advisory  `json:"advisories"`
	Affects     []Affect    `json:"affects"`
	Dependency  Dependency  `json:"dependency,omitempty"`
}

type AllVulnerabilities []Vulnerability

// CycloneDX holds the data structure for the CycloneDX BOM.
type CycloneDX struct {
	Vulnerabilities AllVulnerabilities `json:"vulnerabilities"`
}

func SortVulns(vulns []Vulnerability) []Vulnerability {
	sortVulnerabilitiesByRatingScore(vulns, func(v1, v2 *Vulnerability) bool {
		return v1.Ratings[0].Score > v2.Ratings[0].Score
	})
	return vulns
}

// SortVulnerabilitiesByRatingScore sorts the vulnerabilities by the rating score using the provided comparison function.
func sortVulnerabilitiesByRatingScore(vulnerabilities []Vulnerability, lessFunc func(v1, v2 *Vulnerability) bool) {
	sort.Slice(vulnerabilities, func(i, j int) bool {
		return lessFunc(&vulnerabilities[i], &vulnerabilities[j])
	})
}

// Source holds the source information.
type Source struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// Reference holds the reference information.
type Reference struct {
	ID     string `json:"id"`
	Source Source `json:"source"`
}

// Rating holds the rating information.
type Rating struct {
	Severity string  `json:"severity"`
	Score    float64 `json:"score,omitempty"`
	Method   string  `json:"method,omitempty"`
	Vector   string  `json:"vector,omitempty"`
}

// Advisory holds the advisory information.
type Advisory struct {
	URL string `json:"url"`
}

// Affect holds the affected package information.
type Affect struct {
	Ref string `json:"ref"`
}

func Split(s string) string {
	s = strings.Split(s, ":")[1]
	return strings.Split(s, "?")[0]
}

func Clean(v AllVulnerabilities) AllVulnerabilities {
	var cleanedVulns AllVulnerabilities
	for _, vuln := range v {
		for i, r := range vuln.Ratings {
			var adjustedScore float64
			var adjustedSeverity string
			switch {
			case r.Severity != "" && r.Score > 0:
			case r.Score > 0 && r.Severity == "":
				switch {
				case r.Score >= 9.0:
					adjustedSeverity = "critical"
				case r.Score >= 7.0:
					adjustedSeverity = "high"
				case r.Score >= 4.0:
					adjustedSeverity = "medium"
				case r.Score >= 0.1:
					adjustedSeverity = "low"
				default:
					adjustedSeverity = "unknown"
				}
			case r.Score <= 0 && r.Severity != "":
				switch strings.ToLower(r.Severity) {
				case "critical":
					adjustedScore = 9.0
				case "high":
					adjustedScore = 7.0
				case "medium":
					adjustedScore = 4.0
				case "low":
					adjustedScore = 1.0
				default:
					adjustedScore = 0.0
				}
			case r.Severity == "" && r.Score == 0.0:
				adjustedSeverity = "Unknown"
				adjustedScore = 0
			}

			if r.Severity == "" {
				vuln.Ratings[i].Severity = adjustedSeverity
			}
			if r.Score <= 0.0 {
				vuln.Ratings[i].Score = adjustedScore
			}
			cleanedVulns = append(cleanedVulns, vuln)
		}

	}
	return cleanedVulns
}

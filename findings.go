package site

import "strings"

// Vulnerability holds the data structure for the vulnerability information.
type Vulnerability struct {
	BomRef      string      `json:"bom-ref"`
	ID          string      `json:"id"`
	Source      Source      `json:"source"`
	References  []Reference `json:"references"`
	Ratings     []Rating    `json:"ratings"`
	Description string      `json:"description"`
	Advisories  []Advisory  `json:"advisories"`
	Affects     []Affect    `json:"affects"`
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

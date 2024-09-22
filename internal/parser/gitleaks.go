package parser

import (
	"encoding/json"
	"fmt"
)

type GitLeaksReport []GitLeaksFinding

type GitLeaksFinding struct {
	Description string        `json:"Description"`
	StartLine   int           `json:"StartLine"`
	EndLine     int           `json:"EndLine"`
	StartColumn int           `json:"StartColumn"`
	EndColumn   int           `json:"EndColumn"`
	Match       string        `json:"Match"`
	Secret      string        `json:"Secret"`
	File        string        `json:"File"`
	SymlinkFile string        `json:"SymlinkFile"`
	Commit      string        `json:"Commit"`
	Entropy     float64       `json:"Entropy"`
	Author      string        `json:"Author"`
	Email       string        `json:"Email"`
	Date        string        `json:"Date"`
	Message     string        `json:"Message"`
	Tags        []interface{} `json:"Tags"`
	RuleID      string        `json:"RuleID"`
	Fingerprint string        `json:"Fingerprint"`
}

func (finding GitLeaksFinding) GetFullValue() string {
	if len(finding.Match) < 500 {
		return finding.Match
	}
	return finding.Secret
}

type GitLeaksParser struct {
	prefixesToRemove []string
}

func NewGitLeaksParser(prefixesToRemove []string) Parser {
	return &GitLeaksParser{prefixesToRemove: prefixesToRemove}
}

func (p *GitLeaksParser) Parse(report []byte, taskId string) ([]Finding, error) {
	var gitLeaksReport GitLeaksReport
	err := json.Unmarshal(report, &gitLeaksReport)
	if err != nil {
		return nil, err
	}

	newPrefixesToRemove := make([]string, len(p.prefixesToRemove))
	copy(newPrefixesToRemove, p.prefixesToRemove)
	newPrefixesToRemove = append(newPrefixesToRemove, fmt.Sprintf("/%s", taskId))

	var findings []Finding

	for _, finding := range gitLeaksReport {
		findings = append(findings, Finding{
			Line:      finding.StartLine,
			Value:     finding.Secret,
			FullValue: finding.GetFullValue(),
			File:      removePrefixes(finding.File, newPrefixesToRemove),
			Evidence:  Medium,
			Commit:    "no value available",
			Author:    "no value available",
			Email:     "no value available",
			Date:      "no value available",
		})
	}

	return findings, nil
}

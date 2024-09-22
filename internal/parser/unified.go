package parser

import "encoding/json"

type UnifiedReport struct {
	Low    []Finding `json:"LOW"`
	Medium []Finding `json:"MEDIUM"`
	High   []Finding `json:"HIGH"`
}

type AptbReport struct {
	TaskId       string         `json:"taskUUID"`
	Practice     string         `json:"practice"`
	PracticeTool string         `json:"practiceTool"`
	ScanResult   *UnifiedReport `json:"scanResult"`
}

type Finding struct {
	Date      string   `json:"Date"`
	File      string   `json:"File"`
	Line      int      `json:"Line"`
	Email     string   `json:"Email"`
	Value     string   `json:"Value"`
	Author    string   `json:"Author"`
	Commit    string   `json:"Commit"`
	Evidence  Evidence `json:"Evidence"`
	FullValue string   `json:"Full_value"`
}

type Evidence int

const (
	Low Evidence = iota
	Medium
	High
)

func (e Evidence) String() string {
	switch e {
	case Low:
		return "Low"
	case Medium:
		return "Medium"
	case High:
		return "High"
	default:
		return "Unknown"
	}
}

func (e Evidence) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.String())
}

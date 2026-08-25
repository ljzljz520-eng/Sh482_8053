package ingest

import "strings"

type Report struct {
	Total     int         `json:"total"`
	Succeeded int         `json:"succeeded"`
	Failed    int         `json:"failed"`
	Errors    []string    `json:"errors"`
	Results   []RowResult `json:"results"`
}

func (r Report) Success() bool { return r.Failed == 0 && len(r.Errors) == 0 }

func (r Report) Message() string {
	if r.Success() {
		return "import completed"
	}
	return strings.Join(r.Errors, "; ")
}

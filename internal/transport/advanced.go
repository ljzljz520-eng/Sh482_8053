package transport

import (
	"net/http"
	"strconv"

	"enterpriselead/internal/domain"
)

func (s *Server) insights(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, errMethod)
		return
	}
	query := domain.SearchQuery{Text: r.URL.Query().Get("q"), Company: r.URL.Query().Get("company"), Status: domain.LeadStatus(r.URL.Query().Get("status")), Priority: domain.Priority(r.URL.Query().Get("priority")), Tag: r.URL.Query().Get("tag")}
	query.Limit, _ = strconv.Atoi(r.URL.Query().Get("limit"))
	result, err := s.app.Insights(query)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (s *Server) exportRecords(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, errMethod)
		return
	}
	query := domain.SearchQuery{Text: r.URL.Query().Get("q"), Company: r.URL.Query().Get("company"), Status: domain.LeadStatus(r.URL.Query().Get("status")), Priority: domain.Priority(r.URL.Query().Get("priority")), Tag: r.URL.Query().Get("tag"), IncludeArchived: true}
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment; filename=leadboard.csv")
	if err := s.app.Export(w, query); err != nil {
		writeError(w, http.StatusInternalServerError, err)
	}
}

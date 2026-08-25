package transport

import (
	"encoding/json"
	"net/http"
	"strings"

	"enterpriselead/internal/service"
)

type Server struct {
	app *service.Service
	mux *http.ServeMux
}

func New(app *service.Service) *Server {
	s := &Server{app: app, mux: http.NewServeMux()}
	s.mux.HandleFunc("/health", s.health)
	s.mux.HandleFunc("/records", s.records)
	s.mux.HandleFunc("/records/", s.record)
	s.mux.HandleFunc("/import", s.importRecords)
	s.mux.HandleFunc("/insights", s.insights)
	s.mux.HandleFunc("/export", s.exportRecords)
	return s
}

func (s *Server) Handler() http.Handler { return logging(s.mux) }

func (s *Server) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Leadboard-Service", "enterprise-leadboard")
		next.ServeHTTP(w, r)
	})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, map[string]string{"error": err.Error()})
}

func decodeJSON(r *http.Request, target any) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(target)
}

func (s *Server) records(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		s.createRecord(w, r)
		return
	}
	if r.Method == http.MethodGet {
		s.searchRecords(w, r)
		return
	}
	w.Header().Set("Allow", "GET, POST")
	writeError(w, http.StatusMethodNotAllowed, errMethod)
}

var errMethod = methodError("method not allowed")

type methodError string

func (e methodError) Error() string { return string(e) }

func parseID(path string) (string, string) {
	trimmed := strings.TrimPrefix(path, "/records/")
	parts := strings.Split(strings.Trim(trimmed, "/"), "/")
	if len(parts) == 0 || parts[0] == "" {
		return "", ""
	}
	if len(parts) == 1 {
		return parts[0], ""
	}
	return parts[0], parts[1]
}

func (s *Server) record(w http.ResponseWriter, r *http.Request) {
	id, action := parseID(r.URL.Path)
	if id == "" {
		writeError(w, http.StatusNotFound, methodError("record id is required"))
		return
	}
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, errMethod)
		return
	}
	switch action {
	case "review":
		s.reviewRecord(w, r, id)
	case "approve":
		s.approveRecord(w, r, id)
	case "archive":
		s.archiveRecord(w, r, id)
	case "update":
		s.updateRecord(w, r, id)
	default:
		writeError(w, http.StatusNotFound, methodError("unknown record action"))
	}
}

package transport

import (
	"net/http"
	"strconv"
	"strings"

	"enterpriselead/internal/domain"
	"enterpriselead/internal/ingest"
	"enterpriselead/internal/service"
)

type createRequest struct {
	service.CreateInput
	Actor string `json:"actor"`
}

func (s *Server) createRecord(w http.ResponseWriter, r *http.Request) {
	var request createRequest
	if err := decodeJSON(r, &request); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	record, err := s.app.Create(request.CreateInput, request.Actor)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err)
		return
	}
	writeJSON(w, http.StatusCreated, record)
}

func (s *Server) searchRecords(w http.ResponseWriter, r *http.Request) {
	query := domain.SearchQuery{Text: r.URL.Query().Get("q"), Company: r.URL.Query().Get("company"), Owner: r.URL.Query().Get("owner"), Status: domain.LeadStatus(r.URL.Query().Get("status")), Priority: domain.Priority(r.URL.Query().Get("priority")), Tag: r.URL.Query().Get("tag")}
	query.IncludeArchived, _ = strconv.ParseBool(r.URL.Query().Get("include_archived"))
	query.Limit, _ = strconv.Atoi(r.URL.Query().Get("limit"))
	query.Offset, _ = strconv.Atoi(r.URL.Query().Get("offset"))
	result, err := s.app.Search(query)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, result)
}

type actorRequest struct {
	Actor string `json:"actor"`
	Note  string `json:"note"`
}

func readActor(r *http.Request) (actorRequest, error) {
	var request actorRequest
	if err := decodeJSON(r, &request); err != nil {
		return request, err
	}
	request.Actor = strings.TrimSpace(request.Actor)
	if request.Actor == "" {
		request.Actor = "system"
	}
	return request, nil
}

func (s *Server) reviewRecord(w http.ResponseWriter, r *http.Request, id string) {
	request, err := readActor(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	record, workflow, err := s.app.Review(id, request.Actor)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"record": record, "workflow": workflow})
}

func (s *Server) approveRecord(w http.ResponseWriter, r *http.Request, id string) {
	request, err := readActor(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	record, workflow, err := s.app.Approve(id, request.Actor)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"record": record, "workflow": workflow})
}

func (s *Server) archiveRecord(w http.ResponseWriter, r *http.Request, id string) {
	request, err := readActor(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	record, workflow, err := s.app.Archive(id, request.Actor)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"record": record, "workflow": workflow})
}

type updateRequest struct {
	domain.Change
	Version int    `json:"version"`
	Actor   string `json:"actor"`
}

func (s *Server) updateRecord(w http.ResponseWriter, r *http.Request, id string) {
	var request updateRequest
	if err := decodeJSON(r, &request); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	record, err := s.app.Update(id, request.Version, request.Change, request.Actor)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err)
		return
	}
	writeJSON(w, http.StatusOK, record)
}

func (s *Server) importRecords(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, errMethod)
		return
	}
	actor := r.URL.Query().Get("actor")
	if actor == "" {
		actor = "import"
	}
	report, err := ingest.New(s.app).Import(r.Body, actor)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	status := http.StatusOK
	if !report.Success() {
		status = http.StatusMultiStatus
	}
	writeJSON(w, status, report)
}

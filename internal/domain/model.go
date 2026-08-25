package domain

import "time"

type LeadStatus string

const (
	StatusDraft    LeadStatus = "draft"
	StatusReview   LeadStatus = "review"
	StatusApproved LeadStatus = "approved"
	StatusArchived LeadStatus = "archived"
	StatusRejected LeadStatus = "rejected"
)

type Priority string

const (
	PriorityLow    Priority = "low"
	PriorityNormal Priority = "normal"
	PriorityHigh   Priority = "high"
	PriorityUrgent Priority = "urgent"
)

type Record struct {
	ID             string     `json:"id"`
	Company        string     `json:"company"`
	ContactName    string     `json:"contact_name"`
	ContactEmail   string     `json:"contact_email"`
	Source         string     `json:"source"`
	Need           string     `json:"need"`
	Owner          string     `json:"owner"`
	Status         LeadStatus `json:"status"`
	Priority       Priority   `json:"priority"`
	Summary        string     `json:"summary"`
	Tags           []string   `json:"tags"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	ArchivedAt     *time.Time `json:"archived_at,omitempty"`
	Version        int        `json:"version"`
	LastWorkflowID string     `json:"last_workflow_id"`
}

type AuditEvent struct {
	ID         string            `json:"id"`
	RecordID   string            `json:"record_id"`
	Action     string            `json:"action"`
	Actor      string            `json:"actor"`
	FromStatus LeadStatus        `json:"from_status"`
	ToStatus   LeadStatus        `json:"to_status"`
	Note       string            `json:"note"`
	At         time.Time         `json:"at"`
	Metadata   map[string]string `json:"metadata,omitempty"`
}

type Workflow struct {
	ID          string     `json:"id"`
	RecordID    string     `json:"record_id"`
	Kind        string     `json:"kind"`
	RequestedBy string     `json:"requested_by"`
	ApprovedBy  string     `json:"approved_by"`
	State       string     `json:"state"`
	Steps       []string   `json:"steps"`
	CreatedAt   time.Time  `json:"created_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}

type Attachment struct {
	ID        string    `json:"id"`
	RecordID  string    `json:"record_id"`
	Name      string    `json:"name"`
	MediaType string    `json:"media_type"`
	Size      int64     `json:"size"`
	Checksum  string    `json:"checksum"`
	CreatedAt time.Time `json:"created_at"`
	Content   []byte    `json:"content,omitempty"`
}

type Change struct {
	Company      *string
	ContactName  *string
	ContactEmail *string
	Source       *string
	Need         *string
	Owner        *string
	Priority     *Priority
	Tags         []string
	Summary      *string
}

type SearchQuery struct {
	Text            string
	Company         string
	Owner           string
	Status          LeadStatus
	Priority        Priority
	Tag             string
	IncludeArchived bool
	Limit           int
	Offset          int
}

type SearchResult struct {
	Records []Record `json:"records"`
	Total   int      `json:"total"`
	Offset  int      `json:"offset"`
	Limit   int      `json:"limit"`
}

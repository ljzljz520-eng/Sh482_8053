package policy

import (
	"strings"

	"enterpriselead/internal/domain"
)

type Role string

const (
	RoleOperator Role = "operator"
	RoleReviewer Role = "reviewer"
	RoleManager  Role = "manager"
	RoleImporter Role = "importer"
)

func ParseRole(value string) Role {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "reviewer":
		return RoleReviewer
	case "manager":
		return RoleManager
	case "importer":
		return RoleImporter
	default:
		return RoleOperator
	}
}

func CanPerform(role Role, action string, record domain.Record) bool {
	switch action {
	case "create":
		return role == RoleOperator || role == RoleImporter || role == RoleManager
	case "review":
		return role == RoleReviewer || role == RoleManager
	case "approve":
		return role == RoleReviewer || role == RoleManager
	case "archive":
		return role == RoleOperator || role == RoleManager
	case "update":
		return role == RoleOperator || role == RoleManager || (role == RoleReviewer && record.Status == domain.StatusReview)
	default:
		return false
	}
}

func AllowedActors(action string) []Role {
	roles := make([]Role, 0, 3)
	for _, role := range []Role{RoleOperator, RoleReviewer, RoleManager, RoleImporter} {
		if CanPerform(role, action, domain.Record{Status: domain.StatusReview}) {
			roles = append(roles, role)
		}
	}
	return roles
}

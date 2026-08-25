package domain

import "fmt"

func CanTransition(from, to LeadStatus) bool {
	if from == StatusDraft && to == StatusReview {
		return true
	}
	if from == StatusReview && (to == StatusApproved || to == StatusRejected) {
		return true
	}
	if from == StatusApproved && to == StatusArchived {
		return true
	}
	if from == StatusRejected && to == StatusDraft {
		return true
	}
	return false
}

func Transition(r Record, to LeadStatus) (Record, error) {
	if !CanTransition(r.Status, to) {
		return r, fmt.Errorf("cannot transition %s to %s", r.Status, to)
	}
	r.Status = to
	return r, nil
}

func IsTerminal(status LeadStatus) bool {
	return status == StatusArchived
}

func IsActive(status LeadStatus) bool {
	return status != StatusArchived && status != StatusRejected
}

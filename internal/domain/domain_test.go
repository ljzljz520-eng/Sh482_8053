package domain

import "testing"

func validRecord() Record {
	return Record{ID: "r1", Company: "Acme", ContactName: "Lin", ContactEmail: "lin@acme.test", Need: "procurement portal", Owner: "Ming", Status: StatusDraft, Priority: PriorityHigh, Tags: []string{"saas"}}
}

func TestValidateAndNormalize(t *testing.T) {
	record := NormalizeRecord(Record{Company: " Acme ", ContactName: " Lin ", ContactEmail: "LIN@ACME.TEST", Need: " portal ", Owner: " Ming ", Status: StatusDraft, Priority: PriorityNormal, Tags: []string{" SaaS "}})
	if err := ValidateRecord(record); err != nil {
		t.Fatal(err)
	}
	record.ContactEmail = "bad"
	if err := ValidateRecord(record); err == nil {
		t.Fatal("expected invalid email validation")
	}
	record.ContactEmail = "lin@acme.test"
	if record.Company != "Acme" || record.Tags[0] != "saas" {
		t.Fatalf("unexpected normalization: %#v", record)
	}
}

func TestTransitions(t *testing.T) {
	record := validRecord()
	for _, status := range []LeadStatus{StatusReview, StatusApproved, StatusArchived} {
		changed, err := Transition(record, status)
		if err != nil && status == StatusReview {
			t.Fatal(err)
		}
		record = changed
	}
	if !IsTerminal(record.Status) || IsActive(record.Status) {
		t.Fatalf("unexpected terminal state %s", record.Status)
	}
}

func TestSummaryAndMatch(t *testing.T) {
	record := validRecord()
	record.Summary = BuildSummary(record)
	if !MatchText(record, "portal") || !HasTag(record, "SAAS") {
		t.Fatal("expected record to match")
	}
	if MatchText(record, "missing") {
		t.Fatal("unexpected match")
	}
}

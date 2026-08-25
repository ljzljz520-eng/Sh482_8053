package ingest

import (
	"strings"
	"testing"
	"time"

	"enterpriselead/internal/domain"
	"enterpriselead/internal/service"
	"enterpriselead/internal/storage"
)

func TestParseDeterministicRows(t *testing.T) {
	rows, errors, err := Parse(strings.NewReader("company,contact,email,source,need,owner,priority,tags\nA,B,b@a.test,x,y,z,low,a|b\nshort\n"))
	if err != nil || len(rows) != 1 || len(errors) != 1 {
		t.Fatalf("unexpected parse result: %v %#v %#v", err, rows, errors)
	}
	if rows[0].Tags[1] != "b" || rows[0].Priority != domain.PriorityLow {
		t.Fatalf("unexpected row %#v", rows[0])
	}
}

func TestReportMessage(t *testing.T) {
	if (Report{Total: 1, Succeeded: 1}).Message() != "import completed" {
		t.Fatal("expected success message")
	}
	if (Report{Errors: []string{"bad"}}).Message() != "bad" {
		t.Fatal("expected error message")
	}
	store, err := storage.Open(t.TempDir() + "/ingest.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	_ = service.New(store, service.FixedTime{Value: time.Unix(1, 0)})
}

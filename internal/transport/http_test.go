package transport

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"enterpriselead/internal/service"
	"enterpriselead/internal/storage"
)

func TestHTTPCreateAndSearch(t *testing.T) {
	store, err := storage.Open(t.TempDir() + "/http.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	app := service.New(store, service.FixedTime{Value: time.Unix(5, 0)})
	server := httptest.NewServer(New(app).Handler())
	defer server.Close()
	body := []byte(`{"Company":"Acme","ContactName":"Lin","ContactEmail":"lin@acme.test","Need":"contracts","Owner":"Ops","Priority":"normal","Actor":"tester"}`)
	response, err := http.Post(server.URL+"/records", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	if response.StatusCode != http.StatusCreated {
		t.Fatalf("create status %d", response.StatusCode)
	}
	var created map[string]any
	if err := json.NewDecoder(response.Body).Decode(&created); err != nil {
		t.Fatal(err)
	}
	response.Body.Close()
	if created["company"] != "Acme" {
		t.Fatalf("unexpected response %#v", created)
	}
	searchResponse, err := http.Get(server.URL + "/records?q=Acme")
	if err != nil {
		t.Fatal(err)
	}
	defer searchResponse.Body.Close()
	if searchResponse.StatusCode != http.StatusOK {
		t.Fatalf("search status %d", searchResponse.StatusCode)
	}
}

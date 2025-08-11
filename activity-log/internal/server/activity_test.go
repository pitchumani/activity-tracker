package server

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestActivities_InsertAndRetrieve(t *testing.T) {
	acts := &Activities{}
	activity := Activity{Description: "Test activity", Time: time.Now()}
	id := acts.Insert(activity)
	if id != 1 {
		t.Errorf("expected id 1, got %d", id)
	}
	ret, err := acts.Retrieve(id)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ret.Description != activity.Description {
		t.Errorf("expected description %q, got %q", activity.Description, ret.Description)
	}
}

func TestActivities_RetrieveNotFound(t *testing.T) {
	acts := &Activities{}
	_, err := acts.Retrieve(1)
	if err != ErrIDNotFound {
		t.Errorf("expected ErrIDNotFound, got %v", err)
	}
}

func TestHTTPServer_PostAndGet(t *testing.T) {
	srv := &httpServer{Activities: &Activities{}}
	router := http.NewServeMux()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			srv.handlePost(w, r)
		} else if r.Method == http.MethodGet {
			srv.handleGet(w, r)
		}
	})

	ts := httptest.NewServer(router)
	defer ts.Close()

	// POST
	activity := Activity{Description: "Test via HTTP", Time: time.Now()}
	adoc := ActivityDocument{Activity: activity}
	body, _ := json.Marshal(adoc)
	resp, err := http.Post(ts.URL+"/", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("POST error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("POST status: got %d, want %d", resp.StatusCode, http.StatusOK)
	}
	var idDoc IDDocument
	if err := json.NewDecoder(resp.Body).Decode(&idDoc); err != nil {
		t.Fatalf("decode POST response: %v", err)
	}
	if idDoc.ID == 0 {
		t.Errorf("expected nonzero ID")
	}

	// GET
	idBody, _ := json.Marshal(IDDocument{ID: idDoc.ID})
	getReq, _ := http.NewRequest(http.MethodGet, ts.URL+"/", bytes.NewReader(idBody))
	getReq.Header.Set("Content-Type", "application/json")
	getResp, err := http.DefaultClient.Do(getReq)
	if err != nil {
		t.Fatalf("GET error: %v", err)
	}
	defer getResp.Body.Close()
	if getResp.StatusCode != http.StatusOK {
		t.Fatalf("GET status: got %d, want %d", getResp.StatusCode, http.StatusOK)
	}
	var got ActivityDocument
	if err := json.NewDecoder(getResp.Body).Decode(&got); err != nil {
		t.Fatalf("decode GET response: %v", err)
	}
	if got.Activity.Description != activity.Description {
		t.Errorf("expected description %q, got %q", activity.Description, got.Activity.Description)
	}
}

func TestHTTPServer_GetNotFound(t *testing.T) {
	srv := &httpServer{Activities: &Activities{}}
	router := http.NewServeMux()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			srv.handleGet(w, r)
		}
	})

	ts := httptest.NewServer(router)
	defer ts.Close()

	idBody, _ := json.Marshal(IDDocument{ID: 42})
	getReq, _ := http.NewRequest(http.MethodGet, ts.URL+"/", bytes.NewReader(idBody))
	getReq.Header.Set("Content-Type", "application/json")
	getResp, err := http.DefaultClient.Do(getReq)
	if err != nil {
		t.Fatalf("GET error: %v", err)
	}
	defer getResp.Body.Close()
	if getResp.StatusCode != http.StatusNotFound {
		t.Fatalf("GET status: got %d, want %d", getResp.StatusCode, http.StatusNotFound)
	}
	b, _ := io.ReadAll(getResp.Body)
	if string(b) == "" {
		t.Errorf("expected error message in body")
	}
}

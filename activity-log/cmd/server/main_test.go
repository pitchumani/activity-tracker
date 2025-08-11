package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	api "github.com/pitchumani/activity-tracker/activity-log"
	"github.com/pitchumani/activity-tracker/activity-log/internal/server"
)

func TestMain_StartsServer(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	ts := httptest.NewServer(h)
	defer ts.Close()

	resp, err := http.Get(ts.URL)
	if err != nil {
		t.Fatalf("could not GET: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", resp.StatusCode)
	}
}

func TestNewHTTPServer_Integration(t *testing.T) {
	srv := server.NewHTTPServer(":0")
	ts := httptest.NewUnstartedServer(srv.Handler)
	ts.Start()
	defer ts.Close()

	// POST an activity
	activity := api.Activity{Description: "Integration test", Time: time.Now()}
	body, _ := json.Marshal(api.ActivityDocument{Activity: activity})
	resp, err := http.Post(ts.URL+"/", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("POST error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("POST status: got %d, want %d", resp.StatusCode, http.StatusOK)
	}
	var idDoc api.IDDocument
	if err := json.NewDecoder(resp.Body).Decode(&idDoc); err != nil {
		t.Fatalf("decode POST response: %v", err)
	}
	if idDoc.ID == 0 {
		t.Errorf("expected nonzero ID")
	}

	// GET the activity
	idBody, _ := json.Marshal(api.IDDocument{ID: idDoc.ID})
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
	var got api.ActivityDocument
	if err := json.NewDecoder(getResp.Body).Decode(&got); err != nil {
		t.Fatalf("decode GET response: %v", err)
	}
	if got.Activity.Description != activity.Description {
		t.Errorf("expected description %q, got %q", activity.Description, got.Activity.Description)
	}
}

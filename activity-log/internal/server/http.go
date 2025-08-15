package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	api "github.com/pitchumani/activity-tracker/activity-log"
)

// define type httpServer with activities data
// add Methods to that type to handle GET and POST requests

type httpServer struct {
	Activities *Activities
}

func (s *httpServer) handlePost(w http.ResponseWriter, r *http.Request) {
	log.Printf("handlePost")
	var req api.ActivityDocument
	// decode the request body - assuming json of type ActivityDocument
	err := json.NewDecoder(r.Body).Decode(&req)
	// if couldn't decode to desired type, write error into response
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// add the activity to the database
	id, err := s.Activities.Insert(req.Activity)
	if err != nil {
		log.Fatalf("Error: %s", err.Error())
		id = -1
	}
	// create the response - IDDocument
	res := api.IDDocument{ID: id}
	// write the response
	json.NewEncoder(w).Encode(res)
}

func (s *httpServer) handleGet(w http.ResponseWriter, r *http.Request) {
	log.Printf("handleGet")
	var req api.IDDocument
	// decode the request body - assuming json of type IDDocument
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// get the activity info for the id
	activity, err := s.Activities.Retrieve(req.ID)
	if err == ErrIDNotFound {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	// create response with retrieve activity info
	res := api.ActivityDocument{Activity: activity}
	json.NewEncoder(w).Encode(res)
}

func (s *httpServer) handleList(w http.ResponseWriter, r *http.Request) {
	log.Printf("handleList")
	var query api.ActivityQueryDocument
	var err error
	if r.Body != http.NoBody {
		err = json.NewDecoder(r.Body).Decode(&query)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	activities, err := s.Activities.List(query.Offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Returning %d items\n", len(activities))
	err = json.NewEncoder(w).Encode(activities)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// create new HTTP Server, returns *http.Server
func NewHTTPServer(addr string) *http.Server {
	// create instance of httpServer, initialze the activities data
	activities, err := NewActivities()
	if err != nil {
		log.Fatal("Failed to initialize the server: ", err.Error())
		return nil
	}

	server := &httpServer {
		Activities: activities,
	}

	// create the router
	r := mux.NewRouter()
	// assign the get and post handle functions
	r.HandleFunc("/", server.handlePost).Methods("POST")
	r.HandleFunc("/", server.handleGet).Methods("GET")
	r.HandleFunc("/list", server.handleList).Methods("GET")

	// now create actual http server (instance of http.Server)
	// assign handler (router) and return it.
	srv := &http.Server{
		Addr: addr,
		Handler: r,
	}
	return srv
}

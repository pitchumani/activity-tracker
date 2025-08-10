package server

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

// define types for APIs (get and post activity info)
// included the json encoding field in the type itself
type IDDocument struct {
	ID uint64 `json:"id"`
}

type ActivityDocument struct {
	Activity Activity `json:"activity"`
}
	
// define type httpServer with activities data
// add Methods to that type to handle GET and POST requests

type httpServer struct {
	Activities *Activities
}

func (s *httpServer) handlePost(w http.ResponseWriter, r *http.Request) {
	var req ActivityDocument
	// decode the request body - assuming json of type ActivityDocument
	err := json.NewDecoder(r.Body).Decode(&req)
	// if couldn't decode to desired type, write error into response
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// add the activity to the database
	id := s.Activities.Insert(req.Activity)
	// create the response - IDDocument
	res := IDDocument{ID: id}
	// write the response
	json.NewEncoder(w).Encode(res)
}

func (s *httpServer) handleGet(w http.ResponseWriter, r *http.Request) {
	var req IDDocument
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
	res := ActivityDocument{Activity: activity}
	json.NewEncoder(w).Encode(res)
}

// create new HTTP Server, returns *http.Server
func NewHTTPServer(addr string) *http.Server {
	// create instance of httpServer, initialze the activities data
	server := &httpServer {
		Activities: &Activities{},
	}

	// create the router
	r := mux.NewRouter()
	// assign the get and post handle functions
	r.HandleFunc("/", server.handlePost).Methods("POST")
	r.HandleFunc("/", server.handleGet).Methods("GET")

	// now create actual http server (instance of http.Server)
	// assign handler (router) and return it.
	srv := &http.Server{
		Addr: addr,
		Handler: r,
	}
	return srv
}


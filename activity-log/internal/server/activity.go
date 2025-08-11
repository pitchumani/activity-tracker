package server

import (
	"fmt"
	"log"
	"sync"

	api "github.com/pitchumani/activity-tracker/activity-log"
)

// define activity info for the server package

type Activities struct {
	mtx         sync.Mutex
	activities  []api.Activity
}

// define a insert method for Activities type
func (acts *Activities) Insert(activity api.Activity) int {
	acts.mtx.Lock()
	defer acts.mtx.Unlock()
	activity.ID = len(acts.activities) + 1
	acts.activities = append(acts.activities, activity)
	log.Printf("Added %v\n", activity)
	return activity.ID
}

var ErrIDNotFound = fmt.Errorf("ID not found")

// define a retrieve method for Activities type
func (acts *Activities) Retrieve(id int) (api.Activity, error) {
	log.Printf("Getting id %d\n", id)
	acts.mtx.Lock()
	defer acts.mtx.Unlock()
	if id > len(acts.activities) || id == 0 {
		log.Printf("ID %d is not found\n", id)
		return api.Activity{}, ErrIDNotFound
	}
	activity := acts.activities[id-1]
	return activity, nil
}

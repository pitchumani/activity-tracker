package server

import (
	"fmt"

	api "github.com/pitchumani/activity-tracker/activity-log"
)

// define activity info for the server package

type Activities struct {
	activities  []api.Activity
}

// define a insert method for Activities type
func (acts *Activities) Insert(activity api.Activity) int {
	activity.ID = len(acts.activities) + 1
	acts.activities = append(acts.activities, activity)
	return activity.ID
}

var ErrIDNotFound = fmt.Errorf("ID not found")

// define a retrieve method for Activities type
func (acts *Activities) Retrieve(id int) (api.Activity, error) {
	if id > len(acts.activities) || id == 0 {
		return api.Activity{}, ErrIDNotFound
	}
	activity := acts.activities[id-1]
	return activity, nil
}


package server

import (
	"fmt"
	"time"
)

// define activity info for the server package

type Activity struct {
	ID          uint64     `json:"id"`
	Description string     `json:"description"`
	Time        time.Time  `json:"time"`
}

type Activities struct {
	activities  []Activity
}

// define a insert method for Activities type
func (acts *Activities)Insert(activity Activity) uint64 {
	activity.ID = uint64(len(acts.activities)) + 1
	acts.activities = append(acts.activities, activity)
	return activity.ID
}

var ErrIDNotFound = fmt.Errorf("ID not found")

// define a retrieve method for Activities type
func (acts *Activities)Retrieve(id uint64) (Activity, error) {
	if id > uint64(len(acts.activities)) || id == 0 {
		return Activity{}, ErrIDNotFound
	}
	activity := acts.activities[id-1]
	return activity, nil
}


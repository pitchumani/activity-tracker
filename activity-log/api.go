package api

// define the activity info APIs

import "time"

type Activity struct {
	ID          int        `json:"id"`
	Description string     `json:"description"`
	Time        time.Time  `json:"time"`
}

type ActivityDocument struct {
	Activity Activity `json:"activity"`
}

type IDDocument struct {
	ID int `json:"id"`
}

type ActivityQueryDocument struct {
	Offset int `json:"activities"`
}

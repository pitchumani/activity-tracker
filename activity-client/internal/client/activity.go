package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	
	api "github.com/pitchumani/activity-tracker/activity-log"
)

// client needs server URL to make requests
// define type Activities with URL as string member
type Activities struct {
	URL string
}

// define Methods to Activities
// insert and retrieve

func (c *Activities) Insert(activity api.Activity) (int, error) {
	activityDoc := api.ActivityDocument{Activity: activity}
	jsBytes, err := json.Marshal(activityDoc)
	if err != nil {
		return 0, err
	}

	reader := bytes.NewReader(jsBytes)
	// create new request to send data to server 
	req, err := http.NewRequest(http.MethodPost, c.URL, reader)
	if err != nil {
		return 0, err
	}
	// send the request and get response
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}

	if res.Body != nil {
		defer res.Body.Close()
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return 0, err
	}
	// unmarshal the body to IDDocument
	var document api.IDDocument
	err = json.Unmarshal(body, &document)
	if err != nil {
		return 0, err
	}
	return document.ID, nil
}

func (c *Activities) Retrieve(id int) (api.Activity, error) {
	var document api.ActivityDocument
	idDoc := api.IDDocument{ID: id}
	jsBytes, err := json.Marshal(idDoc)
	if err != nil {
		return document.Activity, err
	}
	reader := bytes.NewReader(jsBytes)
	req, err := http.NewRequest(http.MethodGet, c.URL, reader)
	if err != nil {
		return document.Activity, err
	}
	// send the request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return document.Activity, err
	}
	if res.StatusCode == 404 {
		return document.Activity, errors.New("Not Found")
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return document.Activity, err
	}
	err = json.Unmarshal(body, &document)
	if err != nil {
		return document.Activity, err
	}
	return document.Activity, nil
}

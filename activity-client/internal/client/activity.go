package client

import (
	"context"
	"fmt"
	"log"

	api "github.com/pitchumani/activity-tracker/activity-log/api/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

// moving to grpc client, log client should be used (generated
// by protoc for .proto file already)
type Activities struct {
	client api.Activity_LogClient
}

var ErrIDNotFound = fmt.Errorf("ID not found")

// initialize client with action connection
func NewActivities(URL string) Activities {
	conn, err := grpc.Dial(URL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Did not connect: %v", err.Error())
	}
	client := api.NewActivity_LogClient(conn)
	return Activities{client: client}
}

// define Methods to Activities
func (c *Activities) Insert(ctx context.Context, activity *api.Activity) (int, error) {
	res, err := c.client.Insert(ctx, activity)
	if err != nil {
		return 0, fmt.Errorf("Insert error: %v", err.Error())
	}
	// InsertResponse messge in proto will be the response
	// use Get<FieldName> functions to retrieve it
	return int(res.GetId()), nil
}

func (c *Activities) Retrieve(ctx context.Context, id int) (*api.Activity, error) {
	res, err := c.client.Retrieve(ctx, &api.RetrieveRequest{Id: int32(id)})
	if err != nil {
		st, _ := status.FromError(err)
		if st.Code() == codes.NotFound {
			return &api.Activity{}, ErrIDNotFound
		} else {
			return &api.Activity{}, fmt.Errorf("Unexpected retrieve failure %w", err)
		}
	}
	return res, nil
}

func (c *Activities) List(ctx context.Context, offset int) ([]*api.Activity, error) {
	res, err := c.client.List(ctx, &api.ListRequest{Offset: int32(offset)})
	if err != nil {
		// List from server returns only internal error now
		// so, no need to check status code
		acts := []*api.Activity{}
		return acts, err
	}

	return res.Activities, nil
}

/*
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

func (c *Activities) List(offset int) ([]api.Activity, error) {
	qDoc := api.ActivityQueryDocument{Offset: offset}
	jsBytes, err := json.Marshal(qDoc)
	if err != nil {
		return nil, err
	}
	reader := bytes.NewReader(jsBytes)
	req, err := http.NewRequest(http.MethodGet, c.URL+"/list", reader)
	if err != nil {
		return nil, err
	}
	// send the request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode == 404 {
		return nil, errors.New("Not Found")
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var list []api.Activity
	err = json.Unmarshal(body, &list)
	if err != nil {
		return nil, err
	}
	return list, nil
}
*/

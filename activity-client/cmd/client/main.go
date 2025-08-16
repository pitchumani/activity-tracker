package main

// command line client to add and fetch activities

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/pitchumani/activity-tracker/activity-client/internal/client"
	api "github.com/pitchumani/activity-tracker/activity-log/api/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const defaultURL = "localhost:8080"

func main() {
	add := flag.Bool("add", false, "Add activity")
	get := flag.Bool("get", false, "Get activities")
	list := flag.Bool("list", false, "Get activities from offset id")
	flag.Parse()

	activitiesClient := client.NewActivities(defaultURL)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	switch {
	case *get:
		// get option
		if len(os.Args) != 3 {
			log.Fatal("Usage: -get id")
			os.Exit(1)
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatal("Error: Invalid value in the place of a number")
			os.Exit(1)
		}
		a, err := activitiesClient.Retrieve(ctx, id)
		if err != nil {
			log.Fatalf("Error: %s", err.Error())
			os.Exit(1)
		}
		// get activity as string and print it on the console
		// String() interface method is defined for Activity type
		fmt.Printf("Activity: %s", asString(a))
	case *add:
		// add option
		if len(os.Args) != 3 {
			log.Fatal("Usage: -add \"message\"")
			os.Exit(1)
		}
		a := api.Activity{Time: timestamppb.New(time.Now()), Description: os.Args[2]}
		id, err := activitiesClient.Insert(ctx, &a)
		// update the auto assigned id to activity, used for logging
		a.Id = int32(id)
		if err != nil {
			log.Fatalf("Error: %s", err.Error())
			os.Exit(1)
		}
		log.Printf("Added: %s as %d\n", asString(&a), id)
	case *list:
		// list option
		if len(os.Args) != 3 {
			log.Fatal("Usage: -list id_offset");
			os.Exit(1)
		}
		offset, err := strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatal("Error: Invalid offset value")
			os.Exit(1)
		}
		acts, err := activitiesClient.List(ctx, offset)
		if err != nil {
			log.Fatalf("Error: %s", err.Error())
			os.Exit(1)
		}
		for _, act := range acts {
			fmt.Printf("%s\n", asString(act))
		}
	default:
		flag.Usage()
		os.Exit(1)
	}
}

func asString(a *api.Activity) string {
	return fmt.Sprintf("ID: %d %d-%02d-%02d\t\"%s\"", a.Id,
		a.Time.AsTime().Year(), a.Time.AsTime().Month(), a.Time.AsTime().Day(), a.Description)
}

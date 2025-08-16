package server

import (
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	api "github.com/pitchumani/activity-tracker/activity-log/api/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
	  _ "github.com/mattn/go-sqlite3"
)

// define activity info for the server package
// sql.DB is db handler for any relation databases
type Activities struct {
	mtx sync.Mutex
	db  *sql.DB
}

const dbfile string = "activities.db"

const create_table_sql string = `
    CREATE TABLE IF NOT EXISTS activities (
    id INTEGER NOT NULL PRIMARY KEY,
    time DATETIME NOT NULL,
    description TEXT
);`

func NewActivities() (*Activities, error) {
	db, err := sql.Open("sqlite3", dbfile)
	if err != nil {
		return nil, err
	}
	if _, err := db.Exec(create_table_sql); err != nil {
		return nil, err
	}

	return &Activities{
		db: db,
	}, nil
}

// define a insert method for Activities type
func (acts *Activities) Insert(activity *api.Activity) (int, error) {
	acts.mtx.Lock()
	defer acts.mtx.Unlock()
	res, err := acts.db.Exec("INSERT INTO activities VALUES(NULL,?,?);",
		activity.Time.AsTime(), activity.Description)
	if err != nil {
		return 0, err
	}

	var id int64
	if id, err = res.LastInsertId(); err != nil {
		return 0, err
	}
	log.Printf("Added %v\n", activity)
	return int(id), nil
}

var ErrIDNotFound = fmt.Errorf("ID not found")

// define a retrieve method for Activities type
func (acts *Activities) Retrieve(id int) (*api.Activity, error) {
	log.Printf("Getting id %d\n", id)
	acts.mtx.Lock()
	defer acts.mtx.Unlock()
	// using QueryRow that results almost one row
	// db.Query will result into multiple rows (if possible)
	row := acts.db.QueryRow(`
             SELECT id, time, description
             FROM activities
             WHERE id=?`, id)
	activity := api.Activity{}
	var err error
	// need to convert time into protobuf datatype timestamp
	var time time.Time
	if err = row.Scan(&activity.Id, &time, &activity.Description);
	   err == sql.ErrNoRows {
		log.Printf("ID %d is not found\n", id)
		return &activity, ErrIDNotFound
	}
	activity.Time = timestamppb.New(time)
	return &activity, err
}

// retrieve more rows for List Method
func (acts *Activities) List(offset int) ([]*api.Activity, error) {
	log.Printf("Getting list from offset %d\n", offset)

	rows, err := acts.db.Query(
		"SELECT * FROM activities WHERE ID > ? ORDER BY id DESC LIMIT 100",
		offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data := []*api.Activity{}
	for rows.Next() {
		i := api.Activity{}
		var time time.Time
		err = rows.Scan(&i.Id, &time, &i.Description)
		if err != nil {
			return nil, err
		}
		i.Time = timestamppb.New(time)
		data = append(data, &i)
	}
	return data, nil
}

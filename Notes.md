
### sqlite3
install sqlite-utils

use1:
$ sqlite-utils activities.db "select * from activities" --table
  id  time                  description
----  --------------------  -----------------------
   1  2025-02-09T16:34:04Z  eve bike class
   2  2025-02-10T16:34:04Z  Morning walking 5kms
   3  2025-03-10T16:34:04Z  Evening cycle for 12kms

GoLang SQLite setup
go get github.com/mattn/go-sqlite3


### to test all tests in a module
go test ./... -cover -v

### example GET/POST methods using curl

* check get, post methods
curl -iX GET localhost:8080
curl -iX GET localhost:8080
* exclude the verbose (silent mode)
curl -X GET -s localhost:8080


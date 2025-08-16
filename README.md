# Activity Tracker
Activity tracker using Go

It is a JSON HTTP web service.

## APIs
### insert data
```bash
curl -iX POST localhost:8080 -d '{"activity": {"description": "morning walking", "time":"2025-08-09T12:42:31Z"}}'

HTTP/1.1 200 OK
{"id": 1}
```
### retrieve data
```bash
curl -iX GET localhost:8080 -d '{"id": 1}'

{"activity": {"description": "morning walking", "time":"2025-08-09T12:42:31Z", "id": 1}}
```

## gRPC server
JSON based REST service to gRPC service
- JSON based is human readable, so, more data, expensive to Serialize and De-serialize
- gRPC - binary format, quicker to send, faster to serialize and de-serialize
  - service can be specified in .proto files
  - protoc - protobuf compiler can generate code into many languages

- grpc service can't be invoke like REST api service. With reflection service enabled
in grpc server, we can make grpc service request like:
```bash
$ grpcurl -plaintext -d '{ "description": "evening walking 5kms" }' localhost:8080 api.v1.Activity_Log/Insert
{
  "id": 5
}
```
- To enable this grpc server must have registered the reflection (go package google.golang.org/grpc/reflection)
- Without reflection, we can request service like:
```bash
grpcurl -plaintext -d '{ "id": 1 }' -proto ./activity-log/api/v1/activity.proto  localhost:8080 api.v1.Activity_Log/Retrieve
{
  "id": 1,
  "time": "2025-02-09T16:34:04Z",
  "description": " eve bike class"
}
```

## dependencies
* protobuf - `brew install protobuf`
* grpc - `brew install grpc`
* protoc-gen go package - `go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.26`
  - set path to go modules e.g. PATH="$PATH:$(go env GOPATH)/bin"
  - protoc-gen-go binary must have been installed in the GOPATH
* generate code from .proto file
  - `protoc --go_out=. --go_opt=paths=source_relative --proto_path=.  activity-log/api/v1/*.proto`

* grpc generator is from different go package
  - go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

* grpcurl - to make grpc service request in human readable JSON like form

## Credits
Following the [example by Adam Gordon Bell](https://earthly.dev/blog/golang-http/).


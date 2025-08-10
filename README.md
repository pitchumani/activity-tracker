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

## Credits
Following the [example by Adam Gordon Bell](https://earthly.dev/blog/golang-http/).


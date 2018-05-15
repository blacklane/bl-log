Log 
===

Install with `dep ensure -add github.com/blacklane/bl-log`.

Usage
-----

### Request logger

Uses a `log.Record` to messure the duration of the request. Will appear as a `request_finished` 
or `request_error` event and have information about the response code and URL.

```go
http.ListenAndServe(":1234", log.L(myHandler))
```

### Log durations
```go
r := log.NewRecord("event_name")
timeConsumingOperation()
r.Log("some description")
```

### Simple logging
```go
log.Log("event_name", "some description")
log.Error(errors.New("some error"))
```

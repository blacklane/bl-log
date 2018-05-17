Log 
===

Does much less than the average logger, does more than we want from the average logger.  
Aside from the basic `Log` and `Error`, this logging utility includes a struct used to meassure times between calls and provides a usefull `http.Handler` to be used as a middleware to meassure the duration of requests.

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

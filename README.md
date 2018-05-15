Log 
===

Install with `dep ensure -add github.com/blacklane/bl-log`.

Usage
-----

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

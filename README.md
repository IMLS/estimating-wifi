# go-session-counter

## To run

```
go run session-counter.go
```

This assumes `config.yaml` and `auth.yaml` are in `/etc/session-counter`. Also, `manufacturers.sqlite`. These are set up by an ansible playbook.

## To test

```
go test *.go
go test api/*.go
```



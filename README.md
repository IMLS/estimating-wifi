# go-session-counter

## To run

```
go run session-counter.go
```

This assumes `config.yaml` and `auth.yaml` are in `/opt/imls`. Also, `manufacturers.sqlite`. These are set up by an ansible playbook.

## To test

```
go test *.go
go test api/*.go
```



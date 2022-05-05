# System statistic daemon
Graduate project for Golang courses (OTUS)

Main requirements available at https://github.com/OtusGolang/final_project

Technical requirements available at https://github.com/OtusGolang/final_project/blob/master/05-system-stats-daemon.md

### Main features:
* GRPC server
* Simple client
* In-memory storage
* Supported OS: Mac OS, Linux
* Concurrency in getting snapshots

### How to run daemon
``go run cmd/daemon/main.go``

``go run cmd/daemon/main.go -port=50000``

#### Available flags:
* port - GRPC server port (default 50005)

### How to run client
``go run cmd/client/main.go``

``go run cmd/client/main.go -port=50000``

``go run cmd/client/main.go -port=50000 -n=1 -m=5``

#### Available flags:
* port - GRPC server port (default 50005)
* n - getting stats frequency (in sec, default 5)
* m - average stats interval (in sec, default 15)

### How to build daemon

``go build -v -o ./bin/daemon ./cmd/daemon && ./bin/daemon -port=5000``

### How to build client

``go build -v -o ./bin/client ./cmd/client && ./bin/client -port=5000 -n=1 -m=5``

### Configs
Are available in ``configs/config.yaml``

In "stats" group they are means to get (true) or not to get (false) specified kind of statistics

For example
* get only load average stats
```
stats:
  loadavg: true
  cpu: false
  disk: false
  nettop: false
  netstat: false
```

* don't get any kind of stats (by pass mode)
```
stats:
  loadavg: false
  cpu: false
  disk: false
  nettop: false
  netstat: false
```

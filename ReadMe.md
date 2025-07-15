
```
go-sample-rest
├── cfg
│   └── config.yaml
├── cmd
│   └── main.go
├── go.mod
├── go.sum
├── internal
│   ├── config
│   │   └── config.go
│   ├── http-server
│   │   ├── handlers
│   │   │   ├── make-move.go
│   │   │   └── new_board.go
│   │   ├── middleware
│   │   │   ├── cors.go
│   │   │   ├── logging.go
│   │   │   ├── request_id.go
│   │   │   └── validation.go
│   │   ├── models
│   │   │   └── models.go
│   │   └── openapi.yaml
│   ├── logger
│   │   └── logger.go
│   ├── repository
│   │   ├── db
│   │   │   └── repository.db
│   │   └── sqlite_repo.go
│   └── services
│       └── service.go
├── LICENSE.txt
└── ReadMe.md
```
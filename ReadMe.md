
```
go-sample-rest/
├── cmd/
│     └── api/
│         └── main.go               # Entry point, router and server initialization.
├── internal/
│     ├── handlers/                 # HTTP request handlers (controllers)
│     ├── services/                 # Game business logic
│     ├── models/                   # Game data structures
│     ├── middleware/               # 
│     └── repository/               # Not needed yet, but we keep it in mind for storing game states
├── go.mod
└── ...
```
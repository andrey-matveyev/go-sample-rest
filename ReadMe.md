
```
go-sample-rest/
├── cmd/
│     └── api/
│         └── main.go               # Entry point, router and server initialization.
├── internal/
│     ├── handlers/                 # HTTP request handlers (controllers)
│     │     ├── game_handler.go     # Logic for new-board and make-move
│     ├── services/                 # Game business logic
│     │     └── game_service.go     # Game rules, checking wins, generating moves
│     ├── models/                   # Game data structures
│     │     └── game.go             # Board, GameState, Player, etc.
│     └── repository/               # Not needed yet, but we keep it in mind for storing game states
│           └── game_sqlite_repo.go # If we store games in memory
├── go.mod
└── ...
```
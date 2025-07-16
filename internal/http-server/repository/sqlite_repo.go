package repository

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type Storage interface {
	SaveNewGame(player int) (int64, error)
	UpdateGame() error
	SaveNewMove(idGame int64, board string) error
	Shutdown() error
}

type sqliteStorage struct {
	db *sql.DB
}

func NewStorage(path string) (Storage, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("open sql.DB error: %w", err)
	}

	// Create table 'games'
	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS 
	games(
		id         INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		player     INTEGER NOT NULL,
		is_over    INTEGER NOT NULL,
		win_player INTEGER NOT NULL,
        start_time TEXT NOT NULL,
        stop_time  TEXT
    );
	`)
	if err != nil {
		return nil, fmt.Errorf("error prepare create table 'games': %w", err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("error create table 'games': %w", err)
	}

	// Create table 'moves'
	stmt, err = db.Prepare(`
    CREATE TABLE IF NOT EXISTS 
	moves(
		id        INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		id_game   INTEGER NOT NULL,
        board     TEXT NOT NULL,
		move_time TEXT NOT NULL,
        FOREIGN KEY (id_game) REFERENCES games(id)
    );
	`)
	if err != nil {
		return nil, fmt.Errorf("error prepare create table 'moves': %w", err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("error create table 'moves': %w", err)
	}

	return &sqliteStorage{db: db}, nil
}

var (
	ErrURLNotFound = errors.New("url not found")
	ErrURLExists   = errors.New("url exists")
)

func (item *sqliteStorage) SaveNewGame(player int) (int64, error) {
	stmt, err := item.db.Prepare(`
	INSERT INTO games
        (player, is_over, win_player, start_time) 
    VALUES
        (?, 0, 0, datetime('now'))
	`)
	if err != nil {
		return 0, fmt.Errorf("prepare 'INSERT' sql.DB error: %w", err)
	}

	res, err := stmt.Exec(player)
	if err != nil {
		return 0, fmt.Errorf("execute 'INSERT' sql.Stmt error: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error getting last insert id: %w", err)
	}

	return id, nil
}

func (item *sqliteStorage) UpdateGame() error {
	return nil
}

func (item *sqliteStorage) SaveNewMove(idGame int64, board string) error {
	stmt, err := item.db.Prepare(`
	INSERT INTO moves
        (id_game, board, move_time) 
    VALUES
        (?, ?, datetime('now'))
	`)
	if err != nil {
		return fmt.Errorf("prepare 'INSERT' sql.DB error: %w", err)
	}

	_, err = stmt.Exec(idGame, board)
	if err != nil {
		return fmt.Errorf("execute 'INSERT' sql.Stmt error: %w", err)
	}

	return nil
}

func (item *sqliteStorage) Shutdown() error {
	if item.db == nil {
		return nil
	}
	err := item.db.Close()
	if err != nil {
		return fmt.Errorf("close sql.DB  error: %w", err)
	}
	return nil
}

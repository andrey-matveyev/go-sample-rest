package repository

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type Storage interface {
	SaveNewGame() (int64, error)
	UpdateGame() (int64, error)
	SaveNewMove() (int64, error)
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

	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS 
	games(
		id         INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		is_over    INTEGER NOT NULL,
		win_player INTEGER NOT NULL,
        start_time TEXT NOT NULL,
        stop_time  TEXT
    );
    CREATE TABLE IF NOT EXISTS 
	moves(
		id        INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		game_id   INTEGER NOT NULL,
		player    INTEGER NOT NULL,
        board     TEXT NOT NULL,
		move_time TEXT NOT NULL,
        FOREIGN KEY (game_id) REFERENCES games(id)
    );
	`)
	if err != nil {
		return nil, fmt.Errorf("prepare 'CREATE TABLE' sql.DB error: %w", err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("execute 'CREATE TABLE' sql.Stmt error: %w", err)
	}

	return &sqliteStorage{db: db}, nil
}

var (
	ErrURLNotFound = errors.New("url not found")
	ErrURLExists   = errors.New("url exists")
)

func (item *sqliteStorage) SaveNewGame() (int64, error) {
	stmt, err := item.db.Prepare(`
	INSERT INTO games
        (is_over, win_player, game_start) 
    VALUES
        (0, 0, datetime('now'))
	`)
	if err != nil {
		return 0, fmt.Errorf("prepare 'INSERT' sql.DB error: %w", err)
	}

	res, err := stmt.Exec()
	if err != nil {
		return 0, fmt.Errorf("execute 'INSERT' sql.Stmt error: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error getting last insert id: %w", err)
	}

	return id, nil
}

func (item *sqliteStorage) UpdateGame() (int64, error) {
	return 0, nil
}

func (item *sqliteStorage) SaveNewMove() (int64, error) {
	return 0, nil
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

/*

"github.com/mattn/go-sqlite3"

func (s *Storage) SaveURL(urlToSave string, alias string) (int64, error) {
	const op = "storage.sqlite.SaveURL"

	stmt, err := s.db.Prepare("INSERT INTO url(url, alias) VALUES(?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.Exec(urlToSave, alias)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrURLExists)
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: failed to get last insert id: %w", op, err)
	}

	return id, nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	const op = "storage.sqlite.GetURL"

	stmt, err := s.db.Prepare("SELECT url FROM url WHERE alias = ?")
	if err != nil {
		return "", fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	var resURL string

	err = stmt.QueryRow(alias).Scan(&resURL)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", storage.ErrURLNotFound
		}

		return "", fmt.Errorf("%s: execute statement: %w", op, err)
	}

	return resURL, nil
}
*/

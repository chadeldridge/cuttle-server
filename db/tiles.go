package db

/*
import "fmt"

type Tiles struct {
	DB
}

func NewTiles(db DB) (Tiles, error) {
	r := Tiles{db}
	if db == nil {
		return r, fmt.Errorf("db.NewTiles: db is nil")
	}

	r.migrate()
	return r, nil
}

func (r Tiles) migrate() error {
	query := `
	CREATE TABLE IF NOT EXISTS tiles (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		hideCmd BOOLEAN NOT NULL DEFAULT FALSE,
		hideExp BOOLEAN NOT NULL DEFAULT FALSE,
		name VARCHAR(255) NOT NULL,
		cmd text NOT NULL,
		exp text,
		display_size INTEGER NOT NULL DEFAULT 40,
	)`

	if err := r.Exec(query); err != nil {
		return fmt.Errorf("db.Tiles.Exec: %w", err)
	}

	return nil
}
*/

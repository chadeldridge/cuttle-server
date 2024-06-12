package db

type DB interface {
	Open() error
	Migrate(query string) error
	Close() error
}

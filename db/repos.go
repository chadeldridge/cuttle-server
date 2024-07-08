package db

import "fmt"

type Repos struct {
	DB
	Profiles    Profiles
	Tiles       Tiles
	Groups      Groups
	Servers     Servers
	Connectors  Connectors
	AuthMethods AuthMethods
}

func NewRepos(db DB) (Repos, error) {
	r := Repos{DB: db}
	if db == nil {
		return r, fmt.Errorf("db.NewRepos: db is nil")
	}

	var err error
	// Attach auth_methods.
	r.AuthMethods, err = NewAuthMethods(db)
	if err != nil {
		return r, fmt.Errorf("db.NewRepos: failed to attach auth_methods: %w", err)
	}

	/*
		r.Profiles = NewProfiles(db)
		r.Tiles = NewTiles(db)
		r.Groups = NewGroups(db)
		r.Servers = NewServers(db)
		r.Connectors = NewConnectors(db)
	*/

	return r, nil
}

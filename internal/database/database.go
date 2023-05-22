package database

import (
	surreal "github.com/surrealdb/surrealdb.go"
	"tidal.lol/internal/logging"
)

var (
	DB *surreal.DB
)

func NewDatabase(url, username, password string) *surreal.DB {
	db, err := surreal.New(url)
	if err != nil {
		logging.Logger.Panic().
			Err(err).
			Msg("Failed to connect to database")
		return nil
	}

	_, err = db.Signin(map[string]string{
		"user": username,
		"pass": password,
	})
	if err != nil {
		logging.Logger.Panic().
			Err(err).
			Msg("Failed to sign in to database")
		return nil
	}

	_, err = db.Use("tidal", "lol")
	if err != nil {
		logging.Logger.Panic().
			Err(err).
			Msg("Failed to use database")
		return nil
	}

	return db
}

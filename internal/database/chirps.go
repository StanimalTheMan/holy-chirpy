package database

import (
	"errors"
	"os"
)

type Chirp struct {
	ID       int    `json:"id"`
	Body     string `json:"body"`
	AuthorID int    `json:"author_id"`
}

func (db *DB) CreateChirp(body string, userIDInt int) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	id := len(dbStructure.Chirps) + 1
	chirp := Chirp{
		ID:       id,
		Body:     body,
		AuthorID: userIDInt,
	}
	dbStructure.Chirps[id] = chirp

	err = db.writeDB(dbStructure)
	if err != nil {
		return Chirp{}, err
	}

	return chirp, nil
}

func (db *DB) GetChirps() ([]Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	chirps := make([]Chirp, 0, len(dbStructure.Chirps))
	for _, chirp := range dbStructure.Chirps {
		chirps = append(chirps, chirp)
	}

	return chirps, nil
}

func (db *DB) GetChirp(id int) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	chirp, ok := dbStructure.Chirps[id]
	if !ok {
		return Chirp{}, os.ErrNotExist
	}

	return chirp, nil
}

func (db *DB) DeleteChirp(chirpID, userID int) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}

	chirp, ok := dbStructure.Chirps[chirpID]
	if !ok {
		return os.ErrNotExist
	}

	if chirp.AuthorID != userID {
		return errors.New("id mismatch")
	}

	delete(dbStructure.Chirps, chirpID)
	err = db.writeDB(dbStructure)
	if err != nil {
		return err
	}

	return nil
}

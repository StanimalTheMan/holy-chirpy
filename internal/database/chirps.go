package database

import (
	"os"
	"strconv"
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

func (db *DB) GetChirps(authorID string) ([]Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	chirps := make([]Chirp, 0, len(dbStructure.Chirps))
	for _, chirp := range dbStructure.Chirps {
		if authorID != "" {
			// check if authorID matches chirp's authorID before adding
			authorIDInt, _ := strconv.Atoi(authorID)
			if authorIDInt == chirp.AuthorID {
				chirps = append(chirps, chirp)
			}
		} else {
			chirps = append(chirps, chirp)
		}
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

func (db *DB) DeleteChirp(id int) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}

	delete(dbStructure.Chirps, id)
	err = db.writeDB(dbStructure)
	if err != nil {
		return err
	}

	return nil
}

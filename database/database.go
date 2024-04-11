package database

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}

type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
}

// NewDB creates a new database connection
// and creates the database file if it does not exist
func NewDB(path string) (*DB, error) {
	db := &DB{
		path: path,
		mux:  &sync.RWMutex{},
	}
	if err := db.ensureDB(); err != nil {
		return nil, err
	}
	return db, nil
}

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(body string) (Chirp, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	// Load existing db
	dbStruct, err := db.loadDB()

	if err != nil {
		return Chirp{}, err
	}

	nextID := len(dbStruct.Chirps) + 1
	fmt.Println("in CreateChirp", body)
	chirp := Chirp{
		ID:   nextID,
		Body: body,
	}

	// Add the chirp to db
	dbStruct.Chirps[nextID] = chirp

	// Write updated db back to disk
	if err := db.writeDB(dbStruct); err != nil {
		return Chirp{}, err
	}

	return chirp, nil
}

// GetChirps returns all chirps in the database
func (db *DB) GetChirps() ([]Chirp, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()

	// Load the existing db
	dbStruct, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	// Convert the map of chirps to a slice
	chirps := make([]Chirp, 0, len(dbStruct.Chirps))
	for _, chirp := range dbStruct.Chirps {
		fmt.Println("chirp", chirp)
		chirps = append(chirps, chirp)
	}

	return chirps, nil
}

// ensureDB creates a new database file if it doesn't exist
func (db *DB) ensureDB() error {
	if _, err := os.Stat(db.path); os.IsNotExist(err) {
		// Create a new empty database structure
		dbStruct := DBStructure{
			Chirps: make(map[int]Chirp),
		}

		// Write empty structure to disk
		if err := db.writeDB(dbStruct); err != nil {
			return err
		}
	}
	return nil
}

// loadDB reads the database file into memory
func (db *DB) loadDB() (DBStructure, error) {
	data, err := os.ReadFile(db.path)

	if err != nil {
		return DBStructure{}, err
	}

	var dbStruct DBStructure
	if err := json.Unmarshal(data, &dbStruct); err != nil {
		return DBStructure{}, err
	}
	fmt.Println("loaded data", dbStruct)
	return dbStruct, nil
}

// writeDB writes the database file to disk
func (db *DB) writeDB(dbStructure DBStructure) error {
	data, err := json.MarshalIndent(dbStructure, "", " ")
	if err != nil {
		return err
	}

	return os.WriteFile(db.path, data, 0644)
}

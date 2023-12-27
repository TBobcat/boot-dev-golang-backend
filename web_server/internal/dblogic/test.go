// local server, saves data into a json file on this machine
package dblogic

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
)

// type DBStructure struct {
// 	Chirps map[int]Chirp `json:"chirps"`
// }

type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

type DB struct {
	path string
	mu   *sync.RWMutex
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
}

// Foo prints something, testing importing internal packagae
func Foo() {
	fmt.Println("Using an internal package !")
}

// NewDB creates a new database connection
// and creates the database file if it doesn't exist
// func NewDB(path string) (*DB, error)
func NewDB(path string) (*DB, error) {
	db := &DB{
		path: path,
		mu:   &sync.RWMutex{},
	}
	err := db.ensureDB()
	return db, err
}

// CreateChirp creates a new chirp and saves it to disk
// func (db *DB) CreateChirp(body string) (Chirp, error)

// GetChirps returns all chirps in the database
// receiver is a pointer of a DB object
func (db *DB) GetChirp(id int) (Chirp, error) {
	chirp := Chirp{
		ID:   id,
		Body: "dummy hardcoded",
	}
	return chirp, nil
}

// func (db *DB) GetChirps() ([]Chirp, error) {
// 	return nil, nil
// }

// ensureDB creates a new database file if it doesn't exist
// func (db *DB) ensureDB() error
func (db *DB) ensureDB() error {
	_, err := os.ReadFile(db.path)
	if errors.Is(err, os.ErrNotExist) {
		return db.createDB()
	}
	return err
}

func (db *DB) createDB() error {
	dbStructure := DBStructure{
		Chirps: map[int]Chirp{},
	}
	return db.writeDB(dbStructure)
}

// loadDB reads the database file into memory
//func (db *DB) loadDB() (DBStructure, error)

// writeDB writes the database file to disk
// func (db *DB) writeDB(dbStructure DBStructure) error
func (db *DB) writeDB(dbStructure DBStructure) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	dat, err := json.Marshal(dbStructure)
	if err != nil {
		return err
	}

	err = os.WriteFile(db.path, dat, 0600)
	if err != nil {
		return err
	}
	return nil
}

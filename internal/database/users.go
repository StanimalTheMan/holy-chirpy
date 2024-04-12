package database

import (
	"fmt"
	"os"
)

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (db *DB) CreateUser(email, password string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	id := len(dbStructure.Users) + 1

	// check that no 2 users can be created with same email address
	isEmailUnique := checkValidEmail(email, dbStructure)
	if !isEmailUnique {
		return User{}, err
	}

	user := User{
		ID:       id,
		Email:    email,
		Password: password,
	}
	dbStructure.Users[email] = user

	err = db.writeDB(dbStructure)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func checkValidEmail(email string, dbStructure DBStructure) bool {
	isEmailUnique := true
	for _, user := range dbStructure.Users {
		if user.Email == email {
			isEmailUnique = false
		}
	}
	return isEmailUnique
}

func (db *DB) GetUser(email string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	user, ok := dbStructure.Users[email]
	fmt.Println(email)
	if !ok {
		return User{}, os.ErrNotExist
	}

	return user, nil
}

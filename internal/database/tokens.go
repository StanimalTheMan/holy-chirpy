package database

import (
	"fmt"
	"time"
)

func (db *DB) CreateRevokedToken(token string) error {
	fmt.Println("bruhhh yyyyy")
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}
	dbStructure.RevokedTokens[token] = time.Now()
	err = db.writeDB(dbStructure)
	if err != nil {
		return err
	}
	fmt.Printf("REVOKED MAP %v", dbStructure.RevokedTokens)
	return nil
}

func (db *DB) CheckIfRevoked(token string) bool {
	fmt.Println("yooo")
	fmt.Println("token")
	dbStructure, _ := db.loadDB()
	fmt.Println(dbStructure.RevokedTokens)
	if _, ok := dbStructure.RevokedTokens[token]; ok {
		fmt.Println("bruhhh")
		return true
	}
	return false
}

package db

import (
	"genMaterials/common"
	"log"

	"github.com/jmoiron/sqlx"
)

// var db *sqlx.DB

type DB struct {
	*sqlx.DB
}

func DBConnect() (*DB, error) {
	config := common.GetYamlConfig().Database
	// db, err := sqlx.Open("postgres", "host=localhost user=postgres password=password dbname=imaterials sslmode=disable")
	conString := "host=" + config.Host + " user=" + config.User + " password=" + config.Password + " dbname=" + config.Name + " sslmode=" + config.SslMode
	db, err := sqlx.Open(config.Type, conString)

	if err != nil {
		log.Fatalln(err)
	}
	return &DB{db}, err
}

// func SqlxClose() {
// 	if db != nil {
// 		db.Close()
// 	}
// }

func SqlxConnect() (*DB, error) {
	// return db
	db, err := DBConnect()

	return db, err
}

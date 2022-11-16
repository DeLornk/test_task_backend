package src

import (
	"database/sql"
	"log"
)

var DB *sql.DB

func InitDB() error {
	var err error

	//DB, err = sql.Open("postgres", "host=localhost port=5432 user=avito password=avito dbname=avito sslmode=disable")
	DB, err = sql.Open("postgres", "host=postgres port=5432 user=avito password=avito dbname=avito sslmode=disable")
	if err != nil {
		log.Println(err)
		return err
	}
	// try to ping our DB
	err = DB.Ping()
	if err != nil {
		log.Println(err)
		return err
	}

	log.Println("Подключились к базе данных")
	return nil
}

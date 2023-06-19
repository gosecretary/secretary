package storage

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)


func DatabaseQuery(query string) (*sql.Rows, error) {
	db := OpenDatabase()
	
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err, ": ", query)
		return nil, err
	}
	defer db.Close()

	log.Println("DatabaseQuery:", query)
	return rows, nil
}

func DatabaseExec(query string) (*sql.Result, error) {
	db := OpenDatabase()
	
	result, err := db.Exec(query)
	if err != nil {
		log.Fatal(err, ": ", query)
		return nil, err
	}
	rowCount, err := result.RowsAffected()
	if err != nil {
		log.Fatal(err, ": ", query)
		return nil, err
	}
	defer db.Close()
	
	log.Println("DatabaseExec, Rows Affected:", rowCount, " - Query:", query)
	return &result, nil
}


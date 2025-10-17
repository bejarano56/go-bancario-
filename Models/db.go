package Models

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() error {
	var err error
	// Cambia los datos de conexión según tu entorno
	dsn := "root:@tcp(127.0.0.1:3306)/banco?parseTime=true"
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	return DB.Ping()
}

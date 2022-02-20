package utility

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func NewDbClient() *sql.DB {
	godotenv.Load(".env")
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=require",
		os.Getenv("HOST"), os.Getenv("DB_PORT"), os.Getenv("dbUser"), os.Getenv("dbPassword"), os.Getenv("dbName"))
	conn, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(fmt.Sprintf("Verbindung mit Datenbank kann nicht hergestellt werden: %v", err))
	}
	conn.SetMaxIdleConns(2)
	conn.SetMaxOpenConns(20)
	return conn
}

func GetStringValue(s sql.NullString) string {
	if s.Valid {
		return s.String
	}
	return ""
}

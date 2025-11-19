// database/postgres.go
package database

import (
	"database/sql"
	"fmt"
	"log"
	"uas-pelaporan-prestasi-backend/config"

	_ "github.com/lib/pq"
)

var PostgresDB *sql.DB

func ConnectPostgres() {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Env("DB_HOST", "localhost"),
		config.Env("DB_PORT", "5432"),
		config.Env("DB_USER", "postgres"),
		config.Env("DB_PASSWORD", "postgres"),
		config.Env("DB_NAME", "prestasi_unair"),
		config.Env("DB_SSLMODE", "disable"))

	var err error
	PostgresDB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Gagal koneksi Postgres:", err)
	}

	if err = PostgresDB.Ping(); err != nil {
		log.Fatal("Ping Postgres gagal:", err)
	}

	log.Println("✅ Postgres terkoneksi!")
}
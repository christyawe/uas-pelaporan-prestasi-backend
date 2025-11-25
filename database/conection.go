package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
	"uas-pelaporan-prestasi-backend/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	_ "github.com/lib/pq"
)

var PostgresDB *sql.DB
var MongoClient *mongo.Client
var MongoDB *mongo.Database

func Connect() {
	// Koneksi Postgres
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

	// Koneksi Mongo
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	MongoClient, err = mongo.Connect(ctx, options.Client().ApplyURI(config.Env("MONGO_URI", "mongodb://localhost:27017")))
	if err != nil {
		log.Fatal("Gagal koneksi MongoDB:", err)
	}

	err = MongoClient.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Ping MongoDB gagal:", err)
	}

	MongoDB = MongoClient.Database(config.Env("MONGO_DATABASE", "prestasi_unair"))
	log.Println("✅ MongoDB terkoneksi!")
}
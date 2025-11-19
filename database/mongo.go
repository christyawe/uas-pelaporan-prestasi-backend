// database/mongo.go
package database

import (
	"context"
	"log"
	"time"
	"uas-pelaporan-prestasi-backend/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client
var MongoDB *mongo.Database

func ConnectMongo() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.Env("MONGO_URI", "mongodb://localhost:27017")))
	if err != nil {
		log.Fatal("Gagal koneksi MongoDB:", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Ping MongoDB gagal:", err)
	}

	MongoClient = client
	MongoDB = client.Database(config.Env("MONGO_DATABASE", "prestasi_unair"))
	log.Println("✅ MongoDB terkoneksi!")
}
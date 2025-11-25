package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"
	"uas-pelaporan-prestasi-backend/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func SoftDeleteAchievementMongo(achievementID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(achievementID)
	if err != nil {
		return err
	}

	collection := database.MongoDB.Collection("achievements")
	filter := bson.M{"_id": objectID, "status": "draft"}
	update := bson.M{"$set": bson.M{"status": "deleted", "updated_at": time.Now()}}

	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("prestasi bukan draft atau tidak ditemukan")
	}
	return nil
}

func UpdateAchievementReferencePostgres(achievementID string) error {
	query := `
		UPDATE achievement_references 
		SET status = 'deleted', updated_at = NOW() 
		WHERE mongo_achievement_id = $1 AND status = 'draft'
		RETURNING id
	`
	var id string
	err := database.PostgresDB.QueryRow(query, achievementID).Scan(&id)
	if err == sql.ErrNoRows {
		return fmt.Errorf("prestasi bukan draft atau tidak ditemukan")
	}
	return err
}
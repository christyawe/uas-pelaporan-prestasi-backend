package model

import (
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Achievement struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	StudentID     uuid.UUID          `bson:"student_id"`
	Title         string             `bson:"title"`
	Description   string             `bson:"description"`
	Category      string             `bson:"category"`
	DynamicFields map[string]any     `bson:"dynamic_fields"`
	Status        string             `bson:"status"` // include 'deleted'
	CreatedAt     time.Time          `bson:"created_at"`
	UpdatedAt     time.Time          `bson:"updated_at"`
}

type AchievementReference struct {
	ID                 uuid.UUID  `db:"id"`
	StudentID          uuid.UUID  `db:"student_id"`
	MongoAchievementID string     `db:"mongo_achievement_id"`
	Status             string     `db:"status"` // include 'deleted'
	SubmittedAt        *time.Time `db:"submitted_at"`
	VerifiedAt         *time.Time `db:"verified_at"`
	VerifiedBy         *uuid.UUID `db:"verified_by"`
	RejectionNote      *string    `db:"rejection_note"`
	CreatedAt          time.Time  `db:"created_at"`
	UpdatedAt          time.Time  `db:"updated_at"`
}
// database/migration.go
package database

import (
	"context"
	"database/sql"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// SeedDummyData mengisi data awal untuk testing
func SeedDummyData() {
	log.Println("🌱 Memulai Seeding Data Dummy...")
	ctx := context.Background()

	// 1. SEED ROLES
	roles := map[string]string{
		"Admin":      "Administrator Sistem",
		"Mahasiswa":  "Pengguna Mahasiswa",
		"Dosen Wali": "Dosen Verifikator",
	}
	roleIDs := make(map[string]string)
	for name, desc := range roles {
		var id string
		err := PostgresDB.QueryRow("SELECT id FROM roles WHERE name = $1", name).Scan(&id)
		if err == sql.ErrNoRows {
			err = PostgresDB.QueryRow(
				"INSERT INTO roles (name, description) VALUES ($1, $2) RETURNING id",
				name, desc,
			).Scan(&id)
			if err != nil {
				log.Fatalf("❌ Gagal seed role %s: %v", name, err)
			}
			log.Printf("✅ Role Created: %s", name)
		}
		roleIDs[name] = id
	}

	// 2. SEED USERS
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost) // Password default hashed
	// Admin (ga dipake ID-nya, jadi langsung seed tanpa var)
	seedUser(PostgresDB, "admin", "admin@univ.ac.id", string(passwordHash), "Super Admin", roleIDs["Admin"])

	// Dosen
	dosenUserID := seedUser(PostgresDB, "dosen1", "dosen1@univ.ac.id", string(passwordHash), "Dr. Budi Santoso", roleIDs["Dosen Wali"])

	// Mahasiswa
	mhsUserID := seedUser(PostgresDB, "mhs1", "mhs1@univ.ac.id", string(passwordHash), "Andi Pratama", roleIDs["Mahasiswa"])

	// 3. SEED PROFILES
	// Lecturer (Dosen)
	var lecturerID string
	err := PostgresDB.QueryRow("SELECT id FROM lecturers WHERE user_id = $1", dosenUserID).Scan(&lecturerID)
	if err == sql.ErrNoRows {
		err = PostgresDB.QueryRow(`
			INSERT INTO lecturers (user_id, lecturer_id, department)
			VALUES ($1, $2, $3) RETURNING id`,
			dosenUserID, "19850101201001", "Teknik Informatika",
		).Scan(&lecturerID)
		if err != nil {
			log.Fatalf("❌ Gagal seed lecturer: %v", err)
		}
		log.Println("✅ Profile Lecturer Created")
	}

	// Student (Mahasiswa) - Link advisor ke lecturerID
	var studentID string
	err = PostgresDB.QueryRow("SELECT id FROM students WHERE user_id = $1", mhsUserID).Scan(&studentID)
	if err == sql.ErrNoRows {
		err = PostgresDB.QueryRow(`
			INSERT INTO students (user_id, student_id, program_study, academic_year, advisor_id)
			VALUES ($1, $2, $3, $4, $5) RETURNING id`,
			mhsUserID, "20210001", "D4 Teknik Informatika", "2021", lecturerID,
		).Scan(&studentID)
		if err != nil {
			log.Fatalf("❌ Gagal seed student: %v", err)
		}
		log.Println("✅ Profile Student Created")
	}

	// 4. SEED ACHIEVEMENTS (Mongo + Postgres)
	var count int
	PostgresDB.QueryRow("SELECT COUNT(*) FROM achievement_references WHERE student_id = $1", studentID).Scan(&count)
	if count == 0 {
		log.Println("Creating Dummy Achievement...")
		// Insert Mongo
		achCollection := MongoDB.Collection("achievements")
		mongoDoc := bson.M{
			"student_id": studentID, // UUID as string
			"type":       "competition",
			"title":      "Juara 1 Hackathon Nasional",
			"description": "Menang lomba coding di Jakarta",
			"created_at": time.Now(),
			"updated_at": time.Now(),
			"points":     100,
			"tags":       []string{"coding", "java", "winner"},
			"details": bson.M{
				"competition_name":  "Gemastik 2025",
				"competition_level": "national",
				"rank":              1,
				"location":          "Jakarta",
				"event_date":        time.Now(),
			},
			"attachments": []bson.M{
				{
					"file_name":   "sertifikat.pdf",
					"file_url":    "http://localhost:3000/uploads/dummy.pdf",
					"file_type":   "application/pdf",
					"uploaded_at": time.Now(),
				},
			},
		}
		res, err := achCollection.InsertOne(ctx, mongoDoc)
		if err != nil {
			log.Fatalf("❌ Gagal seed Mongo achievement: %v", err)
		}
		mongoID := res.InsertedID.(primitive.ObjectID).Hex()
		log.Printf("✅ Mongo Achievement Inserted ID: %s", mongoID)

		// Insert Postgres Reference
		_, err = PostgresDB.Exec(`
			INSERT INTO achievement_references (student_id, mongo_achievement_id, status, submitted_at, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6)`,
			studentID, mongoID, "submitted", time.Now(), time.Now(), time.Now(),
		)
		if err != nil {
			log.Fatalf("❌ Gagal seed Postgres reference: %v", err)
		}
		log.Println("✅ Postgres Reference Created")
	}
	log.Println("🎉 Seeding Selesai!")
}

// seedUser helper
func seedUser(db *sql.DB, username, email, passwordHash, fullName, roleID string) string {
	var id string
	err := db.QueryRow("SELECT id FROM users WHERE email = $1", email).Scan(&id)
	if err == sql.ErrNoRows {
		err = db.QueryRow(`
			INSERT INTO users (username, email, password_hash, full_name, role_id)
			VALUES ($1, $2, $3, $4, $5) RETURNING id`,
			username, email, passwordHash, fullName, roleID,
		).Scan(&id)
		if err != nil {
			log.Fatalf("❌ Gagal seed user %s: %v", username, err)
		}
		log.Printf("✅ User Created: %s (%s)", username, email)
	}
	return id
}
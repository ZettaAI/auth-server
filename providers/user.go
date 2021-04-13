package providers

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// User user info
type User struct {
	ID        uint   `gorm:"default:uuid_generate_v3()"`
	Email     string `gorm:"primaryKey"`
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func addUser(email string, name string) int {
	db, err := gorm.Open(mysql.Open(os.Getenv("USER_DB_STRING")), &gorm.Config{})
	if err != nil {
		log.Fatalf("Could not connect to DB server: %v", err.Error())
	}

	user := User{Email: email, Name: name}
	result := db.Create(&user)
	if result.Error != nil {
		log.Fatalf("Could not insert user: %v", result.Error.Error())
	}
	return int(user.ID)
}

// GetUserID retrieve user ID from DB
func GetUserID(email string) int {
	db, err := gorm.Open(mysql.Open(os.Getenv("USER_DB_STRING")), &gorm.Config{})
	if err != nil {
		log.Fatalf("Could not connect to DB server: %v", err.Error())
	}

	user := User{Email: email}
	db.First(&user)
	return int(user.ID)
}

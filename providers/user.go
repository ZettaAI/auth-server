package providers

import (
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// User user info
type User struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	Email     string `gorm:"unique"`
	Name      string
	CreatedAt []uint8
	UpdatedAt []uint8
}

// AddUser ass new user to DB
func AddUser(email string, name string) {
	db, err := gorm.Open(mysql.Open(os.Getenv("USER_DB_STRING")), &gorm.Config{})
	if err != nil {
		log.Fatalf("Could not connect to DB server: %v", err.Error())
	}

	// if !db.Migrator().HasTable(&User{}) {
	// 	db.Migrator().CreateTable(&User{})
	// }

	user := User{}
	result := db.Where(&User{Email: email}).First(&user)
	if result.Error != nil {
		log.Printf("Could not find user %s: %v\n", email, result.Error.Error())
		user = User{Email: email, Name: name}
		result = db.Create(&user)
		if result.Error != nil {
			log.Printf("Could not add new user: %v\n", result.Error.Error())
		}
	}
}

// GetUserID retrieve user ID from DB
func GetUserID(email string) int {
	db, err := gorm.Open(mysql.Open(os.Getenv("USER_DB_STRING")), &gorm.Config{})
	if err != nil {
		log.Fatalf("Could not connect to DB server: %v", err.Error())
	}

	user := User{}
	result := db.Where(&User{Email: email}).First(&user)
	if result.Error != nil {
		log.Printf("Could not find user %s: %v\n", email, result.Error.Error())
	}
	return int(user.ID)
}

package persist

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Keywords is a custom type that represents a slice of strings
type Keywords []string

// Scan implements the sql.Scanner interface
func (k *Keywords) Scan(value interface{}) error {
	// Convert the value to a []byte
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to convert value to []byte: %v", value)
	}

	// Unmarshal the []byte into a []string
	if err := json.Unmarshal(bytes, k); err != nil {
		return fmt.Errorf("failed to unmarshal []byte into []string: %v", err)
	}
	return nil
}

// Value implements the driver.Valuer interface
func (k Keywords) Value() (driver.Value, error) {
	// Marshal the []string into a []byte
	bytes, err := json.Marshal(k)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal []string into []byte: %v", err)
	}
	return bytes, nil
}

// UserData represents the structure of the data in the SQLite database
type UserData struct {
	ChatID    int64    `gorm:"column:chat_id"`    // Telegram chat ID
	Username  string   `gorm:"column:username"`   // Telegram username
	OzbGood   bool     `gorm:"column:ozb_good"`   // watch deals with 25+ upvotes in the last 24 hours
	OzbSuper  bool     `gorm:"column:ozb_super"`  // watch deals with 50+ upvotes in the last 24 hours
	Keywords  Keywords `gorm:"column:keywords"`   // list of keywords / deals to watch for
	OzbSent   Keywords `gorm:"column:ozb_sent"`   // comma separated list of ozb deals sent to user
	AmzDaily  bool     `gorm:"column:amz_daily"`  // watch top daily deals on amazon
	AmzWeekly bool     `gorm:"column:amz_weekly"` // watch top weekly deals on amazon
	AmzSent   Keywords `gorm:"column:amz_sent"`   // comma separated list of amz deals sent to user
}

// TableName sets the table name for the UserData struct
func (UserData) TableName() string {
	return "users"
}

func SqliteToMongoDB(sqliteDBFile, mongoURI, mongoDBName, mongoCollectionName string) error {
	// Connect to the SQLite database
	db, err := gorm.Open("sqlite3", sqliteDBFile)
	if err != nil {
		return err
	}
	defer db.Close()

	db.SingularTable(true)

	// Connect to the MongoDB database
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		return err
	}
	defer client.Disconnect(context.TODO())

	// Get the users collection
	usersCollection := client.Database(mongoDBName).Collection(mongoCollectionName)

	// Get the user data from the SQLite database
	var users []UserData
	if err := db.Find(&users).Error; err != nil {
		return err
	}

	// Insert the user data into the MongoDB collection
	for _, user := range users {
		userDoc := bson.M{
			"chat_id":    user.ChatID,
			"username":   user.Username,
			"ozb_good":   user.OzbGood,
			"ozb_super":  user.OzbSuper,
			"keywords":   user.Keywords,
			"ozb_sent":   user.OzbSent,
			"amz_daily":  user.AmzDaily,
			"amz_weekly": user.AmzWeekly,
			"amz_sent":   user.AmzSent,
		}
		if _, err := usersCollection.InsertOne(context.TODO(), userDoc); err != nil {
			return err
		}
	}
	return nil
}

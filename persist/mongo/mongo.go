package persist

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/intothevoid/kramerbot/models"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

// Use sync mutex to prevent data race
var storeMutex *sync.Mutex = &sync.Mutex{}

type MongoStoreDB struct {
	Coll     *mongo.Collection
	Name     string
	CollName string
	Logger   *zap.Logger
}

// Connect to the mongo database
func New(mongoUri string, dbName string, collName string, logger *zap.Logger) (*MongoStoreDB, error) {
	// Establish a connection to the MongoDB database.
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoUri))
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	// Create a collection to store data.
	usersColl := client.Database(dbName).Collection(collName)

	// Finally, return the UserDB value.
	return &MongoStoreDB{
		Coll:     usersColl,
		Name:     dbName,
		CollName: collName,
		Logger:   logger,
	}, nil
}

// Close the database
func (mdb *MongoStoreDB) Close() error {
	// To close the MongoDB connection, you can use the Disconnect function.
	err := mdb.Coll.Database().Client().Disconnect(context.TODO())
	if err != nil {
		return err
	}
	return nil
}

// Ping checks the connection to the database
func (mdb *MongoStoreDB) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second) // Increased timeout to 15 seconds
	defer cancel()

	if mdb.Coll == nil || mdb.Coll.Database() == nil || mdb.Coll.Database().Client() == nil {
		return fmt.Errorf("mongo client or collection not initialized")
	}

	err := mdb.Coll.Database().Client().Ping(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to ping mongo database: %w", err)
	}
	return nil
}

// Add user to the database
func (mdb *MongoStoreDB) AddUser(user *models.UserData) error {
	// Add sync mutex to the database
	storeMutex.Lock()
	defer storeMutex.Unlock()

	usersColl := mdb.Coll

	// First, you can use the FindOne function to check if the user already exists.
	result := usersColl.FindOne(context.TODO(), bson.M{"chat_id": user.ChatID})
	if result.Err() == nil {
		// create new error
		return fmt.Errorf("user already exists: %v", user.ChatID)
	}

	// Then, you can use the InsertOne function to insert a single document into the collection.
	_, err := usersColl.InsertOne(context.TODO(), user)
	if err != nil {
		return err
	}

	// The user has been successfully inserted into the collection.
	return nil
}

// Update user in the database
func (mdb *MongoStoreDB) UpdateUser(user *models.UserData) error {
	// Add sync mutex to the database
	storeMutex.Lock()
	defer storeMutex.Unlock()

	usersColl := mdb.Coll

	filter := bson.M{"chat_id": user.ChatID}
	update := bson.M{"$set": user}
	result, err := usersColl.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		// create new error
		return fmt.Errorf("user not found: %v", user.ChatID)
	}

	if result.ModifiedCount == 0 {
		// create new error
		return fmt.Errorf("user not modified: %v", user.ChatID)
	}

	// The user has been successfully updated in the collection.
	return nil
}

// Delete user from the database
func (mdb *MongoStoreDB) DeleteUser(user *models.UserData) error {
	// Add sync mutex to the database
	storeMutex.Lock()
	defer storeMutex.Unlock()

	usersColl := mdb.Coll

	filter := bson.M{"chat_id": user.ChatID}
	result, err := usersColl.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		// create new error
		return fmt.Errorf("user not found: %v", user.ChatID)
	}

	// The user has been successfully deleted from the collection.
	return nil
}

// Get user from the database by chat_id
func (mdb *MongoStoreDB) GetUser(chatID int64) (*models.UserData, error) {
	// Add sync mutex to the database
	storeMutex.Lock()
	defer storeMutex.Unlock()

	usersColl := mdb.Coll

	var user models.UserData
	user.ChatID = chatID
	err := usersColl.FindOne(context.TODO(), bson.M{"chat_id": user.ChatID}).Decode(&user)
	if err != nil {
		return nil, err
	}

	// The user has been successfully retrieved from the collection.
	return &user, nil
}

// Read all users from the database
func (mdb *MongoStoreDB) ReadUserStore() (*models.UserStore, error) {
	// Add sync mutex to the database
	storeMutex.Lock()
	defer storeMutex.Unlock()

	userStore := &models.UserStore{
		Users: make(map[int64]*models.UserData),
	}

	cursor, err := mdb.Coll.Find(context.TODO(), bson.M{})
	if err != nil {
		mdb.Logger.Error("Error getting users", zap.Error(err))
		return nil, err
	}

	for cursor.Next(context.TODO()) {
		user := &models.UserData{}
		err := cursor.Decode(user)
		if err != nil {
			mdb.Logger.Error("Error decoding user", zap.Error(err))
			return nil, err
		}

		userStore.Users[user.ChatID] = user
	}

	return userStore, nil
}

// Write *models.UserStore to the database
func (mdb *MongoStoreDB) WriteUserStore(userStore *models.UserStore) error {
	// Add sync mutex to the database
	storeMutex.Lock()
	defer storeMutex.Unlock()

	// Syncrhoniza
	// Convert the userStore map to a slice of users.
	var users []interface{}
	for _, user := range userStore.Users {
		users = append(users, user)
	}

	// Instead of dropping the collection, iterate over the users and update each one.
	for _, user := range users {
		// Cast the user back to *models.UserData
		userData := user.(*models.UserData)

		// Create the filter for the user
		filter := bson.M{"chat_id": userData.ChatID}

		// Create the update document
		update := bson.M{"$set": userData}

		// Update the user
		_, err := mdb.Coll.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			// If the user does not exist, insert it
			if err == mongo.ErrNoDocuments {
				_, err := mdb.Coll.InsertOne(context.TODO(), userData)
				if err != nil {
					return err
				}
			} else {
				return err
			}
		}
	}

	// The users have been successfully updated in the collection.
	return nil
}

package persist

import (
	"context"
	"time"

	"github.com/intothevoid/kramerbot/models"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type MongoStoreDB struct {
	Coll     *mongo.Collection
	Name     string
	CollName string
	Logger   *zap.Logger
}

// Connect to the mongo database
func New(dbName string, collName string, logger *zap.Logger) (*MongoStoreDB, error) {
	// Establish a connection to the MongoDB database.
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
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

// Add user to the database
func (mdb *MongoStoreDB) AddUser(user *models.UserData) error {
	usersColl := mdb.Coll

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
	usersColl := mdb.Coll

	filter := bson.M{"id": user.ChatID}
	update := bson.M{"$set": user}
	_, err := usersColl.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}

	// The user has been successfully updated in the collection.
	return nil
}

// Delete user from the database
func (mdb *MongoStoreDB) DeleteUser(user *models.UserData) error {
	usersColl := mdb.Coll

	filter := bson.M{"id": user.ChatID}
	_, err := usersColl.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}

	// The user has been successfully deleted from the collection.
	return nil
}

// Get user from the database by chat_id
func (mdb *MongoStoreDB) GetUser(chatID int64) (*models.UserData, error) {
	usersColl := mdb.Coll

	var user models.UserData
	err := usersColl.FindOne(context.TODO(), bson.M{"id": user.ChatID}).Decode(&user)
	if err != nil {
		return nil, err
	}

	// The user has been successfully retrieved from the collection.
	return &user, nil
}

// Read all users from the database
func (mdb *MongoStoreDB) ReadUserStore() (*models.UserStore, error) {

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
	// Convert the userStore map to a slice of users.
	var users []interface{}
	for _, user := range userStore.Users {
		users = append(users, user)
	}

	// First, you can use the Drop function to delete all documents in the collection.
	_, err := mdb.Coll.DeleteMany(context.TODO(), bson.M{})
	if err != nil {
		return err
	}

	// Then, you can use the InsertMany function to insert multiple documents into the collection.
	_, err = mdb.Coll.InsertMany(context.TODO(), users)
	if err != nil {
		return err
	}

	// The users have been successfully inserted into the collection.
	return nil
}
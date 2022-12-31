package persist_test

import (
	"testing"

	persist "github.com/intothevoid/kramerbot/persist/mongo"
)

func TestConvertSqltoMongo(t *testing.T) {
	sqliteDBFile := "users.db"
	mongoURI := "mongodb://localhost:27017"
	mongoDBName := "usersdb"
	mongoCollectionName := "users"

	if err := persist.SqliteToMongoDB(sqliteDBFile, mongoURI, mongoDBName, mongoCollectionName); err != nil {
		t.Fatal(err)
	}

	t.Log("Successfully converted SQLite database to MongoDB collection!")
}

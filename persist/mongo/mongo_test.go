package persist

// import (
// 	"context"
// 	"fmt"
// 	"reflect"
// 	"testing"

// 	"github.com/intothevoid/kramerbot/models"
// 	_ "github.com/lib/pq"
// 	_ "github.com/mattn/go-sqlite3"
// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/mongo"
// 	"go.uber.org/zap"
// )

// var NewFunc = New

// const MONGOURI = "mongodb://localhost:27017"

// func TestNew(t *testing.T) {
// 	// Create a mock logger
// 	logger := zap.NewNop()

// 	// Test the case where a connection to the MongoDB can be established
// 	dbName := "test_db"
// 	collName := "test_coll"
// 	expectedDB := &MongoStoreDB{
// 		Coll:     &mongo.Collection{},
// 		Name:     dbName,
// 		CollName: collName,
// 		Logger:   logger,
// 	}
// 	// Set up a mock MongoDB client that returns the expectedDB when Connect is called
// 	client := &mockMongoClient{
// 		expectedDB: expectedDB,
// 	}
// 	// Inject the mock client into the New function
// 	NewFunc = func(mongoUri string, dbName string, collName string, logger *zap.Logger) (*MongoStoreDB, error) {
// 		return client.New(dbName, collName, logger)
// 	}
// 	gotDB, err := NewFunc(MONGOURI, dbName, collName, logger)
// 	if err != nil {
// 		t.Errorf("Unexpected error: %v", err)
// 	}
// 	if !reflect.DeepEqual(gotDB, expectedDB) {
// 		t.Errorf("Unexpected value for gotDB. Expected: %v, got: %v", expectedDB, gotDB)
// 	}

// 	// Test the case where a connection to the MongoDB cannot be established
// 	client = &mockMongoClient{
// 		connectErr: fmt.Errorf("unable to connect to MongoDB"),
// 	}
// 	NewFunc = func(mongoUri string, dbName string, collName string, logger *zap.Logger) (*MongoStoreDB, error) {
// 		return client.New(dbName, collName, logger)
// 	}
// 	gotDB, err = NewFunc(MONGOURI, dbName, collName, logger)
// 	if err == nil {
// 		t.Errorf("Expected an error, got nil")
// 	}
// 	if gotDB != nil {
// 		t.Errorf("Expected nil value for gotDB, got: %v", gotDB)
// 	}
// }

// type mockMongoClient struct {
// 	expectedDB *MongoStoreDB
// 	connectErr error
// }

// func (c *mockMongoClient) New(dbName string, collName string, logger *zap.Logger) (*MongoStoreDB, error) {
// 	if c.connectErr != nil {
// 		return nil, c.connectErr
// 	}
// 	return c.expectedDB, nil
// }

// // The below test should only be run when a MongoDB instance is running on localhost:27017
// // Otherwise, the test will fail / may hang indefinitely
// func TestMongoStoreDB_AddUser(t *testing.T) {
// 	mongoStoreDB, err := New(MONGOURI, "test_db", "test_coll", zap.NewNop())
// 	if err != nil {
// 		t.Fatal("Unable to create MongoStoreDB")
// 	}

// 	type fields struct {
// 		Coll     *mongo.Collection
// 		Name     string
// 		CollName string
// 		Logger   *zap.Logger
// 	}
// 	type args struct {
// 		user *models.UserData
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		args    args
// 		wantErr bool
// 	}{
// 		{
// 			name: "Add user",
// 			fields: fields{
// 				Coll:     mongoStoreDB.Coll,
// 				Name:     mongoStoreDB.Name,
// 				CollName: mongoStoreDB.CollName,
// 				Logger:   mongoStoreDB.Logger,
// 			},
// 			args: args{
// 				user: &models.UserData{
// 					ChatID:         123,
// 					Username:       "JaneDoe",
// 					OzbGood:        true,
// 					OzbSuper:       true,
// 					Keywords:       []string{"test3", "test4"},
// 					OzbSent:        []string{"test3", "test4"},
// 					AmzDaily:       true,
// 					AmzWeekly:      true,
// 					AmzSent:        []string{"test3", "test4"},
// 					UsernameChosen: "JaneDoe",
// 					Password:       "password",
// 				},
// 			},
// 			wantErr: false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			// Clear the database before each test run
// 			_, err := mongoStoreDB.Coll.DeleteMany(context.Background(), bson.M{})
// 			if err != nil {
// 				t.Fatal("Unable to clear database: ", err)
// 			}

// 			mdb := &MongoStoreDB{
// 				Coll:     tt.fields.Coll,
// 				Name:     tt.fields.Name,
// 				CollName: tt.fields.CollName,
// 				Logger:   tt.fields.Logger,
// 			}
// 			if err := mdb.AddUser(tt.args.user); (err != nil) != tt.wantErr {
// 				t.Errorf("MongoStoreDB.AddUser() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 		})
// 	}
// }

// func TestMongoStoreDB_UpdateUser(t *testing.T) {
// 	mongoStoreDB, err := New(MONGOURI, "test_db", "test_coll", zap.NewNop())
// 	if err != nil {
// 		t.Fatal("Unable to create MongoStoreDB")
// 	}

// 	// Insert test user
// 	user := &models.UserData{
// 		ChatID:         123,
// 		Username:       "JohnDoe",
// 		OzbGood:        true,
// 		OzbSuper:       true,
// 		Keywords:       []string{"test1", "test2"},
// 		OzbSent:        []string{"test1", "test2"},
// 		AmzDaily:       true,
// 		AmzWeekly:      true,
// 		AmzSent:        []string{"test1", "test2"},
// 		UsernameChosen: "JohnDoe",
// 		Password:       "password",
// 	}

// 	type fields struct {
// 		Coll     *mongo.Collection
// 		Name     string
// 		CollName string
// 		Logger   *zap.Logger
// 	}
// 	type args struct {
// 		user *models.UserData
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		args    args
// 		wantErr bool
// 	}{
// 		// TODO: Add test cases.
// 		{
// 			name: "Update user",
// 			fields: fields{
// 				Coll:     mongoStoreDB.Coll,
// 				Name:     mongoStoreDB.Name,
// 				CollName: mongoStoreDB.CollName,
// 				Logger:   mongoStoreDB.Logger,
// 			},
// 			args: args{
// 				user: &models.UserData{
// 					ChatID:         123,
// 					Username:       "JaneDoe",
// 					OzbGood:        true,
// 					OzbSuper:       true,
// 					Keywords:       []string{"test3", "test4"},
// 					OzbSent:        []string{"test3", "test4"},
// 					AmzDaily:       true,
// 					AmzWeekly:      true,
// 					AmzSent:        []string{"test3", "test4"},
// 					UsernameChosen: "JaneDoe",
// 					Password:       "password",
// 				},
// 			},
// 			wantErr: false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			// Clear the database before each test run
// 			_, err := mongoStoreDB.Coll.DeleteMany(context.Background(), bson.M{})
// 			if err != nil {
// 				t.Fatal("Unable to clear database: ", err)
// 			}
// 			err = mongoStoreDB.AddUser(user)
// 			if err != nil {
// 				t.Logf("Unable to insert test user %d", user.ChatID)
// 			}

// 			mdb := &MongoStoreDB{
// 				Coll:     tt.fields.Coll,
// 				Name:     tt.fields.Name,
// 				CollName: tt.fields.CollName,
// 				Logger:   tt.fields.Logger,
// 			}
// 			if err := mdb.UpdateUser(tt.args.user); (err != nil) != tt.wantErr {
// 				t.Errorf("MongoStoreDB.UpdateUser() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 		})
// 	}
// }

// func TestMongoStoreDB_DeleteUser(t *testing.T) {

// 	mongoStoreDB, err := New(MONGOURI, "test_db", "test_coll", zap.NewNop())
// 	if err != nil {
// 		t.Fatal("Unable to create MongoStoreDB")
// 	}

// 	// Insert test user
// 	user := &models.UserData{
// 		ChatID:         123,
// 		Username:       "JaneDoe",
// 		OzbGood:        true,
// 		OzbSuper:       true,
// 		Keywords:       []string{"test3", "test4"},
// 		OzbSent:        []string{"test3", "test4"},
// 		AmzDaily:       true,
// 		AmzWeekly:      true,
// 		AmzSent:        []string{"test3", "test4"},
// 		UsernameChosen: "JaneDoe",
// 		Password:       "password",
// 	}

// 	err = mongoStoreDB.AddUser(user)
// 	if err != nil {
// 		t.Logf("Unable to insert test user %d", user.ChatID)
// 	}

// 	type fields struct {
// 		Coll     *mongo.Collection
// 		Name     string
// 		CollName string
// 		Logger   *zap.Logger
// 	}
// 	type args struct {
// 		user *models.UserData
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		args    args
// 		wantErr bool
// 	}{
// 		{
// 			name: "Delete user",
// 			fields: fields{
// 				Coll:     mongoStoreDB.Coll,
// 				Name:     mongoStoreDB.Name,
// 				CollName: mongoStoreDB.CollName,
// 				Logger:   mongoStoreDB.Logger,
// 			},
// 			args: args{
// 				user: &models.UserData{
// 					ChatID:         123,
// 					Username:       "JaneDoe",
// 					OzbGood:        true,
// 					OzbSuper:       true,
// 					Keywords:       []string{"test3", "test4"},
// 					OzbSent:        []string{"test3", "test4"},
// 					AmzDaily:       true,
// 					AmzWeekly:      true,
// 					AmzSent:        []string{"test3", "test4"},
// 					UsernameChosen: "JaneDoe",
// 					Password:       "password",
// 				},
// 			},
// 			wantErr: false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			mdb := &MongoStoreDB{
// 				Coll:     tt.fields.Coll,
// 				Name:     tt.fields.Name,
// 				CollName: tt.fields.CollName,
// 				Logger:   tt.fields.Logger,
// 			}
// 			if err := mdb.DeleteUser(tt.args.user); (err != nil) != tt.wantErr {
// 				t.Errorf("MongoStoreDB.DeleteUser() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 		})
// 	}
// }

// func TestMongoStoreDB_GetUser(t *testing.T) {

// 	mongoStoreDB, err := New(MONGOURI, "test_db", "test_coll", zap.NewNop())
// 	if err != nil {
// 		t.Fatal("Unable to create MongoStoreDB")
// 	}

// 	// Insert test user
// 	user := &models.UserData{
// 		ChatID:         123,
// 		Username:       "JaneDoe",
// 		OzbGood:        true,
// 		OzbSuper:       true,
// 		Keywords:       []string{"test3", "test4"},
// 		OzbSent:        []string{"test3", "test4"},
// 		AmzDaily:       true,
// 		AmzWeekly:      true,
// 		AmzSent:        []string{"test3", "test4"},
// 		UsernameChosen: "JaneDoe",
// 		Password:       "password",
// 	}

// 	err = mongoStoreDB.AddUser(user)
// 	if err != nil {
// 		t.Logf("Unable to insert test user %d", user.ChatID)
// 	}

// 	type fields struct {
// 		Coll     *mongo.Collection
// 		Name     string
// 		CollName string
// 		Logger   *zap.Logger
// 	}
// 	type args struct {
// 		chatID int64
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		args    args
// 		want    *models.UserData
// 		wantErr bool
// 	}{
// 		{
// 			name: "Get user",
// 			fields: fields{
// 				Coll:     mongoStoreDB.Coll,
// 				Name:     mongoStoreDB.Name,
// 				CollName: mongoStoreDB.CollName,
// 				Logger:   mongoStoreDB.Logger,
// 			},
// 			args: args{
// 				chatID: 123,
// 			},
// 			want:    user,
// 			wantErr: false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			mdb := &MongoStoreDB{
// 				Coll:     tt.fields.Coll,
// 				Name:     tt.fields.Name,
// 				CollName: tt.fields.CollName,
// 				Logger:   tt.fields.Logger,
// 			}
// 			got, err := mdb.GetUser(tt.args.chatID)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("MongoStoreDB.GetUser() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("MongoStoreDB.GetUser() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestMongoStoreDB_ReadUserStore(t *testing.T) {

// 	mongoStoreDB, err := New(MONGOURI, "test_db", "test_coll", zap.NewNop())
// 	if err != nil {
// 		t.Fatal("Unable to create MongoStoreDB")
// 	}

// 	// Insert test user
// 	user1 := &models.UserData{
// 		ChatID:         123,
// 		Username:       "JaneDoe",
// 		OzbGood:        true,
// 		OzbSuper:       true,
// 		Keywords:       []string{"test3", "test4"},
// 		OzbSent:        []string{"test3", "test4"},
// 		AmzDaily:       true,
// 		AmzWeekly:      true,
// 		AmzSent:        []string{"test3", "test4"},
// 		UsernameChosen: "JaneDoe",
// 		Password:       "password",
// 	}

// 	user2 := &models.UserData{
// 		ChatID:         456,
// 		Username:       "JohnDoe",
// 		OzbGood:        true,
// 		OzbSuper:       true,
// 		Keywords:       []string{"test3", "test4"},
// 		OzbSent:        []string{"test3", "test4"},
// 		AmzDaily:       true,
// 		AmzWeekly:      true,
// 		AmzSent:        []string{"test3", "test4"},
// 		UsernameChosen: "JohnDoe",
// 		Password:       "password",
// 	}

// 	err = mongoStoreDB.AddUser(user1)
// 	if err != nil {
// 		t.Logf("Unable to insert test user %d", user1.ChatID)
// 	}

// 	err = mongoStoreDB.AddUser(user2)
// 	if err != nil {
// 		t.Logf("Unable to insert test user %d", user2.ChatID)
// 	}

// 	// Create test user store
// 	testStore := &models.UserStore{
// 		Users: map[int64]*models.UserData{
// 			user1.ChatID: user1,
// 			user2.ChatID: user2,
// 		},
// 	}

// 	type fields struct {
// 		Coll     *mongo.Collection
// 		Name     string
// 		CollName string
// 		Logger   *zap.Logger
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		want    *models.UserStore
// 		wantErr bool
// 	}{
// 		{
// 			name: "Read user store",
// 			fields: fields{
// 				Coll:     mongoStoreDB.Coll,
// 				Name:     mongoStoreDB.Name,
// 				CollName: mongoStoreDB.CollName,
// 				Logger:   mongoStoreDB.Logger,
// 			},
// 			want:    testStore,
// 			wantErr: false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			mdb := &MongoStoreDB{
// 				Coll:     tt.fields.Coll,
// 				Name:     tt.fields.Name,
// 				CollName: tt.fields.CollName,
// 				Logger:   tt.fields.Logger,
// 			}
// 			got, err := mdb.ReadUserStore()
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("MongoStoreDB.ReadUserStore() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("MongoStoreDB.ReadUserStore() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestMongoStoreDB_WriteUserStore(t *testing.T) {
// 	mongoStoreDB, err := New(MONGOURI, "test_db", "test_coll", zap.NewNop())
// 	if err != nil {
// 		t.Fatal("Unable to create MongoStoreDB")
// 	}

// 	// Insert test user
// 	user1 := &models.UserData{
// 		ChatID:         123,
// 		Username:       "JaneDoe",
// 		OzbGood:        true,
// 		OzbSuper:       true,
// 		Keywords:       []string{"test3", "test4"},
// 		OzbSent:        []string{"test3", "test4"},
// 		AmzDaily:       true,
// 		AmzWeekly:      true,
// 		AmzSent:        []string{"test3", "test4"},
// 		UsernameChosen: "JaneDoe",
// 		Password:       "password",
// 	}

// 	user2 := &models.UserData{
// 		ChatID:         456,
// 		Username:       "JohnDoe",
// 		OzbGood:        true,
// 		OzbSuper:       true,
// 		Keywords:       []string{"test3", "test4"},
// 		OzbSent:        []string{"test3", "test4"},
// 		AmzDaily:       true,
// 		AmzWeekly:      true,
// 		AmzSent:        []string{"test3", "test4"},
// 		UsernameChosen: "JohnDoe",
// 		Password:       "password",
// 	}

// 	err = mongoStoreDB.AddUser(user1)
// 	if err != nil {
// 		t.Logf("Unable to insert test user %d", user1.ChatID)
// 	}

// 	err = mongoStoreDB.AddUser(user2)
// 	if err != nil {
// 		t.Logf("Unable to insert test user %d", user2.ChatID)
// 	}

// 	// Create test user store
// 	testStore := &models.UserStore{
// 		Users: map[int64]*models.UserData{
// 			user1.ChatID: user1,
// 			user2.ChatID: user2,
// 		},
// 	}
// 	type fields struct {
// 		Coll     *mongo.Collection
// 		Name     string
// 		CollName string
// 		Logger   *zap.Logger
// 	}
// 	type args struct {
// 		userStore *models.UserStore
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		args    args
// 		wantErr bool
// 	}{
// 		{
// 			name: "Write user store",
// 			fields: fields{
// 				Coll:     mongoStoreDB.Coll,
// 				Name:     mongoStoreDB.Name,
// 				CollName: mongoStoreDB.CollName,
// 				Logger:   mongoStoreDB.Logger,
// 			},
// 			args: args{
// 				userStore: testStore,
// 			},
// 			wantErr: false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			mdb := &MongoStoreDB{
// 				Coll:     tt.fields.Coll,
// 				Name:     tt.fields.Name,
// 				CollName: tt.fields.CollName,
// 				Logger:   tt.fields.Logger,
// 			}
// 			if err := mdb.WriteUserStore(tt.args.userStore); (err != nil) != tt.wantErr {
// 				t.Errorf("MongoStoreDB.WriteUserStore() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 		})
// 	}
// }

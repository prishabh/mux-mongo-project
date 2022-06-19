package mongodb

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Config mongodb connection options
type Config struct {
	// TCP host:port or Unix socket depending on Network
	Address string
	// User/Role name which have access to the DB and schema
	User string
	// Password for the user
	Password string
	// Database name
	Database string
	// Collection name
	Collection string
}

// Client One extra client layer
type Client struct {
	mongodbClient     *mongo.Client
	mongodbCollection *mongo.Collection
}

// NewClient Create mongodb client
func NewClient(mongodbConfig *Config) (*Client, error) {
	uri := fmt.Sprintf("mongodb://%s:%s@%s", mongodbConfig.User,
		mongodbConfig.Password, mongodbConfig.Address)
	// mongo.Connect return mongo.Client method
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		return nil, err
	}
	c := new(Client)
	c.mongodbClient = client
	c.mongodbCollection = client.Database(mongodbConfig.Database).Collection(mongodbConfig.Collection)
	return c, nil
}

// Ping is a user defined method that accepts context.Context
// used to ping the mongoDB, return error if any.
func (c *Client) Ping(ctx context.Context) error {
	// mongo.Client has Ping to ping mongoDB, deadline of
	// the Ping method will be determined by cxt
	// Ping method return error if any occurred, then
	// the error can be handled.
	if err := c.mongodbClient.Ping(ctx, readpref.Primary()); err != nil {
		return err
	}
	fmt.Println("mongodb connection successful")
	return nil
}

// Insert is a user defined method, used to insert
// documents into collection returns result of
// InsertMany and error if any.
func (c *Client) Insert(data interface{}) (interface{}, error) {
	result, err := c.mongodbCollection.InsertOne(context.TODO(), data)
	return result, err
}

// Query is user defined method used to query MongoDB
func (c *Client) Query(filter, results interface{}) error {
	cursor, err := c.mongodbCollection.Find(context.TODO(), filter)
	if err != nil {
		return err
	}
	return cursor.All(context.TODO(), results)
}

// CreateIndex is a user define method, used to
// create table index
func (c *Client) CreateIndex(col, index string) (string, error) {
	model := mongo.IndexModel{Keys: bson.D{{col, index}}}
	name, err := c.mongodbCollection.Indexes().CreateOne(context.TODO(), model)
	return name, err
}

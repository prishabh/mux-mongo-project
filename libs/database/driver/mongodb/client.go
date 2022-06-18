package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Config mongodb connection options.
type Config struct {
	// TCP host:port or Unix socket depending on Network.
	Address string
	// User/Role name which have access to the DB and schema.
	User string
	// Password for the user.
	Password string
	// Database name.
	Database string
}

// Client One extra client layer.
type Client struct {
	mongodbClient *mongo.Client
}

// NewClient Create mongodb client
func NewClient(mongodbConfig *Config) (*Client, error) {
	uri := fmt.Sprintf("mongodb://%s:%s@%s", mongodbConfig.User,
		mongodbConfig.Password, mongodbConfig.Address)
	// mongo.Connect return mongo.Client method
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}
	c := new(Client)
	c.mongodbClient = client
	return c, nil
}

// This is a user defined method that accepts context.Context
// This method used to ping the mongoDB, return error if any.
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

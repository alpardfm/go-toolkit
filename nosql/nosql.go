package nosql

import (
	"context"
	"fmt"

	"github.com/alpardfm/go-toolkit/codes"
	"github.com/alpardfm/go-toolkit/errors"
	"github.com/alpardfm/go-toolkit/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Interface interface {
	Close(ctx context.Context) error
	Find(ctx context.Context, collection string, entity interface{}, filter interface{}, opts ...*options.FindOptions) error
	FindOne(ctx context.Context, collection string, dest interface{}, filter interface{}, opts ...*options.FindOneOptions) error
	InsertOne(ctx context.Context, collection string, data interface{}) (*mongo.InsertOneResult, error)
	UpdateOne(ctx context.Context, collection string, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	UpdateMany(ctx context.Context, collection string, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
}

type Config struct {
	WaitingTime int
	DBUrl       string
	DB          string
}

type mongoDB struct {
	client *mongo.Client
	cfg    Config
	log    log.Interface
}

func Init(cfg Config, log log.Interface) Interface {
	ctx := context.Background()
	client := options.Client().ApplyURI(cfg.DBUrl)
	dbClient, err := mongo.Connect(ctx, client)
	if err != nil {
		log.Fatal(ctx, fmt.Sprintf("[FATAL] cannot connect to dbURL %s, err : %v", cfg.DBUrl, err))
	}
	log.Info(ctx, fmt.Sprintf("NoSQL: dbURL=%s db=%s", cfg.DBUrl, cfg.DB))

	nosql := &mongoDB{
		client: dbClient,
		log:    log,
		cfg:    cfg,
	}

	return nosql
}

func (m *mongoDB) Close(ctx context.Context) error {
	err := m.client.Disconnect(ctx)
	if err != nil {
		return errors.NewWithCode(codes.CodeNoSQLClose, err.Error())
	}
	m.log.Info(ctx, "Connection to MongoDB closed...")

	return nil
}

func (m *mongoDB) Find(ctx context.Context, collection string, dest interface{}, filter interface{}, opts ...*options.FindOptions) error {
	cursor, err := m.client.Database(m.cfg.DB).Collection(collection).Find(ctx, filter, opts...)
	if err != nil {
		return errors.NewWithCode(codes.CodeNoSQLRead, err.Error())
	}

	if err := cursor.All(ctx, dest); err != nil {
		return errors.NewWithCode(codes.CodeNoSQLDecode, err.Error())
	}

	return nil
}

func (m *mongoDB) FindOne(ctx context.Context, collection string, dest interface{}, filter interface{}, opts ...*options.FindOneOptions) error {
	err := m.client.Database(m.cfg.DB).Collection(collection).FindOne(ctx, filter, opts...).Decode(dest)
	if err != nil {
		return errors.NewWithCode(codes.CodeNoSQLDecode, err.Error())
	}

	return nil
}

func (m *mongoDB) InsertOne(ctx context.Context, collection string, data interface{}) (*mongo.InsertOneResult, error) {
	insertResult, err := m.client.Database(m.cfg.DB).Collection(collection).InsertOne(ctx, data)
	if err != nil {
		return nil, errors.NewWithCode(codes.CodeNoSQLInsert, err.Error())
	}

	return insertResult, nil
}

func (m *mongoDB) UpdateOne(ctx context.Context, collection string, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	updateResult, err := m.client.Database(m.cfg.DB).Collection(collection).UpdateOne(ctx, filter, update, opts...)
	if err != nil {
		return nil, errors.NewWithCode(codes.CodeNoSQLUpdate, err.Error())
	}

	return updateResult, nil
}

func (m *mongoDB) UpdateMany(ctx context.Context, collection string, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	updateResult, err := m.client.Database(m.cfg.DB).Collection(collection).UpdateMany(ctx, filter, update, opts...)
	if err != nil {
		return nil, errors.NewWithCode(codes.CodeNoSQLUpdate, err.Error())
	}

	return updateResult, nil
}

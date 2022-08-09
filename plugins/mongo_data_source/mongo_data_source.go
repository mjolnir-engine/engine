package mongo_data_source

import (
	"context"
	"fmt"

	"github.com/mjolnir-mud/engine/plugins/data_sources/pkg/constants"
	"github.com/mjolnir-mud/engine/plugins/data_sources/pkg/errors"
	constants2 "github.com/mjolnir-mud/engine/plugins/mongo_data_source/pkg/constants"
	errors2 "github.com/mjolnir-mud/engine/plugins/mongo_data_source/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type plugin struct {
	database *mongo.Database
}

func (p *plugin) Name() string {
	return "mongo_data_source"
}

func (p *plugin) Start() error {
	logger = log.With().Str("plugin", "mongo_data_source").Logger()
	env := viper.GetString("env")
	viper.SetDefault("mongo_url", "mongodb://localhost:27017")
	viper.SetDefault("database", fmt.Sprintf("mjolnir_%s", env))
	err := viper.BindEnv("mongo_url")

	if err != nil {
		logger.Fatal().Err(err).Msg("error binding env")
		panic(err)
	}

	err = viper.BindEnv("database")

	if err != nil {
		logger.Fatal().Err(err).Msg("error binding env")
		panic(err)
	}

	logger = logger.
		With().
		Str("mongo_url", viper.GetString("mongo_url")).
		Str("database", viper.GetString("database")).
		Logger()

	logger.Info().Msg("starting mongo connection")
	c, err := mongo.NewClient(options.Client().ApplyURI(viper.GetString("mongo_url")))

	if err != nil {
		logger.Fatal().Err(err).Msg("error connecting to mongo")
		panic(err)
	}

	err = c.Connect(context.Background())

	if err != nil {
		logger.Fatal().Err(err).Msg("error connecting to mongo")
		panic(err)
	}

	p.database = c.Database(viper.GetString("database"))

	logger.Info().Msg("mongo connection started")

	return nil
}

func (p *plugin) Stop() error {
	logger.Info().Msg("stopping mongo connection")
	_ = p.database.Client().Disconnect(context.Background())
	return nil
}

func (p *plugin) collection(name string) *mongo.Collection {
	return p.database.Collection(name)
}

var logger zerolog.Logger

var Plugin = &plugin{}

type MongoDataSource struct {
	collectionName string
	collection     *mongo.Collection
	logger         zerolog.Logger
}

func New(collection string) MongoDataSource {
	return MongoDataSource{
		collectionName: collection,
		collection:     Plugin.collection(collection),
		logger:         logger.With().Str("collection", collection).Logger(),
	}
}

func (m MongoDataSource) Name() string {
	return m.collectionName
}

func (m MongoDataSource) Start() error {
	return nil
}

func (m MongoDataSource) Stop() error {
	return nil
}

func (m MongoDataSource) Load(entityId string) (map[string]interface{}, error) {
	m.logger.Debug().Str("entityId", entityId).Msg("loading entity")

	filter := map[string]interface{}{
		"_id": entityId,
	}

	result := map[string]interface{}{}

	err := m.collection.FindOne(context.Background(), filter).Decode(&result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.EntityNotFoundError{ID: entityId}
		}
	}

	cleanId(result)

	return result, nil
}

func (m MongoDataSource) LoadAll() (map[string]map[string]interface{}, error) {
	m.logger.Debug().Msg("loading all entities")

	cursor, err := m.collection.Find(context.Background(), map[string]interface{}{})

	if err != nil {
		return nil, err
	}

	var results map[string]map[string]interface{}
	err = cursor.All(context.Background(), &results)

	if err != nil {
		return nil, err
	}

	return results, nil
}

func (m MongoDataSource) Find(search map[string]interface{}) (map[string]map[string]interface{}, error) {
	m.logger.Debug().Interface("search", search).Msg("searching entities")

	cursor, err := m.collection.Find(context.Background(), search)

	if err != nil {
		return nil, err
	}

	var results map[string]map[string]interface{}
	err = cursor.All(context.Background(), &results)

	if err != nil {
		return nil, err
	}

	cleanIds(results)

	return results, nil
}

func (m MongoDataSource) Save(entityId string, entity map[string]interface{}) error {
	metadata, ok := entity[constants.MetadataKey]

	if !ok {
		return errors.MetadataRequiredError{ID: entityId}
	}

	collection, ok := metadata.(map[string]interface{})[constants2.MetadataCollectionKey]

	if !ok {
		return errors2.CollectionMetadataRequiredError{ID: entityId}
	}

	m.logger.Debug().Interface("entity", entity).Msg("saving entity")

	if collection != m.collectionName {
		return errors2.CollectionMismatchError{
			SourceCollection: m.collectionName,
			TargetCollection: collection.(string),
		}
	}

	entity["_id"] = entityId

	_, err := m.collection.InsertOne(context.Background(), entity)
	return err
}

func cleanIds(entities map[string]map[string]interface{}) {
	for _, entity := range entities {
		cleanId(entity)
	}
}

func cleanId(entity map[string]interface{}) {
	delete(entity, "_id")
}
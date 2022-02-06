package apiserver

import (
	"context"
	"net/http"

	"github.com/go-pkgz/lgr"
	"github.com/parMaster/ethtrx/internal/app/store/mongostore"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func Start(config *Config) error {
	client, err := newDB(config.MongoURI)
	if err != nil {
		lgr.Fatalf("FATAL Failed to connect to Mongo:\n%s", err.Error())
	}

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			lgr.Fatalf("ERROR db.Disconnect: %s", err.Error())
		}
	}()

	db := client.Database("eth")

	store := mongostore.NewStore(db)

	srv := newServer(store, *config)
	srv.logger.Logf("Listening to %s", config.BindAddr)
	return http.ListenAndServe(config.BindAddr, srv)
}

func newDB(connectionURI string) (*mongo.Client, error) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(connectionURI))
	if err != nil {
		return nil, err
	}

	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		return nil, err
	}

	return client, nil
}

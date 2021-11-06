package testmongo

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateTestMongoConnection(ctx context.Context, t *testing.T) (coll *mongo.Collection, teardown func()) {
	url := getMongoUrl(t)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(url))
	if err != nil {
		t.Errorf("An error occurred during creation of mongo connection: %v", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		t.Errorf("Cannot connect to mongo: %v", err)
		t.FailNow()
	}

	db := client.Database("test")
	collName := fmt.Sprintf("test_%d", time.Now().Nanosecond())
	coll = db.Collection(collName)

	teardown = func() {
		coll.Drop(ctx)
		client.Disconnect(ctx)
	}

	_ = coll.Drop(ctx)
	return coll, teardown
}

func getMongoUrl(t *testing.T) string {
	adr := os.Getenv("MONGO")
	usr := os.Getenv("MONGO_USER")
	pass := os.Getenv("MONGO_PASS")

	if adr == "" {
		adr = "localhost"
	}

	res := fmt.Sprintf("mongodb://%s:%s@%s:27017", usr, pass, adr)

	return res
}

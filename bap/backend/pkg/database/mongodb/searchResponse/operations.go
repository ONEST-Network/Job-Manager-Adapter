package searchresponse

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type DaoInterface interface {
	GetSearchJobResponse(id string) (*SearchJobResponse, error)
	CreateSearchJobResponse(worker *SearchJobResponse) error
	ListSearchJobResponse(query bson.D) ([]SearchJobResponse, error)
	DeleteSearchJobResponse(worker, name string) error
	UpdateSearchJobResponse(query, update bson.D) error
}

type Dao struct {
	Collection *mongo.Collection
}

func NewSearchJobResponseDao(collection *mongo.Collection) *Dao {
	return &Dao{
		Collection: collection,
	}
}

const dbTimeout = 10 * time.Second

func (d *Dao) GetWorkerProfile(id string) (*SearchJobResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var worker SearchJobResponse
	if err := d.Collection.FindOne(ctx, bson.D{{Key: "id", Value: id}}).Decode(&worker); err != nil {
		return nil, err
	}

	return &worker, nil
}

func (d *Dao) CreateSearchJobResponse(worker *SearchJobResponse) error {
    ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
    defer cancel()
    _, err := d.Collection.InsertOne(ctx, worker)
    return err
}

func (d *Dao) ListSearchJobResponse(query bson.D) ([]SearchJobResponse, error) {
    ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
    defer cancel()

    cursor, err := d.Collection.Find(ctx, query)
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)

    var workers []SearchJobResponse
    if err = cursor.All(ctx, &workers); err != nil {
        return nil, err
    }

    return workers, nil
}

func (d *Dao) DeleteSearchJobResponse(id string) error {
    ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
    defer cancel()

    _, err := d.Collection.DeleteOne(ctx, bson.D{{Key: "id", Value: id}})
    return err
}

func (d *Dao) UpdateSearchJobResponse(query, update bson.D) error {
    ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
    defer cancel()

    _, err := d.Collection.UpdateOne(ctx, query, update)
    return err
}

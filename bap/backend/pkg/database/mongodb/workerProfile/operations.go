package workerProfile

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type DaoInterface interface {
	GetWorkerProfile(id string) (*WorkerProfile, error)
	CreateWorkerProfile(worker *WorkerProfile) error
	ListWorkerProfile(query bson.D) ([]WorkerProfile, error)
	DeleteWorkerProfile(worker, name string) error
	UpdateWorkerProfile(query, update bson.D) error
}

type Dao struct {
	collection *mongo.Collection
}

func NewWorkerDao(collection *mongo.Collection) *Dao {
	return &Dao{
		collection: collection,
	}
}

const dbTimeout = 10 * time.Second

func (d *Dao) GetWorkerProfile(id string) (*WorkerProfile, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var worker WorkerProfile
	if err := d.collection.FindOne(ctx, bson.D{{Key: "id", Value: id}}).Decode(&worker); err != nil {
		return nil, err
	}

	return &worker, nil
}

func (d *Dao) CreateWorkerProfile(worker *WorkerProfile) error {
    ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
    defer cancel()

    _, err := d.collection.InsertOne(ctx, worker)
    return err
}

func (d *Dao) ListWorkerProfile(query bson.D) ([]WorkerProfile, error) {
    ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
    defer cancel()

    cursor, err := d.collection.Find(ctx, query)
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)

    var workers []WorkerProfile
    if err = cursor.All(ctx, &workers); err != nil {
        return nil, err
    }

    return workers, nil
}

func (d *Dao) DeleteWorkerProfile(id string) error {
    ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
    defer cancel()

    _, err := d.collection.DeleteOne(ctx, bson.D{{Key: "id", Value: id}})
    return err
}

func (d *Dao) UpdateWorkerProfile(query, update bson.D) error {
    ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
    defer cancel()

    _, err := d.collection.UpdateOne(ctx, query, update)
    return err
}

package job

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	database "github.com/ONEST-Network/Whatsapp-Chatbot/bpp/backend/pkg/database/mongodb"
)

type DaoInterface interface {
	CreateJob(job *Job) error
	GetJob(jobId string) error
	ListJobs(query bson.D) ([]Job, error)
	DeleteJob(jobID string) error
	UpdateJob(query, update bson.D) error
}

type Dao struct {
	collection *mongo.Collection
}

func NewJobDao(collection *mongo.Collection) *Dao {
	return &Dao{
		collection: collection,
	}
}

const dbTimeout = 10 * time.Second

// CreateJob creates a job post in the database
func (d *Dao) CreateJob(job *Job) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	if _, err := database.Operator.Create(ctx, d.collection, job); err != nil {
		return err
	}

	return nil
}

// GetJob returns the job for the given id
func (d *Dao) GetJob(jobId string) (*Job, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var query = bson.D{{Key: "id", Value: jobId}}

	result := database.Operator.Get(ctx, d.collection, query)
	if result.Err() != nil {
		return nil, result.Err()
	}

	var job Job
	if err := result.Decode(&job); err != nil {
		return nil, err
	}

	return &job, nil
}

// ListJobs lists jobs from the database
func (d *Dao) ListJobs(query bson.D) ([]Job, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	result, err := database.Operator.List(ctx, d.collection, query)
	if err != nil {
		return nil, err
	}

	var jobs []Job
	if err = result.All(ctx, &jobs); err != nil {
		return nil, err
	}

	return jobs, nil
}

// UpdateJob updates a job in the database
func (d *Dao) UpdateJob(query, update bson.D) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	if _, err := database.Operator.Update(ctx, d.collection, query, update); err != nil {
		return err
	}

	return nil
}

// DeleteJob deletes a job from the database
func (d *Dao) DeleteJob(jobID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := bson.D{{Key: "id", Value: jobID}}

	if _, err := database.Operator.Delete(ctx, d.collection, query); err != nil {
		return err
	}

	return nil
}

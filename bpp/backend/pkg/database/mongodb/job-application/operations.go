package jobapplication

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type DaoInterface interface {
	CreateJobApplication(jobApplication *JobApplication) error
	ListJobApplication(query bson.D) ([]JobApplication, error)
	DeleteJobApplication(applicationID, name string) error
	UpdateJobApplication(query, update bson.D) error
}

type Dao struct {
	collection *mongo.Collection
}

func NewJobApplicationDao(collection *mongo.Collection) *Dao {
	return &Dao{
		collection: collection,
	}
}

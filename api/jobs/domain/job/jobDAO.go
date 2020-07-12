package job

import (
	"context"
	jobDb "github.com/nozgurozturk/jobba/api/jobs/database/mongo"
	errors "github.com/nozgurozturk/jobba/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

var (
	collection = jobDb.Client.Database("jobba").Collection("jobs")
)

func (job *Model) Create() *errors.ApplicationError {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, err := collection.InsertOne(ctx, job)
	if err != nil {
		return errors.InternalServer("user can not created")
	}
	job.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (job *Model) FindOne() (*Model, *errors.ApplicationError) {

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	err := collection.FindOne(ctx, bson.M{"_id": job.ID}).Decode(&job)
	if err != nil {
		return nil, errors.InternalServer("job can not found")
	}
	return job, nil
}

func (job *Model) FindAll() ([]Response, *errors.ApplicationError) {
	var jobs []Response
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	cursor, err := collection.Find(ctx, bson.M{"userID": job.UserID})
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var j Response
		if err = cursor.Decode(&j); err != nil {
			return nil, errors.InternalServer("Jobs can not found")
		}
		jobs = append(jobs, j)
	}

	return jobs, nil
}

func (job *Model) Update() (*Model, *errors.ApplicationError) {

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	bsonJob, err := bson.Marshal(&job)
	if err != nil {
		return nil, errors.InternalServer(err.Error())
	}
	result, err := collection.UpdateOne(ctx,
		bson.M{"_id": job.ID},
		bson.D{
			{"$set", bsonJob},
		},
	)
	if err == nil {
		return nil, errors.InternalServer("job can not updated")
	}
	bsonBytes, _ := bson.Marshal(result)
	bson.Unmarshal(bsonBytes, &job)
	return job, nil
}

func (job *Model) Delete() *errors.ApplicationError {

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	_, err := collection.DeleteOne(ctx, bson.D{{"_id", job.ID}})
	if err == nil {
		return errors.InternalServer("job can not deleted")
	}
	return nil
}

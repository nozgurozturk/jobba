package job

import (
	errors "github.com/nozgurozturk/jobba/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type JobStatus int

const (
	Created JobStatus = iota
	Waiting
	Accepted
	Rejected
	Archived
	end
)

type Model struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID      primitive.ObjectID `bson:"userID,omitempty" json:"userID"`
	Title       string             `bson:"title,omitempty" json:"title"`
	Link        string             `bson:"link" json:"link"`
	Company     string             `bson:"company" json:"company"`
	Description string             `bson:"description" json:"description"`
	CreatedAt   int64              `bson:"createdAt" json:"createdAt"`
	UpdatedAt   int64              `bson:"updatedAt" json:"updatedAt"`
	Status      JobStatus          `bson:"status" json:"status"`
}

type Response struct {
	ID          primitive.ObjectID `bson:"_id" json:"id"`
	Title       string             `bson:"title" json:"title"`
	Company     string             `bson:"company" json:"company"`
	Link        string             `bson:"link" json:"link"`
	Description string             `bson:"description" json:"description"`
	CreatedAt   int64              `bson:"createdAt" json:"createdAt"`
	UpdatedAt   int64              `bson:"updatedAt" json:"updatedAt"`
	Status      JobStatus          `bson:"status" json:"status"`
}

type Request struct {
	ID          string    `json:"id" `
	Title       string    `json:"title"`
	Company     string    `json:"company"`
	Link        string    `json:"link"`
	Description string    `json:"description"`
	Status      JobStatus `json:"status"`
}

func (job *Model) Validate() *errors.ApplicationError {
	if job.Title == "" {
		return errors.BadRequest("Title can not be empty")
	}
	if job.Status >= end {
		return errors.BadRequest("Unhandled job status")
	}
	return nil
}

func (job *Model) JobResponse() interface{} {
	return Response{
		ID:          job.ID,
		Title:       job.Title,
		Company:     job.Company,
		Link:        job.Link,
		Description: job.Description,
		CreatedAt:   job.CreatedAt,
		UpdatedAt:   job.UpdatedAt,
		Status:      job.Status,
	}
}

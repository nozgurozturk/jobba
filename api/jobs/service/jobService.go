package service

import (
	"github.com/nozgurozturk/jobba/api/jobs/domain/job"
	errors "github.com/nozgurozturk/jobba/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type jobService struct{}

type jobServiceInterface interface {
	Create(job.Model) (*job.Model, *errors.ApplicationError)
	FindOne(string) (*job.Model, *errors.ApplicationError)
	FindAll(string) (*[]job.Response, *errors.ApplicationError)
	Update(*job.Request) (*job.Model, *errors.ApplicationError)
}

var (
	Job jobServiceInterface = &jobService{}
)

func (s *jobService) Create(j job.Model) (*job.Model, *errors.ApplicationError) {

	if err := j.Validate(); err != nil {
		return nil, err
	}
	now := time.Now().Unix()
	j.CreatedAt = now
	j.UpdatedAt = now
	if err := j.Create(); err != nil {
		return nil, err
	}

	return &j, nil
}

func (s *jobService) FindOne(jobId string) (*job.Model, *errors.ApplicationError) {

	if jobId == "" {
		return nil, errors.BadRequest("job id can not be empty")
	}
	_id, err := primitive.ObjectIDFromHex(jobId)
	if err != nil {
		return nil, errors.InternalServer(err.Error())
	}
	j := &job.Model{ID: _id}
	foundJob, errResult := j.FindOne();
	if errResult != nil {
		return nil, errResult
	}

	return foundJob, nil
}


func (s *jobService) FindAll(userID string) (*[]job.Response, *errors.ApplicationError) {

	if userID == "" {
		return nil, errors.BadRequest("user id can not be empty")
	}
	_userID,err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.InternalServer(err.Error())
	}
	j := &job.Model{UserID: _userID}
	jobs, errResult := j.FindAll();
	if errResult != nil {
		return nil, errResult
	}

	return &jobs, nil
}

func (s *jobService) Update(j *job.Request) (*job.Model, *errors.ApplicationError) {

	if j.ID == "" {
		return nil, errors.BadRequest("job id can not be empty")
	}
	currentJob, err := s.FindOne(j.ID)
	if err != nil {
		return nil, errors.NotFound("job is not found")
	}
	if j.Link != "" {
		currentJob.Link = j.Link
	}

	if j.Status < 5 {
		currentJob.Status = j.Status
	}

	if j.Description != "" {
		currentJob.Description = j.Description
	}

	if j.Company != "" {
		currentJob.Company = j.Company
	}

	if j.Title != "" {
		currentJob.Title = j.Title
	}

	currentJob.UpdatedAt = time.Now().Unix()

	updatedJob, err := currentJob.Update()
	if err != nil {
		return nil, err
	}

	return updatedJob, nil
}




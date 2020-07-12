package user

import (
	"context"
	userDB "github.com/nozgurozturk/jobba/api/user/database/mongo"
	errors "github.com/nozgurozturk/jobba/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

var (
	collection = userDB.Client.Database("jobba").Collection("users")
)

func (user *Model) Create() *errors.ApplicationError {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, err := collection.InsertOne(ctx, user)
	if err != nil {
		return errors.InternalServer("user can not created")
	}
	user.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (user *Model) FindByName() (*Model, *errors.ApplicationError){

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	err := collection.FindOne(ctx, bson.D{{"userName", user.UserName}}).Decode(&user)
	if err != nil {
		return nil, errors.InternalServer("user can not find")
	}
	return user, nil
}


func (user *Model) FindByUserID() (*Model, *errors.ApplicationError){

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	err := collection.FindOne(ctx, bson.D{{"_id", user.ID}}).Decode(&user)
	if err != nil {
		return nil, errors.InternalServer("user can not find")
	}
	return user, nil
}

func (user *Model) FindByEmail() (*Model, *errors.ApplicationError){

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	err := collection.FindOne(ctx, bson.D{{"email", user.Email}}).Decode(&user)
	if err != nil {
		return nil, errors.InternalServer("user can not find")
	}
	return user, nil
}

func (user *Model) Update() (*Model, *errors.ApplicationError){

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	err := collection.FindOne(ctx, bson.D{{"userName", user.UserName}}).Decode(&user)
	if err == nil {
		return nil, errors.InternalServer("user can not find")
	}
	return user, nil
}

func (user *Model) Delete() (*errors.ApplicationError){

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	_, err := collection.DeleteOne(ctx, bson.D{{"userName", user.UserName}})
	if err == nil {
		return errors.InternalServer("user can not deleted")
	}
	return nil
}
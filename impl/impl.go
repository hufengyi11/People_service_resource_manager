package impl

import (
	"context"
	"fmt"
	"log"
	"os"

	gen "github.com/hufengyi11/People_service_resource_manager/gen/go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserServiceServerImpl struct {
	gen.UnimplementedUserServiceServer
}

var db *mongo.Client
var projectdb *mongo.Collection

type PeopleDetail struct {
	ID   primitive.ObjectID `bson:"_id,omitempty"`
	Name string             `bson:"name"`
}

func mongoNewClient() (*mongo.Client, *mongo.Collection, error) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(os.Getenv("MONGODB_URI")))
	if err != nil {
		return nil, nil, err
	}

	resourceManagerDB := client.Database("ResourceManagement")
	projectsCollection := resourceManagerDB.Collection("People_Service")

	return client, projectsCollection, nil
}

func (u *UserServiceServerImpl) CreateUser(ctx context.Context, req *gen.CreateUserReq) (*gen.CreateUserRes, error) {

	user := req.GetUser()

	client, collection, err := mongoNewClient()
	if err != nil {
		return nil, err
	}
	defer client.Disconnect(ctx)

	res, reserr := collection.InsertOne(ctx, bson.D{
		{Key: "Name", Value: user.Name},
	})
	if reserr != nil {
		log.Fatal(reserr)
	}
	fmt.Println(res.InsertedID)

	return &gen.CreateUserRes{User: user}, err
}

func (u *UserServiceServerImpl) ReadUser(ctx context.Context, req *gen.ReadUserReq) (*gen.ReadUserRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReadUser not implemented")
}

func (u *UserServiceServerImpl) UpdateUser(ctx context.Context, req *gen.UpdateUserReq) (*gen.UpdateUserRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateUser not implemented")
}

func (u *UserServiceServerImpl) DeleteUser(ctx context.Context, req *gen.DeleteUserReq) (*gen.DeleteUserRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteUser not implemented")
}

func (u *UserServiceServerImpl) ListUsers(ctx context.Context, req *gen.ListUsersReq) (*gen.ListUsersRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListUsers not implemented")
}

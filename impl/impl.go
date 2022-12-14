package impl

import (
	"context"
	"fmt"
	"log"

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

type UserDetail struct {
	ID   primitive.ObjectID `bson:"_id,omitempty"`
	Name string             `bson:"name"`
}

func mongoNewClient() (*mongo.Client, *mongo.Collection, error) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb+srv://dbUser:dbUserPassword@cluster0.kjucuqb.mongodb.net/?retryWrites=true&w=majority"))
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

	return &gen.CreateUserRes{User: user}, nil
}

func (u *UserServiceServerImpl) ReadUser(ctx context.Context, req *gen.ReadUserReq) (*gen.ReadUserRes, error) {
	userName := req.GetName()

	client, collection, err := mongoNewClient()
	if err != nil {
		return nil, err
	}
	defer client.Disconnect(ctx)

	findUser := collection.FindOne(ctx, bson.M{"Name": userName})

	data := UserDetail{}

	if error := findUser.Decode(&data); error != nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("Could not find user with name %s: %v", req.GetName(), err))
	}

	readResponse := &gen.ReadUserRes{
		User: &gen.User{
			Id:   data.ID.Hex(),
			Name: data.Name,
		},
	}

	return readResponse, nil

}

func (u *UserServiceServerImpl) UpdateUser(ctx context.Context, req *gen.UpdateUserReq) (*gen.UpdateUserRes, error) {
	user := req.GetUser()
	oid, err := primitive.ObjectIDFromHex(user.GetId())
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Couldn't convert the project id to a Object Name: %v", err),
		)
	}
	client, collection, err := mongoNewClient()
	if err != nil {
		return nil, err
	}
	defer client.Disconnect(ctx)

	update := bson.M{
		"Name": user.GetName(),
	}

	filter := bson.M{
		"_id": oid,
	}

	result := collection.FindOneAndUpdate(ctx, filter, bson.M{"$set": update}, options.FindOneAndUpdate().SetReturnDocument(1))

	decoded := UserDetail{}
	err = result.Decode(&decoded)

	if err != nil {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("Could not find user with name %v", err),
		)
	}

	return &gen.UpdateUserRes{
		User: &gen.User{
			Id:   decoded.ID.Hex(),
			Name: decoded.Name,
		},
	}, nil

}

func (u *UserServiceServerImpl) DeleteUser(ctx context.Context, req *gen.DeleteUserReq) (*gen.DeleteUserRes, error) {
	oid, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Could not convert to ObjectId: %v", err))
	}
	client, collection, err := mongoNewClient()
	if err != nil {
		return nil, err
	}
	defer client.Disconnect(ctx)

	deleteresult, deleteerr := collection.DeleteOne(ctx, bson.M{"_id": oid})

	fmt.Printf("%v", deleteresult)

	if deleteerr != nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("Could not find/delete blog with id %s: %v", req.GetId(), err))
	}

	return &gen.DeleteUserRes{
		Success: true,
	}, nil

}

func (u *UserServiceServerImpl) ListUsers(context.Context, *gen.ListUsersReq) (*gen.ListUsersRes, error) {
	data := &UserDetail{}

	// list user needs new client instead of new connect
	client, _ := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://dbUser:dbUserPassword@cluster0.kjucuqb.mongodb.net/?retryWrites=true&w=majority"))
	defer client.Disconnect(context.Background())

	// resourceManagerDB := client.Database("ResourceManagement")
	collection := client.Database("ResourceManagement").Collection("People_Service")

	client.Connect(context.Background())

	cursor, _ := collection.Find(context.Background(), bson.M{})
	defer cursor.Close(context.Background())

	list := []*gen.User{}

	for cursor.Next(context.Background()) {
		cursor.Decode(data)
		list = append(list, &gen.User{
			Id:   data.ID.Hex(),
			Name: data.Name,
		})
	}

	return &gen.ListUsersRes{
		Users: list,
	}, nil
}

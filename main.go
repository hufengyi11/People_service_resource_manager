package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	gen "github.com/hufengyi11/People_service_resource_manager/gen/go"
	impl "github.com/hufengyi11/People_service_resource_manager/impl"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var db *mongo.Client
var projectdb *mongo.Collection
var mongoCtx context.Context

// mongodb+srv://dbUser:dbUserPassword@cluster0.kjucuqb.mongodb.net/?retryWrites=true&w=majority
func main() {

	opts := []grpc.ServerOption{}
	var grpcServer = grpc.NewServer(opts...)
	var projectServer = impl.UserServiceServerImpl{}
	gen.RegisterUserServiceServer(grpcServer, &projectServer)

	address := "localhost:8080"
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("could not listen to %s %v !!!!!!", address, err)
	}

	fmt.Println("successfully connect to grpc at port:", address)

	reflection.Register(grpcServer)

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	uri := os.Getenv("MONGODB_URI")

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	// client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		fmt.Printf("Connect Error: %v \n", err)
	}
	defer client.Disconnect(context.Background())

	if err := client.Ping(context.Background(), readpref.Primary()); err != nil {
		fmt.Printf("Ping Error: %v \n", err)
	}

	// database := client.Database("ResourceManagement")
	// projectsCollection := database.Collection("People_Service")
	// fmt.Println(projectsCollection.Find(context.Background(), bson.D{{}}))

	log.Fatalln(grpcServer.Serve(lis))

	// ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	// eerr := client.Connect(ctx)
	// if eerr != nil {
	// 	log.Fatal(eerr)
	// }
	// defer client.Disconnect(ctx)
}

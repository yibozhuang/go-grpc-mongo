package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	blogpb "github.com/yibozhuang/go-grpc-mongo/proto"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

/* global variables */
type BlogServiceServer struct{}

type BlogItem struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	AuthorID string             `bson:"author_id"`
	Content  string             `bson:"content"`
	Title    string             `bson:"title"`
}

var db *mongo.Client
var blogdb *mongo.Collection
var mongoCtx context.Context

func (s *BlogServiceServer) ReadBlog(ctx context.Context, req *blogpb.ReadBlogReq) (*blogpb.ReadBlogRes, error) {
	// first convert to mongo object ID
	id, err := primitive.ObjectIDFromHex(req.GetId())

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Could not convert to ObjectId - %v", err),
		)
	}

	result := blogdb.FindOne(ctx, bson.M{"_id": id})

	data := BlogItem{}

	if err := result.Decode(&data); err != nil {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("Could not find blog with ObjectId %s - %v", req.GetId(), err),
		)
	}

	response := &blogpb.ReadBlogRes{
		Blog: &blogpb.Blog{
			Id:       id.Hex(),
			AuthorId: data.AuthorID,
			Title:    data.Title,
			Content:  data.Content,
		},
	}

	return response, nil
}

func (s *BlogServiceServer) CreateBlog(ctx context.Context, req *blogpb.CreateBlogReq) (*blogpb.CreateBlogRes, error) {
	blog := req.GetBlog()

	data := BlogItem{
		// no ID
		AuthorID: blog.GetAuthorId(),
		Title:    blog.GetTitle(),
		Content:  blog.GetContent(),
	}

	result, err := blogdb.InsertOne(ctx, data)

	if err != nil {
		// grpc status, grpc codes
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal error - %v", err),
		)
	}

	id := result.InsertedID.(primitive.ObjectID)
	blog.Id = id.Hex()

	return &blogpb.CreateBlogRes{Blog: blog}, nil
}

func (s *BlogServiceServer) UpdateBlog(ctx context.Context, req *blogpb.UpdateBlogReq) (*blogpb.UpdateBlogRes, error) {
	blog := req.GetBlog()

	id, err := primitive.ObjectIDFromHex(blog.GetId())

	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Could bot convert input to valid Blog object - %v", err),
		)
	}

	update := bson.M{
		"author_id": blog.GetAuthorId(),
		"title":      blog.GetTitle(),
		"content":    blog.GetContent(),
	}

	filter := bson.M{"_id": id}

	result := blogdb.FindOneAndUpdate(ctx, filter, bson.M{"$set": update}, options.FindOneAndUpdate().SetReturnDocument(1))

	data := BlogItem{}
	err = result.Decode(&data)
	if err != nil {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("Could not find blog with ObjectId %s - %v", blog.GetId(), err),
		)
	}

	response := &blogpb.UpdateBlogRes{
		Blog: &blogpb.Blog{
			Id:       data.ID.Hex(),
			AuthorId: data.AuthorID,
			Title:    data.Title,
			Content:  data.Content,
		},
	}

	return response, nil
}

func (s *BlogServiceServer) DeleteBlog(ctx context.Context, req *blogpb.DeleteBlogReq) (*blogpb.DeleteBlogRes, error) {
	// first convert to mongo object ID
	id, err := primitive.ObjectIDFromHex(req.GetId())

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Could not convert to ObjectId - %v", err),
		)
	}

	_, err = blogdb.DeleteOne(ctx, bson.M{"_id": id})

	if err != nil {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("Could not find or delete blog with ObjectId %s - %v", req.GetId(), err),
		)
	}

	response := &blogpb.DeleteBlogRes{
		Success: true,
	}

	return response, nil
}

func (s *BlogServiceServer) ListBlogs(req *blogpb.ListBlogsReq, stream blogpb.BlogService_ListBlogsServer) error {
	data := &BlogItem{}

	cursor, err := blogdb.Find(context.Background(), bson.M{})

	if err != nil {
		return status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal error - %v", err),
		)
	}

	// defer closing the database until end of function
	defer cursor.Close(context.Background())

	// leep looping until we don't have anymore data
	for cursor.Next(context.Background()) {
		err := cursor.Decode(data)

		if err != nil {
			return status.Errorf(
				codes.Unavailable,
				fmt.Sprintf("Error decoding data - %v", err),
			)
		}

		stream.Send(&blogpb.ListBlogsRes{
			Blog: &blogpb.Blog{
				Id:       data.ID.Hex(),
				AuthorId: data.AuthorID,
				Content:  data.Content,
				Title:    data.Title,
			},
		})
	}

	if err := cursor.Err(); err != nil {
		return status.Errorf(
			codes.Internal,
			fmt.Sprintf("Unkown cursor error - %v", err),
		)
	}
	return nil
}

func main() {
	grpc_port := ":50051"
	mongdb_url := "mongodb://localhost:27017"

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	fmt.Println("Starting server on port %s", grpc_port)

	listener, err := net.Listen("tcp", grpc_port)

	if err != nil {
		log.Fatalf("Unable to listen on port %s - %v", grpc_port, err)
	}

	opts := []grpc.ServerOption{}
	s := grpc.NewServer(opts...)
	srv := &BlogServiceServer{}

	blogpb.RegisterBlogServiceServer(s, srv)

	// initialize mongodb client
	fmt.Println("Connecting to mongo...")
	mongoCtx = context.Background()
	db, err = mongo.Connect(mongoCtx, options.Client().ApplyURI(mongdb_url))

	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping(mongoCtx, nil)

	if err != nil {
		log.Fatalf("Could not connect to mongo - %v", err)
	} else {
		fmt.Println("Connected to mongo")
	}

	blogdb = db.Database("mydb").Collection("blog")

	// Start the server in a child routine
	go func() {
		if err := s.Serve(listener); err != nil {
			log.Fatalf("Failed to serve - %v", err)
		}
	}()
	fmt.Println("Server succesfully started on port %s", grpc_port)

	// create a channel to receive OS signals
	c := make(chan os.Signal)

	// relay os.Interrupt to channel (os.Interrupt = CTRL+C)
	signal.Notify(c, os.Interrupt)

	// block main routine until a signal is received
	// as long as user doesn't press CTRL+C, our main routine keeps running
	<-c

	// After receiving CTRL+C Properly stop the server
	fmt.Println("Stopping the server...")
	s.Stop()
	listener.Close()

	fmt.Println("Closing mongo connection...")
	db.Disconnect(mongoCtx)

	fmt.Println("END")
}

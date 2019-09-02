package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/crojas/go-grpc-course/blog/blogpb"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

type server struct{}

type blogItem struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	AuthorID string             `bson:"author_id"`
	Title    string             `bson:"title"`
	Content  string             `bson:"content"`
}

var collection *mongo.Collection

func main() {
	// Si el programa se cae, nos dice el archivo y linea donde crasheo
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Conectar a mongodb
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatalf("error al conectarse a mongodb %v", err)
	}
	collection = client.Database("mydb").Collection("blog")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("error al escuchar el puerto 50051 %v", err)
	}
	s := grpc.NewServer()
	blogpb.RegisterBlogServiceServer(s, &server{})

	go func() {
		fmt.Println("Levantando el servidor...")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Error al servir %v", err)
		}
	}()

	// Crea un canal para bloquear hasta que se ingrese ctrl+c
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch

	fmt.Println(("Deteniendo el servidor..."))
	s.Stop()
	fmt.Println("Cerrando listener")
	lis.Close()
	fmt.Println("Fin del programa")
}
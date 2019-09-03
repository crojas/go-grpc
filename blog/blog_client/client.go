package main

import (
	"context"
	"fmt"
	"log"

	"github.com/crojas/go-grpc-course/blog/blogpb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Cliente Blog GRPC")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("No me pude conectar %v", err)
	}
	defer cc.Close()
	c := blogpb.NewBlogServiceClient(cc)
	blog := &blogpb.Blog{
		AuthorId: "Carlos",
		Title:    "First Blog",
		Content:  "Esta es una prueba de GRPC and mondodb",
	}
	res, err := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{Blog: blog})
	if err != nil {
		log.Fatalf("Error inesperado %v", err)
	}
	fmt.Printf("Blog creado %v", res.GetBlog())
}

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
		AuthorId: "Marianela",
		Title:    "Segundo Blog",
		Content:  "Segundo insert",
	}
	createBlogRes, err := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{Blog: blog})
	if err != nil {
		log.Fatalf("Error inesperado %v", err)
	}
	fmt.Printf("Blog creado %v", createBlogRes.GetBlog())
	blogID := createBlogRes.GetBlog().GetId()

	// Llamada para leer pasando el id
	fmt.Println("Reading the blog")

	_, err2 := c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{BlogId: "sdfsdfsd"})
	if err2 != nil {
		fmt.Printf("Error leyendo con id: sdfsdfsd %v \n", err2)
	}

	readBlogReq := &blogpb.ReadBlogRequest{BlogId: blogID}
	readBlogRes, readBlogErr := c.ReadBlog(context.Background(), readBlogReq)
	if readBlogErr != nil {
		fmt.Printf("Error leyendo: %v \n", readBlogErr)
	}

	fmt.Printf("Blog leido: %v \n", readBlogRes)

	// Actualizar Blog
	newBlog := &blogpb.Blog{
		Id:       blogID,
		AuthorId: "Changed Author",
		Title:    "My First Blog (edited)",
		Content:  "Content of the first blog, with some awesome additions!",
	}
	updateRes, updateErr := c.UpdateBlog(context.Background(), &blogpb.UpdateBlogRequest{Blog: newBlog})
	if updateErr != nil {
		fmt.Printf("Error actualizando: %v \n", updateErr)
	}
	fmt.Printf("Blog fue actualizado: %v\n", updateRes)
}

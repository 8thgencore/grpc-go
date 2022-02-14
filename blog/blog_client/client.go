package main

import (
	"context"
	"fmt"
	"log"

	"github.com/8thgencore/grpc-go/blog/blogpb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Blog Client")

	opts := grpc.WithInsecure()

	conn, err := grpc.Dial("localhost:50051", opts)
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer conn.Close()

	c := blogpb.NewBlogServiceClient(conn)

	// CreateBlog(c)
	// ReadBlog(c)
	UpdateBlog(c)
}

func CreateBlog(c blogpb.BlogServiceClient) {
	blog := &blogpb.Blog{
		AuthorId: "Stephane",
		Title:    "My First Blog",
		Content:  "Content of the first blog",
	}
	createBlogRes, err := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{
		Blog: blog,
	})
	if err != nil {
		log.Fatalf("Unexpected error: %v", err)
	}
	fmt.Printf("Blog has been created: %v", createBlogRes)
}

func ReadBlog(c blogpb.BlogServiceClient) {
	blogID := "62065deeda61a5de179fba8b"

	_, err := c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{
		BlogId: blogID,
	})
	if err != nil {
		log.Fatalf("Error happened while reading: %v\n", err)
	}

	readBlogReq := &blogpb.ReadBlogRequest{BlogId: blogID}
	readBlogRes, readBlogErr := c.ReadBlog(context.Background(), readBlogReq)
	if readBlogErr != nil {
		fmt.Printf("Error happened while reading: %v\n", readBlogErr)
	}

	fmt.Printf("Blog was read: %v\n", readBlogRes)
}

func UpdateBlog(c blogpb.BlogServiceClient) {
	blogID := "62065deeda61a5de179fba8b"

	newBlog := &blogpb.Blog{
		Id:       blogID,
		AuthorId: "Changed Author",
		Title:    "My First Blog (edited)",
		Content:  "Content of the first blog, with some awesome additions!",
	}

	updateBlogReq := &blogpb.UpdateBlogRequest{Blog: newBlog}
	updateRes, updateErr := c.UpdateBlog(context.Background(), updateBlogReq)
	if updateErr != nil {
		fmt.Printf("Error happened while reading: %v\n", updateErr)
	}

	fmt.Printf("Blog was read: %v\n", updateRes)
}

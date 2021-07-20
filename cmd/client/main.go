package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/sakiib/grpc-try/gen/pb"
	"google.golang.org/grpc"
	"log"
	"strconv"
	"time"
)

func main() {
	serverAddress := flag.String("address", "", "the server address")
	flag.Parse()
	log.Printf("server address: %s", *serverAddress)

	conn, err := grpc.Dial(*serverAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to dial with %s", err.Error())
	}
	defer conn.Close()

	client := pb.NewBookServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	for i := 0; i < 5; i++ {
		book := &pb.Book{
			Id:   strconv.Itoa(i),
			Name: fmt.Sprintf("book-%d", i),
		}
		res, err := client.CreateBook(ctx, &pb.CreateBookRequest{
			Book: book,
		})
		if err != nil {
			log.Printf("failed to create book with %s", err)
		}
		log.Println("create book response: ", res)
	}

	book, err := client.GetBook(ctx, &pb.GetBookRequest{
		Id: "3",
	})
	if err != nil {
		log.Printf("failed to get book with %s", err.Error())
	}
	log.Println("book: ", book)

	books, err := client.GetBooks(ctx, &pb.EmptyRequest{})
	if err != nil {
		log.Printf("failed to get books with %s", err.Error())
	} else {
		log.Println("books list: ", books)
	}
}

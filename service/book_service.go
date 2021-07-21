package service

import (
	"context"
	"fmt"
	"github.com/sakiib/grpc-try/gen/pb"
	"io"
	"log"
)

type BookService struct {
	pb.UnimplementedBookServiceServer
	store BookStore
}

func NewBookService(store BookStore) *BookService {
	return &BookService{pb.UnimplementedBookServiceServer{}, store}
}

func (bs *BookService) CreateBook(ctx context.Context, req *pb.CreateBookRequest) (*pb.CreateBookResponse, error) {
	book := req.GetBook()
	log.Printf("book details: %+v, %+v", book.Id, book.Name)

	err := bs.store.Set(book)
	if err != nil {
		return nil, fmt.Errorf("failed to set book with %s", err.Error())
	}

	res := &pb.CreateBookResponse{
		Id: book.Id,
	}
	return res, nil
}

func (bs *BookService) GetBook(ctx context.Context, req *pb.GetBookRequest) (*pb.GetBookResponse, error) {
	id := req.GetId()
	log.Printf("get request details: %+v", id)

	book, err := bs.store.Get(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get book with %s", err.Error())
	}

	res := &pb.GetBookResponse{
		Book: book,
	}
	return res, nil
}

func (bs *BookService) GetBooks(ctx context.Context, req *pb.EmptyRequest) (*pb.GetBooksResponse, error) {
	books := bs.store.GetAll()
	res := &pb.GetBooksResponse{
		Book: books,
	}
	return res, nil
}

func (bs *BookService) ListBooks(req *pb.EmptyRequest, stream pb.BookService_ListBooksServer) error {
	books := bs.store.GetAll()
	for _, book := range books {
		if err := stream.Send(&pb.GetBookResponse{
			Book: book,
		}); err != nil {
			return err
		}
	}
	return nil
}

func (bs *BookService) BooksSummary(stream pb.BookService_BooksSummaryServer) error {
	var summary string
	for {
		id, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.SummaryResponse{
				Summary: summary,
			})
		}
		if err != nil {
			return err
		}
		book, err := bs.store.Get(id.GetId())
		if err != nil {
			return err
		}
		summary += book.Name
	}
}

package service

import (
	"acme/internal/repo"
	"context"
	"database/sql"
)

type BookService struct {
	Query *repo.Queries
}

func NewBookService(q *repo.Queries) *BookService {
	return &BookService{Query: q}
}

func (s *BookService) CreateBook(ctx context.Context, arg repo.CreateBookParams) (repo.Book, error) {
	return s.Query.CreateBook(ctx, arg)
}

func (s *BookService) DeleteBook(ctx context.Context, id int64) (error) {
	return s.Query.DeleteBook(ctx, id)
}

func (s *BookService) GetBook(ctx context.Context, id int64) (repo.Book, error) {
	return s.Query.GetBook(ctx, id)
}

func (s *BookService) ListBook(ctx context.Context, arg repo.ListBookParams) ([]repo.Book, error) {
	return s.Query.ListBook(ctx, arg)
}

func (s *BookService) UpdateBook(ctx context.Context, arg repo.UpdateBookParams) (repo.Book, error) {
	return s.Query.UpdateBook(ctx, arg)
}


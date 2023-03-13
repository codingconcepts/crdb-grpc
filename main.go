package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/codingconcepts/crdb-grpc/pb"
	"github.com/jackc/pgx/v5/pgxpool"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("error starting tcp listener: %v", err)
	}

	opts := []grpc.ServerOption{}
	s := grpc.NewServer(opts...)

	db, err := pgxpool.New(context.Background(), "postgres://root@localhost:26257/defaultdb?sslmode=disable")
	if err != nil {
		log.Fatalf("error connecting to db: %v", err)
	}
	defer db.Close()

	server := &server{
		db: db,
	}

	pb.RegisterTodoServiceServer(s, server)
	reflection.Register(s)

	log.Println("ready")
	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

type server struct {
	pb.UnimplementedTodoServiceServer

	db *pgxpool.Pool
}

func (s *server) CreateTodo(ctx context.Context, req *pb.CreateTodoRequest) (*pb.CreateTodoResponse, error) {
	const stmt = `INSERT INTO todo (title) VALUES ($1) RETURNING id`

	row := s.db.QueryRow(ctx, stmt, req.Title)

	var id string
	if err := row.Scan(&id); err != nil {
		return nil, fmt.Errorf("creating todo: %w", err)
	}

	return &pb.CreateTodoResponse{
		Todo: &pb.Todo{
			Id:    id,
			Title: req.Title,
		},
	}, nil
}

func (s *server) GetTodo(ctx context.Context, req *pb.GetTodoRequest) (*pb.GetTodoResponse, error) {
	const stmt = `SELECT title FROM todo WHERE id = $1`

	row := s.db.QueryRow(ctx, stmt, req.Id)

	var title string
	if err := row.Scan(&title); err != nil {
		return nil, fmt.Errorf("getting todo: %w", err)
	}

	return &pb.GetTodoResponse{
		Todo: &pb.Todo{
			Id:    req.Id,
			Title: title,
		},
	}, nil
}

func (s *server) GetTodos(ctx context.Context, req *pb.Empty) (*pb.GetTodosResponse, error) {
	const stmt = `SELECT id, title FROM todo`

	rows, err := s.db.Query(ctx, stmt)
	if err != nil {
		return nil, fmt.Errorf("querying: %w", err)
	}

	var todos []*pb.Todo
	for rows.Next() {
		t := pb.Todo{}
		if err = rows.Scan(&t.Id, &t.Title); err != nil {
			return nil, fmt.Errorf("getting todo: %w", err)
		}

		todos = append(todos, &t)
	}

	return &pb.GetTodosResponse{
		Todos: todos,
	}, nil
}

func (s *server) DeleteTodo(ctx context.Context, req *pb.DeleteTodoRequest) (*pb.DeleteTodoResponse, error) {
	const stmt = `DELETE FROM todo WHERE id = $1`

	result, err := s.db.Exec(ctx, stmt, req.Id)
	if err != nil {
		return nil, fmt.Errorf("deleting: %v", err)
	}

	return &pb.DeleteTodoResponse{
		Affected: result.RowsAffected(),
	}, nil
}

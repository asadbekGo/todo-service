package service

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/asadbekGo/todo-service/genproto"
	l "github.com/asadbekGo/todo-service/pkg/logger"
	"github.com/asadbekGo/todo-service/storage"
)

// TodoService ...
type TodoService struct {
	storage storage.IStorage
	logger  l.Logger
}

// NewTodoService ...
func NewTodoService(db *sqlx.DB, log l.Logger) *TodoService {
	return &TodoService{
		storage: storage.NewStoragePg(db),
		logger:  log,
	}
}

func (s *TodoService) Create(ctx context.Context, req *pb.Todo) (*pb.Todo, error) {
	user, err := s.storage.Todo().Create(*req)
	if err != nil {
		s.logger.Error("falied to create todo", l.Error(err))
		return nil, status.Error(codes.Internal, "failed to create todo")
	}

	return &user, nil
}

func (s *TodoService) Get(ctx context.Context, req *pb.ByIdReq) (*pb.Todo, error) {
	user, err := s.storage.Todo().Get(req.Id)
	if err != nil {
		s.logger.Error("failed to get todo", l.Error(err))
		return nil, status.Error(codes.Internal, "failed to get todo")
	}

	return &user, nil
}

func (s *TodoService) List(ctx context.Context, req *pb.ListReq) (*pb.ListResp, error) {
	todos, count, err := s.storage.Todo().List(req.Page, req.Limit)
	if err != nil {
		s.logger.Error("failed to list todo", l.Error(err))
		return nil, status.Error(codes.Internal, "failed to list todo")
	}

	return &pb.ListResp{
		Todos: todos,
		Count: count,
	}, nil
}

func (s *TodoService) Update(ctx context.Context, req *pb.Todo) (*pb.Todo, error) {
	todo, err := s.storage.Todo().Update(*req)
	if err != nil {
		s.logger.Error("failed to update todo", l.Error(err))
		return nil, status.Error(codes.Internal, "failed to update todo")
	}

	return &todo, nil
}

func (s *TodoService) Delete(ctx context.Context, req *pb.ByIdReq) (*pb.Empty, error) {
	err := s.storage.Todo().Delete(req.Id)
	if err != nil {
		s.logger.Error("failed to delete todo", l.Error(err))
		return nil, status.Error(codes.Internal, "failed to delete todo")
	}

	return &pb.Empty{}, nil
}

func (s *TodoService) ListOverdue(ctx context.Context, req *pb.Time) (*pb.ListOverdueResp, error) {
	layoutISO := "2006-01-02"
	time, err := time.Parse(layoutISO, req.ToTime)
	fmt.Println(time)
	if err != nil {
		s.logger.Error("failed to time parse", l.Error(err))
	}
	todos, err := s.storage.Todo().ListOverdue(time)
	if err != nil {
		s.logger.Error("failed to list overdue todo", l.Error(err))
		return nil, status.Error(codes.Internal, "failed to list overdue todo")
	}

	return &pb.ListOverdueResp{
		Todos: todos,
	}, nil
}

package api

import (
	"context"
	"time"

	"github.com/douglaszuqueto/go-user-microservice/pkg/storage"
	"github.com/douglaszuqueto/go-user-microservice/pkg/util"
	"github.com/douglaszuqueto/go-user-microservice/proto"
	"github.com/google/uuid"

	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UserService UserService
type UserService struct {
	storage storage.UserStorage
}

// NewUserService NewUserService
func NewUserService(s *grpc.Server, storage storage.UserStorage) *UserService {
	server := &UserService{
		storage: storage,
	}

	if s != nil {
		proto.RegisterUserServiceServer(s, server)
	}

	return server
}

// Get Get
func (s *UserService) Get(ctx context.Context, req *proto.GetUserRequest) (*proto.GetUserResponse, error) {
	if _, err := uuid.Parse(req.Id); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	user, err := s.storage.GetUser(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	protoUser, _ := userToProtoStruct(user)

	resp := &proto.GetUserResponse{
		User: &protoUser,
	}

	return resp, nil
}

// List List
func (s *UserService) List(ctx context.Context, req *proto.ListUserRequest) (*proto.ListUserResponse, error) {
	l := []*proto.User{}

	users, err := s.storage.ListUser(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}

	for _, u := range users {
		user, _ := userToProtoStruct(u)

		l = append(l, &user)
	}

	resp := &proto.ListUserResponse{
		User: l,
	}

	return resp, nil
}

// Create Create
func (s *UserService) Create(ctx context.Context, req *proto.CreateUserRequest) (*proto.CreateUserResponse, error) {
	if err := storeValidateOrFail(req.User); err != nil {
		return nil, err
	}

	password, err := util.GeneratePassword(req.User.Password)
	if err != nil {
		return nil, err
	}

	user := storage.User{
		Username:  req.User.Username,
		Password:  password,
		State:     req.User.State,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	id, err := s.storage.CreateUser(ctx, user)

	if err != nil {
		return nil, err
	}

	resp := &proto.CreateUserResponse{
		Id: id,
	}

	return resp, nil
}

// Update Update
func (s *UserService) Update(ctx context.Context, req *proto.UpdateUserRequest) (*proto.UpdateUserResponse, error) {
	userOld, err := s.storage.GetUser(ctx, req.User.Id)
	if err != nil {
		return nil, err
	}

	if err := storeValidateOrFail(req.User); err != nil {
		return nil, err
	}

	password, err := util.GeneratePassword(req.User.Password)
	if err != nil {
		return nil, err
	}

	user := storage.User{
		ID:        req.User.Id,
		Username:  req.User.Username,
		Password:  password,
		State:     req.User.State,
		CreatedAt: userOld.CreatedAt,
		UpdatedAt: time.Now(),
	}

	if err := s.storage.UpdateUser(ctx, user); err != nil {
		return nil, err
	}

	resp := &proto.UpdateUserResponse{
		Result: "ok",
	}

	return resp, nil
}

// Delete Delete
func (s *UserService) Delete(ctx context.Context, req *proto.DeleteUserRequest) (*proto.DeleteUserResponse, error) {
	if _, err := s.storage.GetUser(ctx, req.Id); err != nil {
		return nil, err
	}

	if err := s.storage.DeleteUser(ctx, req.Id); err != nil {
		return nil, err
	}

	resp := &proto.DeleteUserResponse{
		Result: "ok",
	}

	return resp, nil
}

func userToProtoStruct(u storage.User) (proto.User, error) {
	user := proto.User{
		Id:       u.ID,
		Username: u.Username,
		Password: u.Password,
		State:    u.State,
	}

	var err error

	user.CreatedAt, err = ptypes.TimestampProto(u.CreatedAt)
	if err != nil {
		return user, err
	}

	user.UpdatedAt, err = ptypes.TimestampProto(u.UpdatedAt)
	if err != nil {
		return user, err
	}

	return user, nil
}

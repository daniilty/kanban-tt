package server

import (
	"context"

	"github.com/daniilty/kanban-tt/schema"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (g *GRPC) AddUser(ctx context.Context, req *schema.AddUserRequest) (*schema.AddUserResponse, error) {
	err := g.service.AddUser(ctx, convertPBAddUserToCore(req))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &schema.AddUserResponse{}, nil
}

func (g *GRPC) GetUser(ctx context.Context, req *schema.GetUserRequest) (*schema.GetUserResponse, error) {
	user, ok, err := g.service.GetUser(ctx, req.GetId())
	if err != nil {
		if ok {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &schema.GetUserResponse{
		User: convertCoreUserToPB(user),
	}, nil
}

func (g *GRPC) IsUserWithEmailExists(ctx context.Context, req *schema.IsUserWithEmailExistsRequest) (*schema.IsUserWithEmailExistsResponse, error) {
	exists, err := g.service.IsUserWithEmailExists(ctx, req.GetEmail())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &schema.IsUserWithEmailExistsResponse{
		Exists: exists,
	}, nil
}

func (g *GRPC) IsValidUserCredentials(ctx context.Context, req *schema.IsValidUserCredentialsRequest) (*schema.IsValidUserCredentialsResponse, error) {
	isValid, err := g.service.IsValidUserCredentials(ctx, req.GetEmail(), req.GetPasswordHash())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &schema.IsValidUserCredentialsResponse{
		IsValid: isValid,
	}, nil
}

func (g *GRPC) UpdateUser(ctx context.Context, req *schema.UpdateUserRequest) (*schema.UpdateUserResponse, error) {
	err := g.service.UpdateUser(ctx, convertPBUpdateUserToCore(req))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &schema.UpdateUserResponse{}, nil
}

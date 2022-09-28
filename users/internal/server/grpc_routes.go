package server

import (
	"context"
	"errors"
	"strconv"

	"github.com/daniilty/kanban-tt/schema"
	"github.com/daniilty/kanban-tt/users/internal/core"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (g *GRPC) AddUser(ctx context.Context, req *schema.AddUserRequest) (*schema.AddUserResponse, error) {
	id, err := g.service.AddUser(ctx, convertPBAddUserToCore(req))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &schema.AddUserResponse{
		Id: int64(id),
	}, nil
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

func (g *GRPC) GetUserTaskTTL(ctx context.Context, req *schema.GetUserTaskTTLRequest) (*schema.GetUserTaskTTLResponse, error) {
	ttl, err := g.service.GetUserTaskTTL(ctx, req.GetId())
	if err != nil {
		if errors.Is(err, core.ErrNoSuchUser) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &schema.GetUserTaskTTLResponse{
		TaskTtl: int64(ttl),
	}, nil
}

func (g *GRPC) GetTTLs(ctx context.Context, req *schema.GetTTLsRequest) (*schema.GetTTLsResponse, error) {
	ttls := g.service.GetTTLs()

	return &schema.GetTTLsResponse{
		Ttls: ttls,
	}, nil
}

func (g *GRPC) GetUserByEmail(ctx context.Context, req *schema.GetUserByEmailRequest) (*schema.GetUserByEmailResponse, error) {
	user, ok, err := g.service.GetUserByEmail(ctx, req.GetEmail())
	if err != nil {
		if ok {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &schema.GetUserByEmailResponse{
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
		if errors.Is(err, core.ErrNoSuchTTL) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &schema.UpdateUserResponse{}, nil
}

func (g *GRPC) UnconfirmUserEmail(ctx context.Context, req *schema.UnconfirmUserEmailRequest) (*schema.UnconfirmUserEmailResponse, error) {
	err := g.service.UnconfirmEmail(ctx, strconv.Itoa(int(req.GetId())))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &schema.UnconfirmUserEmailResponse{}, nil
}

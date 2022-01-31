package auth

import (
	"context"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/query"
	auth_pb "github.com/caos/zitadel/pkg/grpc/auth"
	obj_grpc "github.com/caos/zitadel/v2/internal/api/grpc/object"
	user_grpc "github.com/caos/zitadel/v2/internal/api/grpc/user"
)

func (s *Server) ListMyZitadelPermissions(ctx context.Context, _ *auth_pb.ListMyZitadelPermissionsRequest) (*auth_pb.ListMyZitadelPermissionsResponse, error) {
	perms, err := s.query.MyZitadelPermissions(ctx, authz.GetCtxData(ctx).UserID)
	if err != nil {
		return nil, err
	}
	return &auth_pb.ListMyZitadelPermissionsResponse{
		Result: perms.Permissions,
	}, nil
}

func (s *Server) ListMyProjectPermissions(ctx context.Context, _ *auth_pb.ListMyProjectPermissionsRequest) (*auth_pb.ListMyProjectPermissionsResponse, error) {
	ctxData := authz.GetCtxData(ctx)
	userGrantOrgID, err := query.NewUserGrantResourceOwnerSearchQuery(ctxData.OrgID)
	if err != nil {
		return nil, err
	}
	userGrantProjectID, err := query.NewUserGrantProjectIDSearchQuery(ctxData.ProjectID)
	if err != nil {
		return nil, err
	}
	userGrantUserID, err := query.NewUserGrantUserIDSearchQuery(ctxData.UserID)
	if err != nil {
		return nil, err
	}
	userGrant, err := s.query.UserGrant(ctx, userGrantOrgID, userGrantProjectID, userGrantUserID)
	if err != nil {
		return nil, err
	}
	return &auth_pb.ListMyProjectPermissionsResponse{
		Result: userGrant.Roles,
	}, nil
}

func (s *Server) ListMyMemberships(ctx context.Context, req *auth_pb.ListMyMembershipsRequest) (*auth_pb.ListMyMembershipsResponse, error) {
	request, err := ListMyMembershipsRequestToModel(ctx, req)
	if err != nil {
		return nil, err
	}
	response, err := s.query.Memberships(ctx, request)
	if err != nil {
		return nil, err
	}
	return &auth_pb.ListMyMembershipsResponse{
		Result: user_grpc.MembershipsToMembershipsPb(response.Memberships),
		Details: obj_grpc.ToListDetails(
			response.Count,
			response.Sequence,
			response.Timestamp,
		),
	}, nil
}

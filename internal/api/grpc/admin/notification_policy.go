package admin

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/grpc/object"
	policy_grpc "github.com/zitadel/zitadel/internal/api/grpc/policy"
	admin_pb "github.com/zitadel/zitadel/pkg/grpc/admin"
)

func (s *Server) GetNotificationPolicy(ctx context.Context, _ *admin_pb.GetNotificationPolicyRequest) (*admin_pb.GetNotificationPolicyResponse, error) {
	policy, err := s.query.DefaultNotificationPolicy(ctx, true)
	if err != nil {
		return nil, err
	}
	return &admin_pb.GetNotificationPolicyResponse{Policy: policy_grpc.ModelNotificationPolicyToPb(policy)}, nil
}

func (s *Server) UpdateNotificationPolicy(ctx context.Context, req *admin_pb.UpdateNotificationPolicyRequest) (*admin_pb.UpdateNotificationPolicyResponse, error) {
	result, err := s.command.ChangeDefaultNotificationPolicy(ctx, UpdateNotificationPolicyToDomain(req))
	if err != nil {
		return nil, err
	}
	return &admin_pb.UpdateNotificationPolicyResponse{
		Details: object.ChangeToDetailsPb(
			result.Sequence,
			result.ChangeDate,
			result.ResourceOwner,
		),
	}, nil
}

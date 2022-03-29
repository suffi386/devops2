package admin

import (
	"context"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/api/grpc/object"
	admin_pb "github.com/caos/zitadel/pkg/grpc/admin"
)

func (s *Server) ListSecretGenerators(ctx context.Context, req *admin_pb.ListSecretGeneratorsRequest) (*admin_pb.ListSecretGeneratorsResponse, error) {
	queries, err := listSecretGeneratorToModel(req)
	if err != nil {
		return nil, err
	}
	result, err := s.query.SearchSecretGenerators(ctx, queries)
	if err != nil {
		return nil, err
	}
	return &admin_pb.ListSecretGeneratorsResponse{
		Details: object.ToListDetails(result.Count, result.Sequence, result.Timestamp),
	}, nil
}

func (s *Server) GetSecretGenerator(ctx context.Context, req *admin_pb.GetSecretGeneratorRequest) (*admin_pb.GetSecretGeneratorResponse, error) {
	generator, err := s.query.SecretGeneratorByType(ctx, SecretGeneratorTypeToDomain(req.GeneratorType))
	if err != nil {
		return nil, err
	}
	return &admin_pb.GetSecretGeneratorResponse{
		SecretGenerator: SecretGeneratorToPb(generator),
	}, nil
}

func (s *Server) UpdateSecretGenerator(ctx context.Context, req *admin_pb.UpdateSecretGeneratorRequest) (*admin_pb.UpdateSecretGeneratorResponse, error) {
	details, err := s.command.ChangeSecretGeneratorConfig(ctx, authz.GetInstance(ctx).ID, SecretGeneratorTypeToDomain(req.GeneratorType), UpdateSecretGeneratorToConfig(req))
	if err != nil {
		return nil, err
	}
	return &admin_pb.UpdateSecretGeneratorResponse{
		Details: object.ChangeToDetailsPb(
			details.Sequence,
			details.EventDate,
			details.ResourceOwner),
	}, nil
}

func (s *Server) GetSMTPConfig(ctx context.Context, req *admin_pb.GetSMTPConfigRequest) (*admin_pb.GetSMTPConfigResponse, error) {
	smtp, err := s.query.SMTPConfigByAggregateID(ctx, authz.GetInstance(ctx).ID)
	if err != nil {
		return nil, err
	}
	return &admin_pb.GetSMTPConfigResponse{
		SmtpConfig: SMTPConfigToPb(smtp),
	}, nil
}

func (s *Server) UpdateSMTPConfig(ctx context.Context, req *admin_pb.UpdateSMTPConfigRequest) (*admin_pb.UpdateSMTPConfigResponse, error) {
	details, err := s.command.ChangeSMTPConfig(ctx, authz.GetInstance(ctx).ID, UpdateSMTPToConfig(req))
	if err != nil {
		return nil, err
	}
	return &admin_pb.UpdateSMTPConfigResponse{
		Details: object.ChangeToDetailsPb(
			details.Sequence,
			details.EventDate,
			details.ResourceOwner),
	}, nil
}

func (s *Server) UpdateSMTPConfigPassword(ctx context.Context, req *admin_pb.UpdateSMTPConfigPasswordRequest) (*admin_pb.UpdateSMTPConfigPasswordResponse, error) {
	details, err := s.command.ChangeSMTPConfigPassword(ctx, authz.GetInstance(ctx).ID, req.Password)
	if err != nil {
		return nil, err
	}
	return &admin_pb.UpdateSMTPConfigPasswordResponse{
		Details: object.ChangeToDetailsPb(
			details.Sequence,
			details.EventDate,
			details.ResourceOwner),
	}, nil
}

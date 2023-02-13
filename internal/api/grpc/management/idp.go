package management

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	idp_grpc "github.com/zitadel/zitadel/internal/api/grpc/idp"
	object_pb "github.com/zitadel/zitadel/internal/api/grpc/object"
	"github.com/zitadel/zitadel/internal/query"
	mgmt_pb "github.com/zitadel/zitadel/pkg/grpc/management"
)

func (s *Server) GetOrgIDPByID(ctx context.Context, req *mgmt_pb.GetOrgIDPByIDRequest) (*mgmt_pb.GetOrgIDPByIDResponse, error) {
	idp, err := s.query.IDPByIDAndResourceOwner(ctx, true, req.Id, authz.GetCtxData(ctx).OrgID, false)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.GetOrgIDPByIDResponse{Idp: idp_grpc.ModelIDPViewToPb(idp)}, nil
}

func (s *Server) ListOrgIDPs(ctx context.Context, req *mgmt_pb.ListOrgIDPsRequest) (*mgmt_pb.ListOrgIDPsResponse, error) {
	queries, err := listIDPsToModel(ctx, req)
	if err != nil {
		return nil, err
	}
	resp, err := s.query.IDPs(ctx, queries, false)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ListOrgIDPsResponse{
		Result:  idp_grpc.IDPViewsToPb(resp.IDPs),
		Details: object_pb.ToListDetails(resp.Count, resp.Sequence, resp.Timestamp),
	}, nil
}

func (s *Server) AddOrgOIDCIDP(ctx context.Context, req *mgmt_pb.AddOrgOIDCIDPRequest) (*mgmt_pb.AddOrgOIDCIDPResponse, error) {
	config, err := s.command.AddIDPConfig(ctx, AddOIDCIDPRequestToDomain(req), authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.AddOrgOIDCIDPResponse{
		IdpId: config.IDPConfigID,
		Details: object_pb.AddToDetailsPb(
			config.Sequence,
			config.ChangeDate,
			config.ResourceOwner,
		),
	}, nil
}

func (s *Server) AddOrgJWTIDP(ctx context.Context, req *mgmt_pb.AddOrgJWTIDPRequest) (*mgmt_pb.AddOrgJWTIDPResponse, error) {
	config, err := s.command.AddIDPConfig(ctx, AddJWTIDPRequestToDomain(req), authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.AddOrgJWTIDPResponse{
		IdpId: config.IDPConfigID,
		Details: object_pb.AddToDetailsPb(
			config.Sequence,
			config.ChangeDate,
			config.ResourceOwner,
		),
	}, nil
}

func (s *Server) DeactivateOrgIDP(ctx context.Context, req *mgmt_pb.DeactivateOrgIDPRequest) (*mgmt_pb.DeactivateOrgIDPResponse, error) {
	objectDetails, err := s.command.DeactivateIDPConfig(ctx, req.IdpId, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.DeactivateOrgIDPResponse{Details: object_pb.DomainToChangeDetailsPb(objectDetails)}, nil
}

func (s *Server) ReactivateOrgIDP(ctx context.Context, req *mgmt_pb.ReactivateOrgIDPRequest) (*mgmt_pb.ReactivateOrgIDPResponse, error) {
	objectDetails, err := s.command.ReactivateIDPConfig(ctx, req.IdpId, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ReactivateOrgIDPResponse{Details: object_pb.DomainToChangeDetailsPb(objectDetails)}, nil
}

func (s *Server) RemoveOrgIDP(ctx context.Context, req *mgmt_pb.RemoveOrgIDPRequest) (*mgmt_pb.RemoveOrgIDPResponse, error) {
	idp, err := s.query.IDPByIDAndResourceOwner(ctx, true, req.IdpId, authz.GetCtxData(ctx).OrgID, true)
	if err != nil {
		return nil, err
	}
	idpQuery, err := query.NewIDPUserLinkIDPIDSearchQuery(req.IdpId)
	if err != nil {
		return nil, err
	}
	userLinks, err := s.query.IDPUserLinks(ctx, &query.IDPUserLinksSearchQuery{
		Queries: []query.SearchQuery{idpQuery},
	}, true)
	if err != nil {
		return nil, err
	}
	_, err = s.command.RemoveIDPConfig(ctx, req.IdpId, authz.GetCtxData(ctx).OrgID, idp != nil, userLinksToDomain(userLinks.Links)...)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.RemoveOrgIDPResponse{}, nil
}

func (s *Server) UpdateOrgIDP(ctx context.Context, req *mgmt_pb.UpdateOrgIDPRequest) (*mgmt_pb.UpdateOrgIDPResponse, error) {
	config, err := s.command.ChangeIDPConfig(ctx, updateIDPToDomain(req), authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.UpdateOrgIDPResponse{
		Details: object_pb.ChangeToDetailsPb(
			config.Sequence,
			config.ChangeDate,
			config.ResourceOwner,
		),
	}, nil
}

func (s *Server) UpdateOrgIDPOIDCConfig(ctx context.Context, req *mgmt_pb.UpdateOrgIDPOIDCConfigRequest) (*mgmt_pb.UpdateOrgIDPOIDCConfigResponse, error) {
	config, err := s.command.ChangeIDPOIDCConfig(ctx, updateOIDCConfigToDomain(req), authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.UpdateOrgIDPOIDCConfigResponse{
		Details: object_pb.ChangeToDetailsPb(
			config.Sequence,
			config.ChangeDate,
			config.ResourceOwner,
		),
	}, nil
}

func (s *Server) UpdateOrgIDPJWTConfig(ctx context.Context, req *mgmt_pb.UpdateOrgIDPJWTConfigRequest) (*mgmt_pb.UpdateOrgIDPJWTConfigResponse, error) {
	config, err := s.command.ChangeIDPJWTConfig(ctx, updateJWTConfigToDomain(req), authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.UpdateOrgIDPJWTConfigResponse{
		Details: object_pb.ChangeToDetailsPb(
			config.Sequence,
			config.ChangeDate,
			config.ResourceOwner,
		),
	}, nil
}

func (s *Server) AddGenericOAuthProvider(ctx context.Context, req *mgmt_pb.AddGenericOAuthProviderRequest) (*mgmt_pb.AddGenericOAuthProviderResponse, error) {
	id, details, err := s.command.AddOrgGenericOAuthProvider(
		ctx,
		authz.GetCtxData(ctx).OrgID,
		addGenericOAuthProviderToCommand(req),
	)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.AddGenericOAuthProviderResponse{
		Id:      id,
		Details: object_pb.DomainToAddDetailsPb(details),
	}, nil
}

func (s *Server) UpdateGenericOAuthProvider(ctx context.Context, req *mgmt_pb.UpdateGenericOAuthProviderRequest) (*mgmt_pb.UpdateGenericOAuthProviderResponse, error) {
	details, err := s.command.UpdateOrgGenericOAuthProvider(ctx, authz.GetCtxData(ctx).OrgID, req.Id, updateGenericOAuthProviderToCommand(req))
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.UpdateGenericOAuthProviderResponse{
		Details: object_pb.DomainToChangeDetailsPb(details),
	}, nil
}

func (s *Server) AddGenericOIDCProvider(ctx context.Context, req *mgmt_pb.AddGenericOIDCProviderRequest) (*mgmt_pb.AddGenericOIDCProviderResponse, error) {
	id, details, err := s.command.AddOrgGenericOIDCProvider(
		ctx,
		authz.GetCtxData(ctx).OrgID,
		addGenericOIDCProviderToCommand(req),
	)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.AddGenericOIDCProviderResponse{
		Id:      id,
		Details: object_pb.DomainToAddDetailsPb(details),
	}, nil
}

func (s *Server) UpdateGenericOIDCProvider(ctx context.Context, req *mgmt_pb.UpdateGenericOIDCProviderRequest) (*mgmt_pb.UpdateGenericOIDCProviderResponse, error) {
	details, err := s.command.UpdateOrgGenericOIDCProvider(ctx, authz.GetCtxData(ctx).OrgID, req.Id, updateGenericOIDCProviderToCommand(req))
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.UpdateGenericOIDCProviderResponse{
		Details: object_pb.DomainToChangeDetailsPb(details),
	}, nil
}

func (s *Server) AddJWTProvider(ctx context.Context, req *mgmt_pb.AddJWTProviderRequest) (*mgmt_pb.AddJWTProviderResponse, error) {
	id, details, err := s.command.AddOrgJWTProvider(
		ctx,
		authz.GetCtxData(ctx).OrgID,
		addJWTProviderToCommand(req),
	)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.AddJWTProviderResponse{
		Id:      id,
		Details: object_pb.DomainToAddDetailsPb(details),
	}, nil
}

func (s *Server) UpdateJWTProvider(ctx context.Context, req *mgmt_pb.UpdateJWTProviderRequest) (*mgmt_pb.UpdateJWTProviderResponse, error) {
	details, err := s.command.UpdateOrgJWTProvider(ctx, authz.GetCtxData(ctx).OrgID, req.Id, updateJWTProviderToCommand(req))
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.UpdateJWTProviderResponse{
		Details: object_pb.DomainToChangeDetailsPb(details),
	}, nil
}

func (s *Server) AddAzureADProvider(ctx context.Context, req *mgmt_pb.AddAzureADProviderRequest) (*mgmt_pb.AddAzureADProviderResponse, error) {
	id, details, err := s.command.AddOrgAzureADProvider(ctx, authz.GetCtxData(ctx).OrgID, req.ClientId, req.ClientSecret, idp_grpc.OptionsToCommand(req.ProviderOptions))
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.AddAzureADProviderResponse{
		Id:      id,
		Details: object_pb.DomainToAddDetailsPb(details),
	}, nil
}

func (s *Server) UpdateAzureADProvider(ctx context.Context, req *mgmt_pb.UpdateAzureADProviderRequest) (*mgmt_pb.UpdateAzureADProviderResponse, error) {
	details, err := s.command.UpdateOrgAzureADProvider(ctx, authz.GetCtxData(ctx).OrgID, req.Id, req.ClientId, req.ClientSecret, idp_grpc.OptionsToCommand(req.ProviderOptions))
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.UpdateAzureADProviderResponse{
		Details: object_pb.DomainToChangeDetailsPb(details),
	}, nil
}

func (s *Server) AddGitHubProvider(ctx context.Context, req *mgmt_pb.AddGitHubProviderRequest) (*mgmt_pb.AddGitHubProviderResponse, error) {
	id, details, err := s.command.AddOrgGitHubProvider(ctx, authz.GetCtxData(ctx).OrgID, req.ClientId, req.ClientSecret, idp_grpc.OptionsToCommand(req.ProviderOptions))
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.AddGitHubProviderResponse{
		Id:      id,
		Details: object_pb.DomainToAddDetailsPb(details),
	}, nil
}

func (s *Server) UpdateGitHubProvider(ctx context.Context, req *mgmt_pb.UpdateGitHubProviderRequest) (*mgmt_pb.UpdateGitHubProviderResponse, error) {
	details, err := s.command.UpdateOrgGitHubProvider(ctx, authz.GetCtxData(ctx).OrgID, req.Id, req.ClientId, req.ClientSecret, idp_grpc.OptionsToCommand(req.ProviderOptions))
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.UpdateGitHubProviderResponse{
		Details: object_pb.DomainToChangeDetailsPb(details),
	}, nil
}

func (s *Server) AddGitHubEnterpriseProvider(ctx context.Context, req *mgmt_pb.AddGitHubEnterpriseProviderRequest) (*mgmt_pb.AddGitHubEnterpriseProviderResponse, error) {
	id, details, err := s.command.AddOrgGitHubEnterpriseProvider(ctx, authz.GetCtxData(ctx).OrgID, req.ClientId, req.ClientSecret, idp_grpc.OptionsToCommand(req.ProviderOptions))
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.AddGitHubEnterpriseProviderResponse{
		Id:      id,
		Details: object_pb.DomainToAddDetailsPb(details),
	}, nil
}

func (s *Server) UpdateGitHubEnterpriseProvider(ctx context.Context, req *mgmt_pb.UpdateGitHubEnterpriseProviderRequest) (*mgmt_pb.UpdateGitHubEnterpriseProviderResponse, error) {
	details, err := s.command.UpdateOrgGitHubEnterpriseProvider(ctx, authz.GetCtxData(ctx).OrgID, req.Id, req.ClientId, req.ClientSecret, idp_grpc.OptionsToCommand(req.ProviderOptions))
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.UpdateGitHubEnterpriseProviderResponse{
		Details: object_pb.DomainToChangeDetailsPb(details),
	}, nil
}

func (s *Server) AddGoogleProvider(ctx context.Context, req *mgmt_pb.AddGoogleProviderRequest) (*mgmt_pb.AddGoogleProviderResponse, error) {
	id, details, err := s.command.AddOrgGoogleProvider(ctx, authz.GetCtxData(ctx).OrgID, req.ClientId, req.ClientSecret, idp_grpc.OptionsToCommand(req.ProviderOptions))
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.AddGoogleProviderResponse{
		Id:      id,
		Details: object_pb.DomainToAddDetailsPb(details),
	}, nil
}

func (s *Server) UpdateGoogleProvider(ctx context.Context, req *mgmt_pb.UpdateGoogleProviderRequest) (*mgmt_pb.UpdateGoogleProviderResponse, error) {
	details, err := s.command.UpdateOrgGoogleProvider(ctx, authz.GetCtxData(ctx).OrgID, req.Id, req.ClientId, req.ClientSecret, idp_grpc.OptionsToCommand(req.ProviderOptions))
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.UpdateGoogleProviderResponse{
		Details: object_pb.DomainToChangeDetailsPb(details),
	}, nil
}

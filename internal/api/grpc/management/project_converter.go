package management

import (
	"encoding/json"

	"github.com/caos/logging"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/caos/zitadel/internal/eventstore/models"
	proj_model "github.com/caos/zitadel/internal/project/model"
	"github.com/caos/zitadel/pkg/management/grpc"
)

func projectFromModel(project *proj_model.Project) *grpc.Project {
	creationDate, err := ptypes.TimestampProto(project.CreationDate)
	logging.Log("GRPC-iejs3").OnError(err).Debug("unable to parse timestamp")

	changeDate, err := ptypes.TimestampProto(project.ChangeDate)
	logging.Log("GRPC-di7rw").OnError(err).Debug("unable to parse timestamp")

	return &grpc.Project{
		Id:           project.AggregateID,
		State:        projectStateFromModel(project.State),
		CreationDate: creationDate,
		ChangeDate:   changeDate,
		Name:         project.Name,
		Sequence:     project.Sequence,
	}
}

func projectSearchResponseFromModel(response *proj_model.ProjectViewSearchResponse) *grpc.ProjectSearchResponse {
	return &grpc.ProjectSearchResponse{
		Offset:      response.Offset,
		Limit:       response.Limit,
		TotalResult: response.TotalResult,
		Result:      projectViewsFromModel(response.Result),
	}
}

func projectViewsFromModel(projects []*proj_model.ProjectView) []*grpc.ProjectView {
	converted := make([]*grpc.ProjectView, len(projects))
	for i, project := range projects {
		converted[i] = projectViewFromModel(project)
	}
	return converted
}

func projectViewFromModel(project *proj_model.ProjectView) *grpc.ProjectView {
	creationDate, err := ptypes.TimestampProto(project.CreationDate)
	logging.Log("GRPC-dlso3").OnError(err).Debug("unable to parse timestamp")

	changeDate, err := ptypes.TimestampProto(project.ChangeDate)
	logging.Log("GRPC-sope3").OnError(err).Debug("unable to parse timestamp")

	return &grpc.ProjectView{
		ProjectId:     project.ProjectID,
		State:         projectStateFromModel(project.State),
		CreationDate:  creationDate,
		ChangeDate:    changeDate,
		Name:          project.Name,
		Sequence:      project.Sequence,
		ResourceOwner: project.ResourceOwner,
	}
}

func projectRoleSearchResponseFromModel(response *proj_model.ProjectRoleSearchResponse) *grpc.ProjectRoleSearchResponse {
	return &grpc.ProjectRoleSearchResponse{
		Offset:      response.Offset,
		Limit:       response.Limit,
		TotalResult: response.TotalResult,
		Result:      projectRoleViewsFromModel(response.Result),
	}
}

func projectRoleViewsFromModel(roles []*proj_model.ProjectRoleView) []*grpc.ProjectRoleView {
	converted := make([]*grpc.ProjectRoleView, len(roles))
	for i, role := range roles {
		converted[i] = projectRoleViewFromModel(role)
	}
	return converted
}

func projectRoleViewFromModel(role *proj_model.ProjectRoleView) *grpc.ProjectRoleView {
	creationDate, err := ptypes.TimestampProto(role.CreationDate)
	logging.Log("GRPC-dlso3").OnError(err).Debug("unable to parse timestamp")

	return &grpc.ProjectRoleView{
		ProjectId:    role.ProjectID,
		CreationDate: creationDate,
		Key:          role.Key,
		Group:        role.Group,
		DisplayName:  role.DisplayName,
		Sequence:     role.Sequence,
	}
}

func projectStateFromModel(state proj_model.ProjectState) grpc.ProjectState {
	switch state {
	case proj_model.ProjectStateActive:
		return grpc.ProjectState_PROJECTSTATE_ACTIVE
	case proj_model.ProjectStateInactive:
		return grpc.ProjectState_PROJECTSTATE_INACTIVE
	default:
		return grpc.ProjectState_PROJECTSTATE_UNSPECIFIED
	}
}

func projectUpdateToModel(project *grpc.ProjectUpdateRequest) *proj_model.Project {
	return &proj_model.Project{
		ObjectRoot: models.ObjectRoot{
			AggregateID: project.Id,
		},
		Name: project.Name,
	}
}

func projectRoleFromModel(role *proj_model.ProjectRole) *grpc.ProjectRole {
	creationDate, err := ptypes.TimestampProto(role.CreationDate)
	logging.Log("GRPC-due83").OnError(err).Debug("unable to parse timestamp")

	changeDate, err := ptypes.TimestampProto(role.ChangeDate)
	logging.Log("GRPC-id93s").OnError(err).Debug("unable to parse timestamp")

	return &grpc.ProjectRole{
		CreationDate: creationDate,
		ChangeDate:   changeDate,
		Sequence:     role.Sequence,
		Key:          role.Key,
		DisplayName:  role.DisplayName,
		Group:        role.Group,
	}
}

func projectRoleAddBulkToModel(bulk *grpc.ProjectRoleAddBulk) []*proj_model.ProjectRole {
	roles := make([]*proj_model.ProjectRole, len(bulk.ProjectRoles))
	for i, role := range bulk.ProjectRoles {
		roles[i] = &proj_model.ProjectRole{
			ObjectRoot: models.ObjectRoot{
				AggregateID: bulk.Id,
			},
			Key:         role.Key,
			DisplayName: role.DisplayName,
			Group:       role.Group,
		}
	}
	return roles
}

func projectRoleAddToModel(role *grpc.ProjectRoleAdd) *proj_model.ProjectRole {
	return &proj_model.ProjectRole{
		ObjectRoot: models.ObjectRoot{
			AggregateID: role.Id,
		},
		Key:         role.Key,
		DisplayName: role.DisplayName,
		Group:       role.Group,
	}
}

func projectRoleChangeToModel(role *grpc.ProjectRoleChange) *proj_model.ProjectRole {
	return &proj_model.ProjectRole{
		ObjectRoot: models.ObjectRoot{
			AggregateID: role.Id,
		},
		Key:         role.Key,
		DisplayName: role.DisplayName,
		Group:       role.Group,
	}
}

func projectSearchRequestsToModel(project *grpc.ProjectSearchRequest) *proj_model.ProjectViewSearchRequest {
	return &proj_model.ProjectViewSearchRequest{
		Offset:  project.Offset,
		Limit:   project.Limit,
		Queries: projectSearchQueriesToModel(project.Queries),
	}
}
func grantedProjectSearchRequestsToModel(request *grpc.GrantedProjectSearchRequest) *proj_model.ProjectGrantViewSearchRequest {
	return &proj_model.ProjectGrantViewSearchRequest{
		Offset:  request.Offset,
		Limit:   request.Limit,
		Queries: grantedPRojectSearchQueriesToModel(request.Queries),
	}
}

func projectSearchQueriesToModel(queries []*grpc.ProjectSearchQuery) []*proj_model.ProjectViewSearchQuery {
	converted := make([]*proj_model.ProjectViewSearchQuery, len(queries))
	for i, q := range queries {
		converted[i] = projectSearchQueryToModel(q)
	}
	return converted
}

func projectSearchQueryToModel(query *grpc.ProjectSearchQuery) *proj_model.ProjectViewSearchQuery {
	return &proj_model.ProjectViewSearchQuery{
		Key:    projectSearchKeyToModel(query.Key),
		Method: searchMethodToModel(query.Method),
		Value:  query.Value,
	}
}

func projectSearchKeyToModel(key grpc.ProjectSearchKey) proj_model.ProjectViewSearchKey {
	switch key {
	case grpc.ProjectSearchKey_PROJECTSEARCHKEY_PROJECT_NAME:
		return proj_model.ProjectViewSearchKeyName
	default:
		return proj_model.ProjectViewSearchKeyUnspecified
	}
}

func grantedPRojectSearchQueriesToModel(queries []*grpc.ProjectSearchQuery) []*proj_model.ProjectGrantViewSearchQuery {
	converted := make([]*proj_model.ProjectGrantViewSearchQuery, len(queries))
	for i, q := range queries {
		converted[i] = grantedProjectSearchQueryToModel(q)
	}
	return converted
}

func grantedProjectSearchQueryToModel(query *grpc.ProjectSearchQuery) *proj_model.ProjectGrantViewSearchQuery {
	return &proj_model.ProjectGrantViewSearchQuery{
		Key:    projectGrantSearchKeyToModel(query.Key),
		Method: searchMethodToModel(query.Method),
		Value:  query.Value,
	}
}

func projectGrantSearchKeyToModel(key grpc.ProjectSearchKey) proj_model.ProjectGrantViewSearchKey {
	switch key {
	case grpc.ProjectSearchKey_PROJECTSEARCHKEY_PROJECT_NAME:
		return proj_model.GrantedProjectSearchKeyName
	default:
		return proj_model.GrantedProjectSearchKeyUnspecified
	}
}

func projectRoleSearchRequestsToModel(role *grpc.ProjectRoleSearchRequest) *proj_model.ProjectRoleSearchRequest {
	return &proj_model.ProjectRoleSearchRequest{
		Offset:  role.Offset,
		Limit:   role.Limit,
		Queries: projectRoleSearchQueriesToModel(role.Queries),
	}
}

func projectRoleSearchQueriesToModel(queries []*grpc.ProjectRoleSearchQuery) []*proj_model.ProjectRoleSearchQuery {
	converted := make([]*proj_model.ProjectRoleSearchQuery, len(queries))
	for i, q := range queries {
		converted[i] = projectRoleSearchQueryToModel(q)
	}
	return converted
}

func projectRoleSearchQueryToModel(query *grpc.ProjectRoleSearchQuery) *proj_model.ProjectRoleSearchQuery {
	return &proj_model.ProjectRoleSearchQuery{
		Key:    projectRoleSearchKeyToModel(query.Key),
		Method: searchMethodToModel(query.Method),
		Value:  query.Value,
	}
}

func projectRoleSearchKeyToModel(key grpc.ProjectRoleSearchKey) proj_model.ProjectRoleSearchKey {
	switch key {
	case grpc.ProjectRoleSearchKey_PROJECTROLESEARCHKEY_KEY:
		return proj_model.ProjectRoleSearchKeyKey
	case grpc.ProjectRoleSearchKey_PROJECTROLESEARCHKEY_DISPLAY_NAME:
		return proj_model.ProjectRoleSearchKeyDisplayName
	default:
		return proj_model.ProjectRoleSearchKeyUnspecified
	}
}

func projectChangesToResponse(response *proj_model.ProjectChanges, offset uint64, limit uint64) (_ *grpc.Changes) {
	return &grpc.Changes{
		Limit:   limit,
		Offset:  offset,
		Changes: projectChangesToMgtAPI(response),
	}
}

func projectChangesToMgtAPI(changes *proj_model.ProjectChanges) (_ []*grpc.Change) {
	result := make([]*grpc.Change, len(changes.Changes))

	for i, change := range changes.Changes {
		b, err := json.Marshal(change.Data)
		data := &structpb.Struct{}
		err = protojson.Unmarshal(b, data)
		if err != nil {
		}
		result[i] = &grpc.Change{
			ChangeDate: change.ChangeDate,
			EventType:  change.EventType,
			Sequence:   change.Sequence,
			Editor:     change.ModifierName,
			EditorId:   change.ModifierId,
			Data:       data,
		}
	}

	return result
}

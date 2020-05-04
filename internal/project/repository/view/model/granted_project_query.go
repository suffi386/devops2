package model

import (
	global_model "github.com/caos/zitadel/internal/model"
	proj_model "github.com/caos/zitadel/internal/project/model"
	"github.com/caos/zitadel/internal/view"
)

type GrantedProjectSearchRequest proj_model.GrantedProjectSearchRequest
type GrantedProjectSearchQuery proj_model.GrantedProjectSearchQuery
type GrantedProjectSearchKey proj_model.GrantedProjectSearchKey

func (req GrantedProjectSearchRequest) GetLimit() uint64 {
	return req.Limit
}

func (req GrantedProjectSearchRequest) GetOffset() uint64 {
	return req.Offset
}

func (req GrantedProjectSearchRequest) GetSortingColumn() view.ColumnKey {
	if req.SortingColumn == proj_model.GRANTEDPROJECTSEARCHKEY_UNSPECIFIED {
		return nil
	}
	return GrantedProjectSearchKey(req.SortingColumn)
}

func (req GrantedProjectSearchRequest) GetAsc() bool {
	return req.Asc
}

func (req GrantedProjectSearchRequest) GetQueries() []view.SearchQuery {
	result := make([]view.SearchQuery, 0)
	for _, q := range req.Queries {
		result = append(result, GrantedProjectSearchQuery{Key: q.Key, Value: q.Value, Method: q.Method})
	}
	return result
}

func (req GrantedProjectSearchQuery) GetKey() view.ColumnKey {
	return GrantedProjectSearchKey(req.Key)
}

func (req GrantedProjectSearchQuery) GetMethod() global_model.SearchMethod {
	return req.Method
}

func (req GrantedProjectSearchQuery) GetValue() interface{} {
	return req.Value
}

func (key GrantedProjectSearchKey) ToColumnName() string {
	switch proj_model.GrantedProjectSearchKey(key) {
	case proj_model.GRANTEDPROJECTSEARCHKEY_NAME:
		return GrantedProjectNameKey
	case proj_model.GRANTEDPROJECTSEARCHKEY_GRANTID:
		return GrantedProjectGrantIDKey
	case proj_model.GRANTEDPROJECTSEARCHKEY_ORGID:
		return GrantedProjectOrgIDKey
	case proj_model.GRANTEDPROJECTSEARCHKEY_PROJECTID:
		return GrantedProjectIDKey
	default:
		return ""
	}
}

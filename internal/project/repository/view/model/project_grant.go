package model

import (
	"encoding/json"
	"time"

	"github.com/caos/logging"
	"github.com/lib/pq"

	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/project/model"
	"github.com/caos/zitadel/internal/repository/project"
)

const (
	ProjectGrantKeyProjectID     = "project_id"
	ProjectGrantKeyGrantID       = "grant_id"
	ProjectGrantKeyOrgID         = "org_id"
	ProjectGrantKeyResourceOwner = "resource_owner"
	ProjectGrantKeyName          = "project_name"
	ProjectGrantKeyRoleKeys      = "granted_role_keys"
)

type ProjectGrantView struct {
	GrantID           string         `json:"-" gorm:"column:grant_id;primary_key"`
	ProjectID         string         `json:"-" gorm:"column:project_id"`
	OrgID             string         `json:"-" gorm:"column:org_id"`
	Name              string         `json:"name" gorm:"column:project_name"`
	CreationDate      time.Time      `json:"-" gorm:"column:creation_date"`
	ChangeDate        time.Time      `json:"-" gorm:"column:change_date"`
	State             int32          `json:"-" gorm:"column:project_state"`
	ResourceOwner     string         `json:"-" gorm:"column:resource_owner"`
	ResourceOwnerName string         `json:"-" gorm:"column:resource_owner_name"`
	OrgName           string         `json:"-" gorm:"column:org_name"`
	Sequence          uint64         `json:"-" gorm:"column:sequence"`
	GrantedRoleKeys   pq.StringArray `json:"-" gorm:"column:granted_role_keys"`
}

type ProjectGrant struct {
	GrantID      string   `json:"grantId"`
	GrantedOrgID string   `json:"grantedOrgId"`
	RoleKeys     []string `json:"roleKeys"`
	InstanceID   string   `json:"instanceID"`
}

func ProjectGrantFromModel(project *model.ProjectGrantView) *ProjectGrantView {
	return &ProjectGrantView{
		ProjectID:         project.ProjectID,
		OrgID:             project.OrgID,
		Name:              project.Name,
		ChangeDate:        project.ChangeDate,
		CreationDate:      project.CreationDate,
		State:             int32(project.State),
		ResourceOwner:     project.ResourceOwner,
		ResourceOwnerName: project.ResourceOwnerName,
		OrgName:           project.OrgName,
		GrantID:           project.GrantID,
		GrantedRoleKeys:   project.GrantedRoleKeys,
		Sequence:          project.Sequence,
	}
}

func ProjectGrantToModel(project *ProjectGrantView) *model.ProjectGrantView {
	return &model.ProjectGrantView{
		ProjectID:         project.ProjectID,
		OrgID:             project.OrgID,
		Name:              project.Name,
		ChangeDate:        project.ChangeDate,
		CreationDate:      project.CreationDate,
		State:             model.ProjectState(project.State),
		ResourceOwner:     project.ResourceOwner,
		ResourceOwnerName: project.ResourceOwnerName,
		OrgName:           project.OrgName,
		GrantID:           project.GrantID,
		Sequence:          project.Sequence,
		GrantedRoleKeys:   project.GrantedRoleKeys,
	}
}

func ProjectGrantsToModel(projects []*ProjectGrantView) []*model.ProjectGrantView {
	result := make([]*model.ProjectGrantView, len(projects))
	for i, p := range projects {
		result[i] = ProjectGrantToModel(p)
	}
	return result
}

func (p *ProjectGrantView) AppendEvent(event *models.Event) (err error) {
	p.ChangeDate = event.CreationDate
	p.Sequence = event.Sequence
	switch eventstore.EventType(event.Type) {
	case project.GrantAddedType:
		p.State = int32(model.ProjectStateActive)
		p.CreationDate = event.CreationDate
		p.setRootData(event)
		err = p.setProjectGrantData(event)
	case project.GrantChangedType, project.GrantCascadeChangedType:
		err = p.setProjectGrantData(event)
	case project.GrantDeactivatedType:
		p.State = int32(model.ProjectStateInactive)
	case project.GrantReactivatedType:
		p.State = int32(model.ProjectStateActive)
	}
	return err
}

func (p *ProjectGrantView) setRootData(event *models.Event) {
	p.ProjectID = event.AggregateID
	p.ResourceOwner = event.ResourceOwner
}

func (p *ProjectGrantView) setData(event *models.Event) error {
	if err := json.Unmarshal(event.Data, p); err != nil {
		logging.Log("EVEN-dlo92").WithError(err).Error("could not unmarshal event data")
		return err
	}
	return nil
}

func (p *ProjectGrantView) setProjectGrantData(event *models.Event) error {
	grant := new(ProjectGrant)
	err := grant.SetData(event)
	if err != nil {
		return err
	}
	if grant.GrantedOrgID != "" {
		p.OrgID = grant.GrantedOrgID
	}
	p.GrantID = grant.GrantID
	p.GrantedRoleKeys = grant.RoleKeys
	return nil
}

func (p *ProjectGrant) SetData(event *models.Event) error {
	if err := json.Unmarshal(event.Data, p); err != nil {
		logging.Log("EVEN-dlo92").WithError(err).Error("could not unmarshal event data")
		return caos_errs.ThrowInternal(err, "MODEL-s9ols", "Could not unmarshal data")
	}
	return nil
}

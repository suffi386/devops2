package command

import (
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/restrictions"
)

type restrictionsWriteModel struct {
	eventstore.WriteModel
	publicOrgRegistrationIsNotAllowed bool
	allowedLanguages                  []language.Tag
}

// newRestrictionsWriteModel aggregateId is filled by reducing unit matching events
func newRestrictionsWriteModel(instanceId, resourceOwner string) *restrictionsWriteModel {
	return &restrictionsWriteModel{
		WriteModel: eventstore.WriteModel{
			InstanceID:    instanceId,
			ResourceOwner: resourceOwner,
		},
	}
}

func (wm *restrictionsWriteModel) Query() *eventstore.SearchQueryBuilder {
	query := eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		InstanceID(wm.InstanceID).
		AddQuery().
		AggregateTypes(restrictions.AggregateType).
		EventTypes(restrictions.SetEventType)

	return query.Builder()
}

func (wm *restrictionsWriteModel) Reduce() error {
	for _, event := range wm.Events {
		wm.ChangeDate = event.CreatedAt()
		e, ok := event.(*restrictions.SetEvent)
		if !ok {
			continue
		}
		if e.PublicOrgRegistrationIsNotAllowed != nil {
			wm.publicOrgRegistrationIsNotAllowed = *e.PublicOrgRegistrationIsNotAllowed
		}
		if e.AllowedLanguages != nil {
			wm.allowedLanguages = *e.AllowedLanguages
		}
	}
	return wm.WriteModel.Reduce()
}

// NewChanges returns all changes that need to be applied to the aggregate.
// nil properties in setRestrictions are ignored
func (wm *restrictionsWriteModel) NewChanges(setRestrictions *SetRestrictions) (changes []restrictions.RestrictionsChange, languagesChanged bool) {
	if setRestrictions == nil {
		return nil, false
	}
	changes = make([]restrictions.RestrictionsChange, 0, 1)
	if setRestrictions.DisallowPublicOrgRegistration != nil && (wm.publicOrgRegistrationIsNotAllowed != *setRestrictions.DisallowPublicOrgRegistration) {
		changes = append(changes, restrictions.ChangePublicOrgRegistrations(*setRestrictions.DisallowPublicOrgRegistration))
	}
	if setRestrictions.AllowedLanguages != nil && domain.LanguagesDiffer(wm.allowedLanguages, setRestrictions.AllowedLanguages) {
		changes = append(changes, restrictions.ChangeAllowedLanguages(setRestrictions.AllowedLanguages))
		languagesChanged = true
	}
	return changes, languagesChanged
}

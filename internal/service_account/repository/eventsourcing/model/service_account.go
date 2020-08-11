package model

import (
	"encoding/json"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/service_account/model"
)

type ServiceAccount struct {
	models.ObjectRoot

	Name        string
	Email       string
	Description string
	State       int32
}

func (sa *ServiceAccount) AppendEvents(events ...*models.Event) error {
	for _, event := range events {
		if err := sa.AppendEvent(event); err != nil {
			return err
		}
	}
	return nil
}

func (sa *ServiceAccount) AppendEvent(event *models.Event) (err error) {
	sa.ObjectRoot.AppendEvent(event)

	switch event.Type {
	case ServiceAccountAdded:
		sa.State = model.ServiceAccountStateActive
		fallthrough
	case ServiceAccountChanged:
		err = sa.setData(event)
	case ServiceAccountLocked:
		sa.State = model.ServiceAccountStateLocked
	case ServiceAccountDeactivated:
		sa.State = model.ServiceAccountStateInactive
	case ServiceAccountRemoved:
		sa.State = model.ServiceAccountStateDeleted
	case ServiceAccountUnlocked, ServiceAccountReactivated:
		sa.State = model.ServiceAccountStateActive

	case KeyAdded:
		fallthrough
	case KeyRemoved:
		logging.Log("MODEL-iBgOc").Panic("key unimplemented")
	}

	return err
}

func (sa *ServiceAccount) setData(event *models.Event) error {
	if err := json.Unmarshal(event.Data, sa); err != nil {
		logging.Log("EVEN-8ujgd").WithError(err).Error("could not unmarshal event data")
		return errors.ThrowInternal(err, "MODEL-GwjY9", "could not unmarshal event")
	}
	return nil
}

func ServiceAccountFromModel(serviceAccount *model.ServiceAccount) *ServiceAccount {
	return &ServiceAccount{
		ObjectRoot:  serviceAccount.ObjectRoot,
		Description: serviceAccount.Description,
		Email:       serviceAccount.Email,
		Name:        serviceAccount.Name,
	}
}

func ServiceAccountToModel(serviceAccount *ServiceAccount) *model.ServiceAccount {
	return &model.ServiceAccount{
		ObjectRoot:  serviceAccount.ObjectRoot,
		Description: serviceAccount.Description,
		Email:       serviceAccount.Email,
		Name:        serviceAccount.Name,
	}
}

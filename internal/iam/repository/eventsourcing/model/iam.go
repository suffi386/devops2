package model

import (
	"encoding/json"
	"github.com/caos/logging"
	caos_errs "github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/iam/model"
)

const (
	IamVersion = "v1"
)

type Iam struct {
	es_models.ObjectRoot
	SetUpStarted bool   `json:"-"`
	SetUpDone    bool   `json:"-"`
	GlobalOrgID  string `json:"globalOrgId,omitempty"`
	IamProjectID string `json:"iamProjectId,omitempty"`
}

func IamFromModel(iam *model.Iam) *Iam {
	converted := &Iam{
		ObjectRoot:   iam.ObjectRoot,
		SetUpStarted: iam.SetUpStarted,
		SetUpDone:    iam.SetUpDone,
		GlobalOrgID:  iam.GlobalOrgID,
		IamProjectID: iam.IamProjectID,
	}
	return converted
}

func IamToModel(iam *Iam) *model.Iam {
	converted := &model.Iam{
		ObjectRoot:   iam.ObjectRoot,
		SetUpStarted: iam.SetUpStarted,
		SetUpDone:    iam.SetUpDone,
		GlobalOrgID:  iam.GlobalOrgID,
		IamProjectID: iam.IamProjectID,
	}
	return converted
}

func (i *Iam) AppendEvents(events ...*es_models.Event) error {
	for _, event := range events {
		if err := i.AppendEvent(event); err != nil {
			return err
		}
	}
	return nil
}

func (i *Iam) AppendEvent(event *es_models.Event) (err error) {
	i.ObjectRoot.AppendEvent(event)
	switch event.Type {
	case IamSetupStarted:
		i.SetUpStarted = true
	case IamSetupDone:
		i.SetUpDone = true
	case IamProjectSet,
		GlobalOrgSet:
		i.setData(event)
	}
	if err != nil {
		return err
	}
	return nil
}

func (i *Iam) setData(event *es_models.Event) error {
	i.ObjectRoot.AppendEvent(event)
	if err := json.Unmarshal(event.Data, i); err != nil {
		logging.Log("EVEN-9sie4").WithError(err).Error("could not unmarshal event data")
		return caos_errs.ThrowInternal(err, "MODEL-slwi3", "could not unmarshal event")
	}
	return nil
}

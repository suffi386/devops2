package model

import (
	"encoding/json"
	"github.com/caos/logging"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/user/model"
	es_model "github.com/caos/zitadel/internal/user/repository/eventsourcing/model"
	"time"
)

const (
	NotifyUserKeyUserID = "id"
)

type NotifyUser struct {
	ID                string    `json:"-" gorm:"column:id;primary_key"`
	CreationDate      time.Time `json:"-" gorm:"column:creation_date"`
	ChangeDate        time.Time `json:"-" gorm:"column:change_date"`
	ResourceOwner     string    `json:"-" gorm:"column:resource_owner"`
	State             int32     `json:"-" gorm:"column:user_state"`
	UserName          string    `json:"userName" gorm:"column:user_name"`
	FirstName         string    `json:"firstName" gorm:"column:first_name"`
	LastName          string    `json:"lastName" gorm:"column:last_name"`
	NickName          string    `json:"nickName" gorm:"column:nick_name"`
	DisplayName       string    `json:"displayName" gorm:"column:display_name"`
	PreferredLanguage string    `json:"preferredLanguage" gorm:"column:preferred_language"`
	Gender            int32     `json:"gender" gorm:"column:gender"`
	LastEmail         string    `json:"email" gorm:"column:last_email"`
	VerifiedEmail     string    `json:"-" gorm:"column:verified_email"`
	LastPhone         string    `json:"phone" gorm:"column:last_phone"`
	VerifiedPhone     string    `json:"-" gorm:"column:verified_phone"`
	Sequence          uint64    `json:"-" gorm:"column:sequence"`
}

func NotifyUserFromModel(user *model.NotifyUser) *NotifyUser {
	return &NotifyUser{
		ID:                user.ID,
		ChangeDate:        user.ChangeDate,
		CreationDate:      user.CreationDate,
		ResourceOwner:     user.ResourceOwner,
		State:             int32(user.State),
		UserName:          user.UserName,
		FirstName:         user.FirstName,
		LastName:          user.LastName,
		NickName:          user.NickName,
		DisplayName:       user.DisplayName,
		PreferredLanguage: user.PreferredLanguage,
		Gender:            int32(user.Gender),
		LastEmail:         user.LastEmail,
		VerifiedEmail:     user.VerifiedEmail,
		LastPhone:         user.LastPhone,
		VerifiedPhone:     user.VerifiedPhone,
		Sequence:          user.Sequence,
	}
}

func NotifyUserToModel(user *NotifyUser) *model.NotifyUser {
	return &model.NotifyUser{
		ID:                user.ID,
		ChangeDate:        user.ChangeDate,
		CreationDate:      user.CreationDate,
		ResourceOwner:     user.ResourceOwner,
		State:             model.UserState(user.State),
		UserName:          user.UserName,
		FirstName:         user.FirstName,
		LastName:          user.LastName,
		NickName:          user.NickName,
		DisplayName:       user.DisplayName,
		PreferredLanguage: user.PreferredLanguage,
		Gender:            model.Gender(user.Gender),
		LastEmail:         user.LastEmail,
		VerifiedEmail:     user.VerifiedEmail,
		LastPhone:         user.LastPhone,
		VerifiedPhone:     user.VerifiedPhone,
		Sequence:          user.Sequence,
	}
}

func (p *NotifyUser) AppendEvent(event *models.Event) (err error) {
	p.ChangeDate = event.CreationDate
	p.Sequence = event.Sequence
	switch event.Type {
	case es_model.UserAdded,
		es_model.UserRegistered:
		p.CreationDate = event.CreationDate
		p.setRootData(event)
		err = p.setData(event)
	case es_model.UserProfileChanged:
		err = p.setData(event)
	case es_model.UserEmailChanged:
		err = p.setData(event)
	case es_model.UserEmailVerified:
		p.VerifiedEmail = p.LastEmail
	case es_model.UserPhoneChanged:
		err = p.setData(event)
	case es_model.UserPhoneVerified:
		p.VerifiedPhone = p.LastPhone
	}
	return err
}

func (u *NotifyUser) setRootData(event *models.Event) {
	u.ID = event.AggregateID
	u.ResourceOwner = event.ResourceOwner
}

func (u *NotifyUser) setData(event *models.Event) error {
	if err := json.Unmarshal(event.Data, u); err != nil {
		logging.Log("EVEN-lso9e").WithError(err).Error("could not unmarshal event data")
		return caos_errs.ThrowInternal(nil, "MODEL-8iows", "could not unmarshal data")
	}
	return nil
}

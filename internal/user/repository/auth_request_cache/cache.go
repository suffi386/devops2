package cache

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/caos/zitadel/internal/config/types"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/user/model"
)

type Config struct {
	Connection types.SQL
}

type AuthRequestCache struct {
	client *sql.DB
}

func Start(conf Config) (*AuthRequestCache, error) {
	client, err := sql.Open("postgres", conf.Connection.ConnectionString())
	if err != nil {
		return nil, caos_errs.ThrowPreconditionFailed(err, "SQL-9qBtr", "unable to open database connection")
	}
	return &AuthRequestCache{
		client: client,
	}, nil
}

func (c *AuthRequestCache) Health(ctx context.Context) error {
	return c.client.PingContext(ctx)
}

func (c *AuthRequestCache) GetAuthRequestByID(_ context.Context, id string) (*model.AuthRequest, error) {
	var b []byte
	err := c.client.QueryRow("SELECT request FROM auth.authrequests WHERE request_id = ?", id).Scan(&b)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, caos_errs.ThrowNotFound(err, "CACHE-d24aD", "auth request not found")
		}
		return nil, caos_errs.ThrowInternal(err, "CACHE-as3kj", "unable to get auth request from database")
	}
	request := new(model.AuthRequest)
	err = json.Unmarshal(b, &request)
	if err != nil {
		return nil, caos_errs.ThrowInternal(err, "CACHE-2wshg", "unable to unmarshal auth request")
	}
	return request, nil
}

func (c *AuthRequestCache) SaveAuthRequest(_ context.Context, request *model.AuthRequest) error {
	b, err := json.Marshal(request)
	if err != nil {
		return caos_errs.ThrowInternal(err, "CACHE-32FH9", "unable to marshal auth request")
	}
	stmt, err := c.client.Prepare("INSERT INTO auth.authrequests (request_id, request) VALUES($1, $2)")
	if err != nil {
		return caos_errs.ThrowInternal(err, "CACHE-dswfF", "sql prepare failed")
	}
	_, err = stmt.Exec(request.AggregateID, b)
	if err != nil {
		return caos_errs.ThrowInternal(err, "CACHE-sw4af", "unable to save auth request")
	}
	return nil
}

//
//func (c *AuthRequestCache) GetAuthRequestByID(ctx context.Context, id string) (*model.AuthRequest, error) {
//	request := new(model.AuthRequest)
//	row := c.client.QueryRow(query, id)
//	err := row.Scan(
//		&request.AggregateID,
//		&request.ObjectRoot.ChangeDate,
//		&request.ObjectRoot.CreationDate,
//		&request.ObjectRoot.Sequence,
//		&request.BrowserInfo.UserAgent,
//		&request.BrowserInfo.RemoteIP,
//		&request.BrowserInfo.AcceptLanguage,
//		&request.ApplicationID,
//		&request.CallbackURI,
//		&request.TransferState,
//		&request.Prompt,
//		&request.PossibleLOAs,
//		&request.UiLocales,
//		&request.LoginHint,
//		&request.PreselectedUserID,
//		&request.MaxAuthAge,
//		&request.Request,
//	)
//	if err == nil {
//		return request, nil
//	}
//	if errors.Is(err, sql.ErrNoRows) {
//		return nil, caos_errs.ThrowNotFound(err, "CACHE-d24aD", "auth request not found")
//	}
//	return nil, caos_errs.ThrowInternal(err, "CACHE-as3kj", "unable to get auth request from database")
//}
//
//func (c *AuthRequestCache) SaveAuthRequest(ctx context.Context, request *model.AuthRequest) (*model.AuthRequest, error) {
//
//	row := c.client.QueryRow(query, id)
//}

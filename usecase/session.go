package usecase

import (
	"context"
	db "gateway_mrc/db/sqlc"
	"gateway_mrc/helpers/server"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"time"
)

type SessionUsecase struct {
	store db.Store
}

func NewSessionUsecase(store db.Store) SessionUsecase {
	return SessionUsecase{store}
}

func (uc *SessionUsecase) UpdateSession(ctx context.Context, sessionIdStr string, refreshToken, accessToken string) (db.Session, int32, error) {
	var session db.Session
	var err error
	sessionId, err := uuid.Parse(sessionIdStr)
	if err != nil {
		return session, server.SESSION_PARSING_ERR_CODE, err
	}
	sessionArgs := &db.UpdateSessionDataParams{
		AccessToken:  pgtype.Text{String: accessToken, Valid: accessToken != ""},
		RefreshToken: pgtype.Text{String: refreshToken, Valid: refreshToken != ""},
		ID:           sessionId,
	}
	session, err = uc.store.UpdateSessionData(ctx, *sessionArgs)
	if err != nil {
		return session, db.ErrorCode(err), err
	}
	return session, server.SUCCESS_CODE, nil
}

func (uc *SessionUsecase) UpdateSessionLastActive(ctx context.Context, sessionId uuid.UUID) (db.Session, int32, error) {
	sessionArgs := &db.UpdateSessionDataParams{
		LastActive: pgtype.Timestamptz{Time: time.Now(), Valid: true},
		ID:         sessionId,
	}
	session, err := uc.store.UpdateSessionData(ctx, *sessionArgs)
	if err != nil {
		return session, db.ErrorCode(err), err
	}
	return session, server.SUCCESS_CODE, nil
}

func (uc *SessionUsecase) GetSession(ctx context.Context, id string) (db.Session, int32, error) {
	var session db.Session
	sessionId, err := uuid.Parse(id)
	if err != nil {
		return session, db.ErrorCode(err), err
	}
	session, err = uc.store.GetSession(ctx, sessionId)
	if err != nil {
		return session, db.ErrorCode(err), err
	}
	return session, db.ErrorCode(err), nil
}

package usecase

import (
	"context"
	db "gateway_mrc/db/sqlc"
	"gateway_mrc/helpers/server"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"time"
)

type AuthUsecase struct {
	store db.Store
}

func NewAuthUsecase(store db.Store) *AuthUsecase {
	return &AuthUsecase{store}
}

func (uc *AuthUsecase) UpdateSession(ctx context.Context, sessionIdStr string, refreshToken, accessToken string) (db.Session, int32, error) {
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

func (uc *AuthUsecase) UpdateSessionLastActive(ctx context.Context, sessionId uuid.UUID) (db.Session, int32, error) {
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

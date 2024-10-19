// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

import (
	"context"

	"github.com/google/uuid"
)

type Querier interface {
	AVGSessionLength(ctx context.Context) (float64, error)
	BlockSession(ctx context.Context, arg BlockSessionParams) (Session, error)
	CreateSession(ctx context.Context, arg CreateSessionParams) (Session, error)
	DeleteSession(ctx context.Context, id uuid.UUID) error
	GetSession(ctx context.Context, id uuid.UUID) (Session, error)
	UpdateSessionData(ctx context.Context, arg UpdateSessionDataParams) (Session, error)
}

var _ Querier = (*Queries)(nil)

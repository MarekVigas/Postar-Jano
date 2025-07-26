package repository

import (
	"context"
	"github.com/MarekVigas/Postar-Jano/internal/model"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

func FindOwner(ctx context.Context, db sqlx.QueryerContext, username string) (*model.Owner, error) {
	var owner model.Owner
	if err := sqlx.GetContext(ctx, db, &owner, `SELECT * FROM owners WHERE username = $1`, username); err != nil {
		return nil, errors.Wrap(err, "failed to find user")
	}
	return &owner, nil
}

package auth

import (
	"context"

	"github.com/MarekVigas/Postar-Jano/internal/model"
	"github.com/MarekVigas/Postar-Jano/internal/repository"

	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type FromDB struct {
	repo *repository.PostgresRepo
}

func NewFromDB(repo *repository.PostgresRepo) *FromDB {
	return &FromDB{
		repo: repo,
	}
}

func (a *FromDB) Authenticate(ctx context.Context, username string, password string) (*model.Owner, error) {
	owner, err := a.repo.FindOwner(ctx, username)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find user")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(owner.Pass), []byte(password)); err != nil {
		return nil, errors.WithStack(err)
	}
	return owner, nil

}

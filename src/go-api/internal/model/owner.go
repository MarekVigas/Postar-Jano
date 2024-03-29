package model

import (
	"context"

	"golang.org/x/crypto/bcrypt"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type Owner struct {
	ID       int    `db:"id"`
	Name     string `db:"name"`
	Surname  string `db:"surname"`
	Email    string `db:"email"`
	Username string `db:"username"`
	Pass     string `db:"pass"`
	Phone    string `db:"phone"`
	Photo    string `db:"photo"`
	Gender   string `db:"gender"`
}

func (o Owner) Create(ctx context.Context, db sqlx.ExtContext, plainTextPass string) (*Owner, error) {
	pass, err := bcrypt.GenerateFromPassword([]byte(plainTextPass), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var owner Owner
	if err := sqlx.GetContext(ctx, db, &owner, `
		INSERT INTO owners (
			name,
			surname,
			gender,
			username,
			pass,
			email,
			phone,
			photo
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING *
	`, o.Name, o.Surname, o.Gender, o.Username, pass, o.Email,
		o.Phone, o.Photo); err != nil {
		return nil, errors.WithStack(err)
	}
	return &owner, nil
}

package auth

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt"
	"github.com/labstack/echo/v4"
	"time"

	"github.com/MarekVigas/Postar-Jano/internal/model"
	"github.com/MarekVigas/Postar-Jano/internal/repository"

	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

const tokenLifetime = 3 * time.Hour

type FromDB struct {
	postgresDB *repository.PostgresDB
	jwtSecret  []byte
}

func NewFromDB(db *repository.PostgresDB, jwtSecret []byte) *FromDB {
	return &FromDB{
		postgresDB: db,
		jwtSecret:  jwtSecret,
	}
}

func (a *FromDB) generateToken(owner *model.Owner) (string, error) {
	now := time.Now()
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(tokenLifetime)),
			//Id:        uuid.,
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    "sbb.sk",
			NotBefore: jwt.NewNumericDate(now),
			Subject:   owner.Email,
		},
	})

	return tok.SignedString(a.jwtSecret)
}

func (a *FromDB) Authenticate(ctx context.Context, username string, password string) (string, error) {
	owner, err := repository.FindOwner(ctx, a.postgresDB.QueryerContext(), username)
	if err != nil {
		return "", errors.Wrap(err, "failed to find user")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(owner.Pass), []byte(password)); err != nil {
		return "", errors.WithStack(err)
	}
	return a.generateToken(owner)

}

func (a *FromDB) Middleware() echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		SigningKey: a.jwtSecret,
		ErrorHandler: func(c echo.Context, err error) error {
			return echo.ErrUnauthorized
		},
	})
}

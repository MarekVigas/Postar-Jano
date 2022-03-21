package db

import (
	"bytes"
	"database/sql"
	"fmt"

	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

type Config struct {
	User     string `envconfig:"POSTGRES_USER" required:"true"`
	Password string `envconfig:"POSTGRES_PASSWORD" required:"true"`
	Host     string `envconfig:"POSTGRES_HOST" required:"true"`
	Port     uint   `envconfig:"POSTGRES_PORT" required:"true"`
	Database string `envconfig:"POSTGRES_DB" required:"true"`
}

func Connect() (*sql.DB, error) {
	var c Config
	if err := envconfig.Process("", &c); err != nil {
		return nil, errors.Wrap(err, "failed to load config")
	}
	return c.Connect()
}

func (c *Config) Connect() (*sql.DB, error) {
	conn, err := c.ConnectionString()
	if err != nil {
		return nil, err
	}
	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open database")
	}
	return db, nil
}

func (c *Config) ConnectionString() (string, error) {
	var b bytes.Buffer

	add := func(key string, val interface{}) {
		fmt.Fprintf(&b, "%s=%v ", key, val)
	}

	add("user", c.User)
	add("password", c.Password)
	add("host", c.Host)
	add("port", c.Port)
	add("dbname", c.Database)

	fmt.Fprint(&b, "sslmode=disable")
	return b.String(), nil
}

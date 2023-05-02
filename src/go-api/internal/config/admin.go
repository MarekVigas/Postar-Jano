package config

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

type AdminSettings struct {
	Mailer
	Server
	DB
}

func LoadAdminSetting() (*AdminSettings, error) {
	var c AdminSettings
	if err := envconfig.Process("", &c); err != nil {
		return nil, errors.Wrap(err, "failed to load admin settings")
	}
	return &c, nil
}

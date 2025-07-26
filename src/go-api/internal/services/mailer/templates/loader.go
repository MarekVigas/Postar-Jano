package templates

import (
	"html/template"
	"os"

	"github.com/pkg/errors"
)

func LoadFromFile(path string) (*template.Template, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read template file ")
	}
	tpl, err := template.New("").Parse(string(data))
	if err != nil {
		return nil, errors.Wrap(err, "failed to load template")
	}
	return tpl, nil
}

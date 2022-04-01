package templates

import (
	"html/template"
	"os"
)

type ConfirmationReq struct {
	Mail          string
	ParentName    string
	ParentSurname string
	EventName     string
	Name          string
	Surname       string
	Pills         string
	Restrictions  string
	Text          string
	PhotoURL      string
	Sum           int
	Owner         string
	Info          string
	Days          []string
	RegInfo       string
}

const (
	emailTemplatePath = "EMAIL_TEMPLATE_PATH"
)

var Confirmation *template.Template

// TODO: the errors should be handled properly and the configuration processed in the main.
func init() {
	path := os.Getenv(emailTemplatePath)
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	Confirmation = template.Must(template.New("").Parse(string(data)))
}

package templates

import (
    "text/template"
    "bufio"
    "fmt"
    "io"
    "os"
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}

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
}

const (
	emailTemplatePath                = "EMAIL_TEMPLATE_PATH"
)

path := os.Getenv(emailTemplatePath)
dat, err := os.ReadFile(path)
check(err)
const confirmationMail = string(dat)
var Confirmation = template.Must(template.New("").Parse(confirmationMail))

package templates

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

package templates

type PaymentReminderReq struct {
	Mail          string
	ParentName    string
	ParentSurname string
	EventName     string
	Name          string
	Surname       string
	Sum           int
	Payment       PaymentDetails
}

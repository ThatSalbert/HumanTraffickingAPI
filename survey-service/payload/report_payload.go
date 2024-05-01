package payload

type ReportPayload struct {
	ReportId          string `json:"report_id"`
	Date              string `json:"date"`
	Time              string `json:"time"`
	Anonymous         bool   `json:"anonymous"`
	Email             string `json:"email"`
	ReportDescription string `json:"report_description"`
	Country           string `json:"country"`
	City              string `json:"city"`
	Street            string `json:"street"`
}

package payload

type StatsPayload struct {
	TotalSurveyAnswers int `json:"total_survey_answers"`
	TotalReports       int `json:"total_reports"`
}

package payload

type SurveyAnswerPayload struct {
	SurveyAnswerId  string `json:"survey_answer_id"`
	SurveyId        int    `json:"survey_id"`
	Date            string `json:"date"`
	Time            string `json:"time"`
	Anonymous       bool   `json:"anonymous"`
	Email           string `json:"email"`
	QuestionAnswers []struct {
		QuestionId int    `json:"question_id"`
		Answer     string `json:"answer"`
	} `json:"question_answers"`
}

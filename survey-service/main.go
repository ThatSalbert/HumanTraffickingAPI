package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"os"
	"strconv"
	"survey-service/database"
	"survey-service/payload"
)

var db *mongo.Client
var err error
var port = ":8080"

var dbIp = os.Getenv("DB_IP")
var dbPort = os.Getenv("DB_PORT")
var dbUser = os.Getenv("DB_USER")
var dbPassword = os.Getenv("DB_PASSWORD")

func main() {
	db, err = database.ConnectToDatabase(dbIp, dbPort, dbUser, dbPassword)
	if err != nil {
		log.Fatal(err)
	}
	defer func(db *mongo.Client) {
		err := db.Disconnect(nil)
		if err != nil {
			log.Fatal(err)
		}
	}(db)

	router := mux.NewRouter()

	router.HandleFunc("/api/v1/health", HealthCheck).Methods("GET")

	router.HandleFunc("/api/v1/submit-answer", SubmitAnswer).Methods("POST")
	router.HandleFunc("/api/v1/submit-report", SubmitReport).Methods("POST")

	router.HandleFunc("/api/v1/get-survey-answer", GetAllSurveyAnswers).Methods("GET")
	router.HandleFunc("/api/v1/get-survey-answer/{survey_answer_id}", GetSurveyAnswerByID).Methods("GET")

	router.HandleFunc("/api/v1/get-report", GetAllReports).Methods("GET")
	router.HandleFunc("/api/v1/get-report/{report_id}", GetReportByID).Methods("GET")

	router.HandleFunc("/api/v1/get-stats", GetStats).Methods("GET")

	corsOrigins := handlers.AllowedOrigins([]string{"*"})
	corsMethods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	corsHeaders := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})

	log.Fatal(http.ListenAndServe(port, handlers.CORS(corsOrigins, corsMethods, corsHeaders)(router)))
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	log.Println("Health check endpoint accessed")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(`{"status": "ok"}`))
	if err != nil {
		return
	}
	return
}

func SubmitAnswer(w http.ResponseWriter, r *http.Request) {
	log.Println("Submit answer endpoint accessed")
	w.Header().Set("Content-Type", "application/json")
	jsonDecoder := json.NewDecoder(r.Body)
	var answerPayload payload.SurveyAnswerPayload
	err := jsonDecoder.Decode(&answerPayload)
	if err != nil {
		log.Println("Invalid JSON provided")
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(`{"error": "invalid json"}`))
		if err != nil {
			return
		}
		return
	}
	answerPayload.SurveyAnswerId = uuid.New().String()
	fmt.Println(answerPayload)
	err = database.SubmitAnswerToDatabase(db, answerPayload)
	if err != nil {
		log.Println("Internal server error")
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(`{"error": "` + err.Error() + `"}`))
		if err != nil {
			return
		}
		return
	} else {
		w.WriteHeader(http.StatusCreated)
		_, err = w.Write([]byte(`{"status": "created with answer ID: ` + answerPayload.SurveyAnswerId + `"}`))
		if err != nil {
			return
		}
		return
	}
}

func SubmitReport(w http.ResponseWriter, r *http.Request) {
	log.Println("Submit report endpoint accessed")
	w.Header().Set("Content-Type", "application/json")
	jsonDecoder := json.NewDecoder(r.Body)
	var reportPayload payload.ReportPayload
	err := jsonDecoder.Decode(&reportPayload)
	if err != nil {
		log.Println("Invalid JSON provided")
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(`{"error": "invalid json"}`))
		if err != nil {
			return
		}
		return
	}
	reportPayload.ReportId = uuid.New().String()
	fmt.Println(reportPayload)
	err = database.SubmitReportToDatabase(db, reportPayload)
	if err != nil {
		log.Println("Internal server error")
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(`{"error": "` + err.Error() + `"}`))
		if err != nil {
			return
		}
		return
	} else {
		w.WriteHeader(http.StatusCreated)
		_, err = w.Write([]byte(`{"status": "created with report ID: ` + reportPayload.ReportId + `"}`))
		if err != nil {
			return
		}
		return
	}
}

func GetAllSurveyAnswers(w http.ResponseWriter, r *http.Request) {
	log.Println("Get all survey answers endpoint accessed")
	w.Header().Set("Content-Type", "application/json")
	query := r.URL.Query()
	pageNum := query.Get("page")
	pageSize := query.Get("pagesize")
	surveyId := query.Get("surveyid")
	if len(pageNum) == 0 || len(pageSize) == 0 {
		log.Println("Invalid query parameters")
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(`{"error": "invalid query parameters"}`))
		if err != nil {
			return
		}
		return
	}
	pageNumInt, err := strconv.Atoi(pageNum)
	if err != nil {
		log.Println("Invalid query parameters")
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(`{"error": "invalid query parameters"}`))
		if err != nil {
			return
		}
		return
	}
	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		log.Println("Invalid query parameters")
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(`{"error": "invalid query parameters"}`))
		if err != nil {
			return
		}
		return
	}
	if len(surveyId) != 0 {
		surveyIdInt, err := strconv.Atoi(surveyId)
		if err != nil {
			log.Println("Invalid query parameters")
			w.WriteHeader(http.StatusBadRequest)
			_, err := w.Write([]byte(`{"error": "invalid query parameters"}`))
			if err != nil {
				return
			}
			return
		}
		surveyAnswers, err := database.GetAllSurveyAnswersBySurveyIDFromDatabase(db, surveyIdInt, pageNumInt, pageSizeInt)
		if err != nil {
			log.Println("Internal server error")
			w.WriteHeader(http.StatusInternalServerError)
			_, err := w.Write([]byte(`{"error": "` + err.Error() + `"}`))
			if err != nil {
				return
			}
			return
		}
		if len(surveyAnswers) == 0 {
			w.WriteHeader(http.StatusNoContent)
			_, err := w.Write([]byte(`{"error": "no survey answers found"}`))
			if err != nil {
				return
			}
			return
		}
		w.WriteHeader(http.StatusOK)
		jsonData, err := json.Marshal(surveyAnswers)
		if err != nil {
			return
		}
		_, err = w.Write(jsonData)
		if err != nil {
			return
		}
		return
	} else {
		surveyAnswers, err := database.GetAllSurveyAnswersFromDatabase(db, pageNumInt, pageSizeInt)
		if err != nil {
			log.Println("Internal server error")
			w.WriteHeader(http.StatusInternalServerError)
			_, err := w.Write([]byte(`{"error": "` + err.Error() + `"}`))
			if err != nil {
				return
			}
			return
		}
		if len(surveyAnswers) == 0 {
			w.WriteHeader(http.StatusNoContent)
			_, err := w.Write([]byte(`{"error": "no survey answers found"}`))
			if err != nil {
				return
			}
			return
		}
		w.WriteHeader(http.StatusOK)
		jsonData, err := json.Marshal(surveyAnswers)
		if err != nil {
			return
		}
		_, err = w.Write(jsonData)
		if err != nil {
			return
		}
		return
	}
}

func GetSurveyAnswerByID(w http.ResponseWriter, r *http.Request) {
	log.Println("Get survey answer by ID endpoint accessed")
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	surveyAnswerID := vars["survey_answer_id"]
	log.Println(surveyAnswerID)
	answers, err := database.GetSurveyAnswerByIDFromDatabase(db, surveyAnswerID)
	if err != nil {
		log.Println("Internal server error")
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(`{"error": "` + err.Error() + `"}`))
		if err != nil {
			return
		}
		return
	}
	w.WriteHeader(http.StatusOK)
	jsonData, err := json.Marshal(answers)
	if err != nil {
		return
	}
	_, err = w.Write(jsonData)
	if err != nil {
		return
	}
	return
}

func GetAllReports(w http.ResponseWriter, r *http.Request) {
	log.Println("Get all reports endpoint accessed")
	w.Header().Set("Content-Type", "application/json")
	query := r.URL.Query()
	pageNum := query.Get("page")
	pageSize := query.Get("pagesize")
	if len(pageNum) == 0 || len(pageSize) == 0 {
		log.Println("Invalid query parameters")
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(`{"error": "invalid query parameters"}`))
		if err != nil {
			return
		}
		return
	}
	pageNumInt, err := strconv.Atoi(pageNum)
	if err != nil {
		log.Println("Invalid query parameters")
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(`{"error": "invalid query parameters"}`))
		if err != nil {
			return
		}
		return
	}
	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		log.Println("Invalid query parameters")
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(`{"error": "invalid query parameters"}`))
		if err != nil {
			return
		}
		return
	}
	reports, err := database.GetAllReportsFromDatabase(db, pageNumInt, pageSizeInt)
	if err != nil {
		log.Println("Internal server error")
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(`{"error": "` + err.Error() + `"}`))
		if err != nil {
			return
		}
		return
	}
	if len(reports) == 0 {
		log.Println("No reports found")
		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte(`{"error": "no reports found"}`))
		if err != nil {
			return
		}
		return
	}
	w.WriteHeader(http.StatusOK)
	jsonData, err := json.Marshal(reports)
	if err != nil {
		return
	}
	_, err = w.Write(jsonData)
	if err != nil {
		return
	}
	return
}

func GetReportByID(w http.ResponseWriter, r *http.Request) {
	log.Println("Get report by ID endpoint accessed")
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	reportID := vars["report_id"]
	reports, err := database.GetReportByIDFromDatabase(db, reportID)
	if err != nil {
		log.Println("Internal server error")
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(`{"error": "` + err.Error() + `"}`))
		if err != nil {
			return
		}
		return
	}
	w.WriteHeader(http.StatusOK)
	jsonData, err := json.Marshal(reports)
	if err != nil {
		return
	}
	_, err = w.Write(jsonData)
	if err != nil {
		return
	}
	return
}

func GetStats(w http.ResponseWriter, r *http.Request) {
	log.Println("Get stats endpoint accessed")
	w.Header().Set("Content-Type", "application/json")
	stats, err := database.GetStatsFromDatabase(db)
	if err != nil {
		log.Println("Internal server error")
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(`{"error": "` + err.Error() + `"}`))
		if err != nil {
			return
		}
		return
	}
	w.WriteHeader(http.StatusOK)
	jsonData, err := json.Marshal(stats)
	if err != nil {
		return
	}
	_, err = w.Write(jsonData)
	if err != nil {
		return
	}
	return
}

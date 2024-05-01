package database

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"survey-service/payload"
	"time"
)

func ConnectToDatabase() (db *mongo.Client, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	db, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if db == nil {
		return nil, fmt.Errorf("could not connect to database: %v", err)
	}

	return db, nil
}

func SubmitAnswerToDatabase(db *mongo.Client, surveyAnswer payload.SurveyAnswerPayload) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := db.Database("survey_service_database").Collection("survey_answer_collection")
	_, err = collection.InsertOne(ctx, surveyAnswer)
	if err != nil {
		return fmt.Errorf("could not submit answer to database: %v", err)
	}

	return nil
}

func SubmitReportToDatabase(db *mongo.Client, report payload.ReportPayload) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := db.Database("survey_service_database").Collection("report_collection")
	_, err = collection.InsertOne(ctx, report)
	if err != nil {
		return fmt.Errorf("could not submit report to database: %v", err)
	}

	return nil
}

func GetAllSurveyAnswersFromDatabase(db *mongo.Client, pageNum int, pageSize int) (surveyAnswers []payload.SurveyAnswerPayload, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := db.Database("survey_service_database").Collection("survey_answer_collection")
	findOptions := options.Find()
	findOptions.SetSkip(int64((pageNum - 1) * pageSize))
	findOptions.SetLimit(int64(pageSize))
	cursor, err := collection.Find(ctx, bson.D{}, findOptions)
	if err != nil {
		return nil, fmt.Errorf("could not get survey answers from database: %v", err)
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			return
		}
	}(cursor, ctx)
	for cursor.Next(ctx) {
		var result payload.SurveyAnswerPayload
		err := cursor.Decode(&result)
		if err != nil {
			return nil, fmt.Errorf("could not decode result: %v", err)
		}
		surveyAnswers = append(surveyAnswers, result)
	}

	return surveyAnswers, nil
}

func GetSurveyAnswerByIDFromDatabase(db *mongo.Client, surveyAnswerID string) (surveyAnswer payload.SurveyAnswerPayload, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := db.Database("survey_service_database").Collection("survey_answer_collection")
	err = collection.FindOne(ctx, bson.D{{"surveyanswerid", surveyAnswerID}}).Decode(&surveyAnswer)
	if err != nil {
		return surveyAnswer, fmt.Errorf("could not get survey answer by ID from database: %v", err)
	}

	return surveyAnswer, nil
}

func GetAllReportsFromDatabase(db *mongo.Client, pageNum int, pageSize int) (reports []payload.ReportPayload, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := db.Database("survey_service_database").Collection("report_collection")
	findOptions := options.Find()
	findOptions.SetSkip(int64((pageNum - 1) * pageSize))
	findOptions.SetLimit(int64(pageSize))
	cursor, err := collection.Find(ctx, bson.D{}, findOptions)
	if err != nil {
		return nil, fmt.Errorf("could not get reports from database: %v", err)
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			return
		}
	}(cursor, ctx)
	for cursor.Next(ctx) {
		var result payload.ReportPayload
		err := cursor.Decode(&result)
		if err != nil {
			return nil, fmt.Errorf("could not decode result: %v", err)
		}
		reports = append(reports, result)
	}

	return reports, nil
}

func GetReportByIDFromDatabase(db *mongo.Client, reportID string) (report payload.ReportPayload, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := db.Database("survey_service_database").Collection("report_collection")
	err = collection.FindOne(ctx, bson.D{{"reportid", reportID}}).Decode(&report)
	if err != nil {
		return report, fmt.Errorf("could not get report by ID from database: %v", err)
	}

	return report, nil
}

func GetAllSurveyAnswersBySurveyIDFromDatabase(db *mongo.Client, surveyID int, pageNum int, pageSize int) (surveyAnswers []payload.SurveyAnswerPayload, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := db.Database("survey_service_database").Collection("survey_answer_collection")
	findOptions := options.Find()
	findOptions.SetSkip(int64((pageNum - 1) * pageSize))
	findOptions.SetLimit(int64(pageSize))
	cursor, err := collection.Find(ctx, bson.D{{"surveyid", surveyID}}, findOptions)
	if err != nil {
		return nil, fmt.Errorf("could not get survey answers by survey ID from database: %v", err)
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			return
		}
	}(cursor, ctx)
	for cursor.Next(ctx) {
		var result payload.SurveyAnswerPayload
		err := cursor.Decode(&result)
		if err != nil {
			return nil, fmt.Errorf("could not decode result: %v", err)
		}
		surveyAnswers = append(surveyAnswers, result)
	}

	return surveyAnswers, nil
}

func GetStatsFromDatabase(db *mongo.Client) (stats payload.StatsPayload, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := db.Database("survey_service_database").Collection("survey_answer_collection")
	totalSurveyAnswers, err := collection.CountDocuments(ctx, bson.D{})
	if err != nil {
		return stats, fmt.Errorf("could not get total survey answers from database: %v", err)
	}
	collection = db.Database("survey_service_database").Collection("report_collection")
	totalReports, err := collection.CountDocuments(ctx, bson.D{})
	if err != nil {
		return stats, fmt.Errorf("could not get total reports from database: %v", err)
	}
	stats.TotalSurveyAnswers = int(totalSurveyAnswers)
	stats.TotalReports = int(totalReports)

	return stats, nil
}

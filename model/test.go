package model

import (
	"database/sql"
	"time"
)

type Test struct {
	Id        int64     `json:"id"`
	UserId    int64     `json:"user_id"`
	VariantId int64     `json:"variant_id"`
	StartAt   time.Time `json:"start_at"`
}

type TestResponse struct {
	Id int64 `json:"id"`
}

type TestResultResponse struct {
	Result int `json:"result"`
}

type TestAnswerRequest struct {
	Answer string `json:"answer"`
	TestId int64  `json:"test_id"`
}

func CreateTest(db *sql.DB, userId, variantId int64) (Test, error) {
	query := `
		INSERT INTO user_test
		(user_id, variant_id, start_at)
		VALUES
		($1, $2, $3)
		RETURNING id
	`

	timeNow := time.Now()
	var id int64

	err := db.QueryRow(
		query,
		userId,
		variantId,
		timeNow,
	).Scan(&id)

	if err != nil {
		return Test{}, err
	}

	return Test{
			Id:        id,
			UserId:    userId,
			VariantId: variantId,
			StartAt:   timeNow,
		},
		nil
}

func AddTestAnswer(db *sql.DB, testId int64, answer string) error {
	query := `
		INSERT INTO test_answer
		(test_id, answer)
		VALUES
		($1, $2)
	`

	_, err := db.Exec(query, testId, answer)
	return err
}

func GetTestResult(db *sql.DB, testId, variantId int64) (TestResultResponse, error) {
	query := `
		SELECT answer FROM test_answer 
		WHERE test_id = $1 ORDER BY id ASC;
	`
	rightAnswers := 0

	rows, err := db.Query(query, testId)
	if err != nil {
		return TestResultResponse{}, err
	}
	defer rows.Close()

	answers := make([]string, 0)

	for rows.Next() {
		var answer string

		err := rows.Scan(
			&answer,
		)

		if err != nil {
			return TestResultResponse{}, err
		}
		answers = append(answers, answer)
	}

	tasks, err := GetAllTasksByVariantId(db, variantId)
	if err != nil {
		return TestResultResponse{}, err
	}

	for i, task := range tasks {
		taskData, err := GetTask(db, task, variantId)
		if err != nil {
			return TestResultResponse{}, err
		}

		if taskData.Answer == answers[i] {
			rightAnswers++
		}
	}

	percent := int((float64(rightAnswers) / float64(len(tasks))) * 100)
	return TestResultResponse{Result: percent}, nil
}

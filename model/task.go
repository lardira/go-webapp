package model

import (
	"database/sql"
	"encoding/json"
)

type Task struct {
	Id        int64    `json:"id"`
	VariantId int64    `json:"variant_id"`
	Task      string   `json:"task"`
	Answer    string   `json:"answer"`
	Options   []string `json:"options"`
}

func GetTask(db *sql.DB, id int64, variantId int64) (Task, error) {
	query := `
		SELECT id, variant_id, task, answer, options 
		FROM task WHERE id = $1 AND variant_id = $2
	`
	output := Task{}
	var optionsAsString string

	db.QueryRow(
		query,
		id,
		variantId,
	).Scan(
		&output.Id,
		&output.VariantId,
		&output.Task,
		&output.Answer,
		&optionsAsString,
	)

	json.Unmarshal([]byte(optionsAsString), &output.Options)
	return output, nil
}

func GetAllTasksByVariantId(db *sql.DB, variantId int64) ([]int64, error) {
	query := `SELECT id FROM task WHERE variant_id = $1`
	output := make([]int64, 0)

	rows, err := db.Query(query, variantId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id int64

		err := rows.Scan(
			&id,
		)

		if err != nil {
			return nil, err
		}

		output = append(output, id)
	}
	return output, nil
}

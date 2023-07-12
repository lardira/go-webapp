package model

import (
	"database/sql"
)

type Variant struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

func GetAllVariants(db *sql.DB) ([]Variant, error) {
	query := `SELECT * FROM test_variant`
	output := make([]Variant, 0)

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var variant Variant

		err := rows.Scan(&variant.Id, &variant.Name)
		if err != nil {
			return nil, err
		}

		output = append(output, variant)
	}
	return output, nil
}

// func GetVariantById(db *sql.DB, id int64) (Variant, error) {
// 	query := `
// 		SELECT id, is_auth FROM users
// 		WHERE login=$1 AND password=$2
// 	`
// 	var user User = User{}

// 	err := db.QueryRow(query, login, password).Scan(
// 		&user.Id,
// 		&user.IsAuth,
// 	)

// 	return user, err
// }

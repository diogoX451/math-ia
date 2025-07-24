package formater

import (
	"database/sql"
	"fmt"
	"math-ia/internal/db"
	"strings"
)

func FormatRows(table string, rows *sql.Rows) []string {
	cols, _ := rows.Columns()
	var result []string

	for rows.Next() {
		values := make([]interface{}, len(cols))
		valuePtrs := make([]interface{}, len(cols))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		rows.Scan(valuePtrs...)
		var parts []string
		for i, col := range cols {
			val := values[i]
			parts = append(parts, fmt.Sprintf("%s=%v", col, val))
		}

		result = append(result, fmt.Sprintf("%s: %s", table, strings.Join(parts, ", ")))
	}

	return result
}

func FormatRow(table string, row *sql.Row) (string, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", table)
	stmt, err := db.GetDB().Prepare(query)
	if err != nil {
		return "", err
	}
	defer stmt.Close()

	rows, err := stmt.Query(1) // id fictício só para pegar os metadados
	if err != nil {
		return "", err
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return "", err
	}

	values := make([]interface{}, len(cols))
	valuePtrs := make([]interface{}, len(cols))
	for i := range values {
		valuePtrs[i] = &values[i]
	}

	// Use row.Scan com ponteiros corretos
	if err := row.Scan(valuePtrs...); err != nil {
		return "", err
	}

	var sb strings.Builder
	for i, col := range cols {
		if values[i] != nil {
			sb.WriteString(fmt.Sprintf("%s: %v; ", col, values[i]))
		}
	}

	return sb.String(), nil
}

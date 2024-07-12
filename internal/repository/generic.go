package repository

import (
	"fmt"
	"strings"
)

// sqlInsert generates an SQL INSERT statement.
func sqlInsert(table string, fields []string) string {
	columns := strings.Join(fields, ", ")
	placeholders := ":" + strings.Join(fields, ", :")
	return fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", table, columns, placeholders)
}

// sqlUpdate generates an SQL UPDATE statement.
func sqlUpdate(table string, fields []string, where []string) string {
	setClauses := make([]string, len(fields))
	for i, field := range fields {
		setClauses[i] = fmt.Sprintf("%s = :%s", field, field)
	}
	whereClauses := strings.Join(where, " AND ")
	return fmt.Sprintf("UPDATE %s SET %s WHERE %s", table, strings.Join(setClauses, ", "), whereClauses)
}

// sqlReplace generates an SQL REPLACE statement.
func sqlReplace(table string, fields []string) string {
	columns := strings.Join(fields, ", ")
	placeholders := ":" + strings.Join(fields, ", :")
	return fmt.Sprintf("REPLACE INTO %s (%s) VALUES (%s)", table, columns, placeholders)
}

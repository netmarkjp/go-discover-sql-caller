package main

import (
	"testing"
)

func TestIsSQL(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"SelectStatement", "SELECT * FROM table", true},
		{"InsertStatement", "INSERT INTO table (column1) VALUES ('value1')", true},
		{"UpdateStatement", "UPDATE table SET column1 = 'value1'", true},
		{"DeleteStatement", "DELETE FROM table", true},
		{"ReplaceStatement", "REPLACE INTO table (column1) VALUES ('value1')", true},
		{"AlterStatement", "ALTER TABLE table_name ADD COLUMN column_name", true},
		{"CreateStatement", "CREATE TABLE table_name (column_name char(3))", true},
		{"DropStatement", "DROP TABLE table_name", true},
		{"TruncateStatement", "TRUNCATE TABLE table_name", true},
		{"GrantStatement", "GRANT SELECT ON table_name TO user_name", true},
		{"RevokeStatement", "REVOKE SELECT ON table_name FROM user_name", true},
		{"BeginStatement", "BEGIN TRANSACTION", true},
		{"CommitStatement", "COMMIT", true},
		{"RollbackStatement", "ROLLBACK", true},
		{"NonSQL", "This is not an sql statement", false},
		{"EmptyString", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isSQL(tt.input)
			if result != tt.expected {
				t.Errorf("isSQL(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

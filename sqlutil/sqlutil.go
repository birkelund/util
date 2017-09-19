// Package sqlutil includes database/sql utility functions.
package sqlutil

import (
	"fmt"
	"os"
	"regexp"

	"github.com/jmoiron/sqlx"
	"github.com/kbj/util/cliutil"
	"github.com/kbj/util/pqutil"
	"github.com/lib/pq"
	"github.com/pkg/errors"
)

// CreateStatement represents an SQL CREATE statement.
type CreateStatement struct {
	// Table is the table to create.
	Name string

	// SQL is the raw SQL statement.
	SQL string
}

// Parse parses a schema and returns a slice of SQL CREATE TABLE statements.
func ParseCreateTableStatements(schema []byte) (stmts []CreateStatement) {
	createTableMatcher := regexp.MustCompile("(CREATE TABLE (.*) \\((?s:.*?)\\);)")

	for _, stmt := range createTableMatcher.FindAllSubmatch(schema, -1) {
		stmts = append(stmts, CreateStatement{
			Name: string(stmt[2]),
			SQL:  string(stmt[1]),
		})
	}

	return
}

// ParseCreateTypeEnum parses a schema and returns a slice of SQL CREATE TYPE
// name AS ENUM statements.
func ParseCreateTypeEnumStatements(schema []byte) (stmts []CreateStatement) {
	createTypeEnumMatcher := regexp.MustCompile("(CREATE TYPE (.*) AS ENUM \\((?s:.*?)\\);)")

	for _, stmt := range createTypeEnumMatcher.FindAllSubmatch(schema, -1) {
		stmts = append(stmts, CreateStatement{
			Name: string(stmt[2]),
			SQL:  string(stmt[1]),
		})
	}

	return
}

// Connect connects to the database and returns a handle.
func Connect(driver, username, password, dbname string) (*sqlx.DB, error) {
	dbstr := fmt.Sprintf("user=%s password=%s dbname=%s", username, password, dbname)
	db, err := sqlx.Open(driver, dbstr)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func InitializeSchema(db *sqlx.DB, schema []byte) error {
	createTypeEnumStatements := ParseCreateTypeEnumStatements(schema)
	createTableStatements := ParseCreateTableStatements(schema)

	for _, stmt := range createTableStatements {
		err := cliutil.Run(fmt.Sprintf("DROP TABLE   %s CASCADE", stmt.Name), func() error {
			_, err := db.Exec(fmt.Sprintf("DROP TABLE %s CASCADE", stmt.Name))
			return err
		})

		if err != nil {
			if err, ok := err.(*pq.Error); ok {
				if !pqutil.IsUndefinedTableError(err) {
					fmt.Fprintf(os.Stderr, "[-] failed to drop table '%s': %v\n", stmt.Name, err)
					return err
				}
			}
		}
	}

	for _, stmt := range createTypeEnumStatements {
		err := cliutil.Run(fmt.Sprintf("DROP TYPE    %s", stmt.Name), func() error {
			_, err := db.Exec(fmt.Sprintf("DROP TYPE %s", stmt.Name))
			return err
		})

		if err != nil {
			if err, ok := err.(*pq.Error); ok {
				if !pqutil.IsUndefinedObjectError(err) {
					fmt.Fprintf(os.Stderr, "[-] failed to drop type '%s': %v\n", stmt.Name, err)
					return err
				}
			}
		}
	}

	for _, stmt := range createTypeEnumStatements {
		err := cliutil.Run(fmt.Sprintf("CREATE TYPE  %s", stmt.Name), func() error {
			_, err := db.Exec(stmt.SQL)
			return err
		})

		if err != nil {
			return errors.Wrapf(err, "failed to create type '%s'", stmt.Name)
		}
	}

	for _, stmt := range createTableStatements {
		err := cliutil.Run(fmt.Sprintf("CREATE TABLE %s", stmt.Name), func() error {
			_, err := db.Exec(stmt.SQL)
			return err
		})

		if err != nil {
			return errors.Wrapf(err, "failed to create table '%s'", stmt.Name)
		}
	}

	return nil
}

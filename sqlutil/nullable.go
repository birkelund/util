package sqlutil

import "database/sql"

// NullableInt constructs sql.NullInt64 from an integer.
func NullableInt(i int) (ni sql.NullInt64) {
	if i != 0 {
		ni.Int64 = int64(i)
		ni.Valid = true
	}

	return
}

// NullableString constructs a sql.NullString from a string.
func NullableString(s string) (ni sql.NullString) {
	if s != "" {
		ni.String = s
		ni.Valid = true
	}

	return
}

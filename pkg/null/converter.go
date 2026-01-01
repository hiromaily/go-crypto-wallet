package null

import (
	"database/sql"

	"github.com/guregu/null/v6"
)

// ConvertSQLNullTimeToNullTime converts sql.NullTime to null.Time
func ConvertSQLNullTimeToNullTime(t sql.NullTime) null.Time {
	if !t.Valid {
		return null.Time{}
	}
	return null.TimeFrom(t.Time)
}

// ConvertNullTimeToSQLNullTime converts null.Time to sql.NullTime
func ConvertNullTimeToSQLNullTime(t null.Time) sql.NullTime {
	if !t.Valid {
		return sql.NullTime{}
	}
	return sql.NullTime{Time: t.Time, Valid: true}
}

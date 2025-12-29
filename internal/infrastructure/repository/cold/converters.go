package cold

import (
	"database/sql"

	"github.com/guregu/null/v6"
)

// convertSQLNullTimeToNullTime converts sql.NullTime to null.Time
func convertSQLNullTimeToNullTime(t sql.NullTime) null.Time {
	if !t.Valid {
		return null.Time{}
	}
	return null.TimeFrom(t.Time)
}

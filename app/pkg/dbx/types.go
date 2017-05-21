package dbx

import (
	"database/sql"
	"encoding/json"

	"github.com/lib/pq"
)

// NullInt representa a nullable integer
type NullInt struct {
	sql.NullInt64
}

// MarshalJSON interface redefinition
func (r NullInt) MarshalJSON() ([]byte, error) {
	if r.Valid {
		return json.Marshal(r.Int64)
	}
	return json.Marshal(nil)
}

// NullString representa a nullable string
type NullString struct {
	sql.NullString
}

// MarshalJSON interface redefinition
func (r NullString) MarshalJSON() ([]byte, error) {
	if r.Valid {
		return json.Marshal(r.String)
	}
	return json.Marshal(nil)
}

// NullTime representa a nullable time.Time
type NullTime struct {
	pq.NullTime
}

// MarshalJSON interface redefinition
func (r NullTime) MarshalJSON() ([]byte, error) {
	if r.Valid {
		return json.Marshal(r.Time)
	}
	return json.Marshal(nil)
}

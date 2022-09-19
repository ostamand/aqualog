package helper

import "database/sql"

func ToFloatSQL(value *float64) sql.NullFloat64 {
	if value != nil {
		return sql.NullFloat64{
			Float64: *value,
			Valid:   true,
		}
	} else {
		return sql.NullFloat64{
			Valid: false,
		}
	}
}

func ToIntSQL(value *int64) sql.NullInt64 {
	if value != nil {
		return sql.NullInt64{
			Int64: *value,
			Valid: true,
		}
	} else {
		return sql.NullInt64{
			Valid: false,
		}
	}
}

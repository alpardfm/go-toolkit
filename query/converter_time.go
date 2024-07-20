package query

import (
	"time"

	"github.com/alpardfm/go-toolkit/sql"
)

func convertTimeArgs(_f interface{}) (primitiveType int8, isMany bool, args interface{}) {
	switch f := _f.(type) {
	case time.Time:
		if !f.IsZero() {
			args = f
		}
		primitiveType = Time

	case []time.Time:
		if len(f) > 0 {
			isMany = true
			args = f
		}
		primitiveType = TimeArr

	case []*time.Time:
		if len(f) > 0 {
			isMany = true
			args = f
		}
		primitiveType = TimeArr

	case []sql.NullTime:
		if len(f) > 0 {
			var _args []time.Time
			for _, r := range f {
				if r.Valid {
					_args = append(_args, r.Time)
				}
			}
			isMany = true
			args = _args
		}
		primitiveType = TimeArr

	case []*sql.NullTime:
		if len(f) > 0 {
			var _args []time.Time
			for _, r := range f {
				if r != nil {
					if r.Valid {
						_args = append(_args, r.Time)
					}
				}
			}
			isMany = true
			args = _args
		}
		primitiveType = TimeArr

	case []sql.NullDate:
		if len(f) > 0 {
			var _args []time.Time
			for _, r := range f {
				if r.Valid {
					_args = append(_args, r.Time)
				}
			}
			isMany = true
			args = _args
		}
		primitiveType = TimeArr

	case []*sql.NullDate:
		if len(f) > 0 {
			var _args []time.Time
			for _, r := range f {
				if r != nil {
					if r.Valid {
						_args = append(_args, r.Time)
					}
				}
			}
			isMany = true
			args = _args
		}
		primitiveType = TimeArr

	case sql.NullTime:
		if f.Valid {
			args = f.Time
		}
		primitiveType = Time

	case sql.NullDate:
		if f.Valid {
			args = f.Time
		}
		primitiveType = Time
	}

	return primitiveType, isMany, args
}

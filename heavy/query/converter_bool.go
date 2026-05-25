package query

import "github.com/alpardfm/go-toolkit/heavy/sql"

func convertBoolArgs(_f any) (primitiveType int8, isMany bool, args any) {
	switch f := _f.(type) {
	case []bool:
		if len(f) > 0 {
			isMany = true
			args = f
		}
		primitiveType = Bool

	case []*bool:
		if len(f) > 0 {
			isMany = true
			args = f
		}
		primitiveType = Bool

	case bool:
		if f {
			args = f
		}
		primitiveType = Bool

	case []sql.NullBool:
		if len(f) > 0 {
			var _args []bool
			for _, r := range f {
				if r.Valid {
					_args = append(_args, r.Bool)
				}
			}
			isMany = true
			args = _args
		}
		primitiveType = Bool

	case []*sql.NullBool:
		if len(f) > 0 {
			var _args []bool
			for _, r := range f {
				if r != nil {
					if r.Valid {
						_args = append(_args, r.Bool)
					}
				}
			}
			isMany = true
			args = _args
		}
		primitiveType = Bool

	case sql.NullBool:
		if f.Valid {
			args = f.Bool
		}
		primitiveType = Bool
	}

	return primitiveType, isMany, args
}

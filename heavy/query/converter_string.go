package query

import (
	"strings"

	"github.com/alpardfm/go-toolkit/heavy/sql"
)

func convertStringArgs(_f any) (primitiveType int8, isMany, isLike bool, args any) {
	switch f := _f.(type) {
	case []string:
		if len(f) > 0 {
			isMany = true
			args = f
		}
		primitiveType = StringArr

	case []*string:
		if len(f) > 0 {
			isMany = true
			args = f
		}
		primitiveType = StringArr

	case string:
		if len(f) > 0 {
			if strings.ContainsRune(f, '%') {
				isLike = true
			}
			args = f
			primitiveType = String
		}

	case []sql.NullString:
		if len(f) > 0 {
			var _args []string
			for _, r := range f {
				if r.Valid && len(r.String) > 0 {
					_args = append(_args, r.String)
				}
			}
			isMany = true
			args = _args
		}
		primitiveType = StringArr

	case []*sql.NullString:
		if len(f) > 0 {
			var _args []string
			for _, r := range f {
				if r != nil {
					if r.Valid && len(r.String) > 0 {
						_args = append(_args, r.String)
					}
				}
			}
			isMany = true
			args = _args
		}
		primitiveType = StringArr

	case sql.NullString:
		if f.Valid {
			if strings.ContainsRune(f.String, '%') {
				isLike = true
			}
			args = f.String
		}
		primitiveType = String
	}

	return primitiveType, isMany, isLike, args
}

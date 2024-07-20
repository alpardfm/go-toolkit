package query

import (
	"bytes"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/alpardfm/go-toolkit/codes"
	"github.com/alpardfm/go-toolkit/errors"
	"github.com/alpardfm/go-toolkit/sql"
	"github.com/jmoiron/sqlx"
)

const cursorField = "cursorField"

type Cursor interface {
	DecodeCursor(v string) error
	EncodeCursor() (string, error)
}

type sqlQueryBuilderOption func(*sqlClausebuilder) error

type sqlClausebuilder struct {
	rawQuery                      *bytes.Buffer
	suffixQuery                   string
	param                         reflect.Value
	args                          []interface{}
	paramToDBMap, paramToFieldMap map[string]string
	sortableParams                map[string]bool
	dbTag                         string
	paramTag                      string
	fieldTag                      string
	sortClause                    string
	paginationClause              string
	paramSortBy                   []string
	dbSortBy                      []string
	limit                         int64
	page                          int64
	db                            sql.Interface
	aliasMap                      map[string]string

	// cursors
	useCursor        bool
	rawCursor        string
	cursorArgCounter int
}

func NewSQLQueryBuilder(db sql.Interface, paramTag, dbTag string, options ...sqlQueryBuilderOption) (*sqlClausebuilder, error) {
	qb := sqlClausebuilder{
		db:              db,
		rawQuery:        bytes.NewBufferString(" WHERE 1=1"),
		args:            nil,
		fieldTag:        cursorField,
		dbTag:           dbTag,
		paramTag:        paramTag,
		paramSortBy:     nil,
		useCursor:       false,
		paramToDBMap:    make(map[string]string),
		paramToFieldMap: make(map[string]string),
		sortableParams:  make(map[string]bool),
		aliasMap:        make(map[string]string),
		limit:           0,
		page:            0,
	}

	for _, opt := range options {
		err := opt(&qb)
		if err != nil {
			return nil, err
		}
	}

	return &qb, nil
}

func (s *sqlClausebuilder) AddPrefixQuery(prefix string) *sqlClausebuilder {
	if len(prefix) > 0 {
		_, _ = s.rawQuery.WriteString(" AND " + prefix)
	}
	return s
}

func (s *sqlClausebuilder) AddSuffixQuery(suffix string) *sqlClausebuilder {
	if len(suffix) > 0 {
		s.suffixQuery = " " + suffix
	}
	return s
}

func (s *sqlClausebuilder) AddAliasPrefix(alias string, ptr interface{}) *sqlClausebuilder {
	p := reflect.ValueOf(ptr)
	if p.Kind() != reflect.Ptr {
		panic(errors.NewWithCode(codes.CodeInvalidValue, "passed interface{} should be a pointer"))
	}
	v := p.Elem()
	var address string
	if v.CanAddr() {
		address = fmt.Sprint(v.Addr().Pointer())
	}

	s.aliasMap[address] = alias
	return s
}

func (s *sqlClausebuilder) Build(param interface{}) (string, []interface{}, string, []interface{}, error) {
	// return error if the param is not a pointer or has nil value
	p := reflect.ValueOf(param)
	if p.Kind() != reflect.Ptr || p.IsNil() {
		return "", nil, "", nil, errors.NewWithCode(codes.CodeInvalidValue, "passed param should be a pointer and cannot be nil")
	}

	//copy param to struct
	s.param = p

	traverseOnParam(s.paramTag, s.dbTag, s.fieldTag, "$", "", "", s.aliasMap, s.param, s.buildSQLQueryString)

	// copy buffer to get count query
	countquery := s.rawQuery.Bytes()

	// page pagination
	// TODO: implement cursor pagination
	if !s.useCursor || len(s.rawCursor) < 1 {
		// sort must be done first before page pagination
		s.sort()
		if len(s.sortClause) > 0 {
			s.rawQuery.WriteString(s.sortClause)
		}

		s.pagePagination()
		if len(s.paginationClause) > 0 {
			s.rawQuery.WriteString(s.paginationClause)
		}
	}

	newQuery, newArgs, err := sqlx.In(s.rawQuery.String()+s.suffixQuery+";", s.args...)
	if err != nil {
		return "", nil, "", nil, err
	}
	newQuery = s.db.Leader().Rebind(newQuery)

	newCountQuery, newCountArgs, err := sqlx.In(string(countquery)+s.suffixQuery+";", s.args[0:len(s.args)-s.cursorArgCounter]...)
	if err != nil {
		return "", nil, "", nil, err
	}
	newCountQuery = s.db.Leader().Rebind(newCountQuery)

	return newQuery, newArgs, newCountQuery, newCountArgs, nil
}

func (s *sqlClausebuilder) sort() {
	for _, param := range s.paramSortBy {
		reg := regexp.MustCompile(`(?P<sign>-)?(?P<col>[a-zA-Z_\.0-9]+),?`)
		if reg.MatchString(param) {
			for _, _s := range strings.Split(param, ",") {
				direction := "ASC"
				match := reg.FindStringSubmatch(_s)
				for i, name := range reg.SubexpNames() {
					if i == 0 || name == "" {
						continue
					}
					if match != nil {
						if name == "sign" && match[i] == "-" {
							direction = "DESC"
						} else if name == "col" {
							if s.useCursor && len(s.rawCursor) > 0 && !s.sortableParams[match[i]] {
								continue
							}
							if db, ok := s.paramToDBMap[match[i]]; ok {
								db = db + " " + direction
								s.dbSortBy = append(s.dbSortBy, db)
							}
						}
					}
				}
			}
		}
	}
	if len(s.dbSortBy) > 0 {
		s.sortClause = " ORDER BY " + strings.Join(s.dbSortBy, ", ")
	}
}

func (s *sqlClausebuilder) pagePagination() {
	if s.page > 0 || s.limit > 0 {
		offset := getOffset(s.page, s.limit)
		s.paginationClause = fmt.Sprintf(" LIMIT %d, %d", offset, s.limit)
	}
}

func (s *sqlClausebuilder) buildSQLQueryString(primitiveType int8, isLike, isMany bool, fieldName, paramTag, dbTag string, args interface{}) {
	//map param to field name
	s.paramToFieldMap[paramTag] = fieldName
	// map param to db column name
	s.paramToDBMap[paramTag] = dbTag

	if dbTag == "" {
		return
	}

	if isSortBy(paramTag) {
		v, _ := args.([]string)
		if v != nil {
			s.paramSortBy = normalizeSortBy(v)
		}
		return
	}

	if isPage(paramTag) {
		v, _ := args.(int64)
		s.page = validatePage(v)
		return
	}

	if isLimit(paramTag) {
		v, _ := args.(int64)
		s.limit = validateLimit(v)
		return
	}

	// we only remap if the args is not nil
	if args == nil {
		return
	}

	if !isMany {
		if isLike {
			_, _ = s.rawQuery.WriteString(" AND " + dbTag + " LIKE " + s.getBindVar())
			s.args = append(s.args, args)
			return
		}
		if strings.Contains(paramTag, "__gte") {
			_, _ = s.rawQuery.WriteString(" AND " + dbTag + ">=" + s.getBindVar())
			s.args = append(s.args, args)
			return
		} else if strings.Contains(paramTag, "__lte") {
			_, _ = s.rawQuery.WriteString(" AND " + dbTag + "<=" + s.getBindVar())
			s.args = append(s.args, args)
			return
		} else if strings.Contains(paramTag, "__lt") {
			_, _ = s.rawQuery.WriteString(" AND " + dbTag + "<" + s.getBindVar())
			s.args = append(s.args, args)
			return
		} else if strings.Contains(paramTag, "__gt") {
			_, _ = s.rawQuery.WriteString(" AND " + dbTag + ">" + s.getBindVar())
			s.args = append(s.args, args)
			return
		} else if strings.Contains(paramTag, "__ne") {
			_, _ = s.rawQuery.WriteString(" AND " + dbTag + "<>" + s.getBindVar())
			s.args = append(s.args, args)
			return
		} else if strings.Contains(paramTag, "__opt") {
			_, _ = s.rawQuery.WriteString(" OR " + dbTag + "=" + s.getBindVar())
			s.args = append(s.args, args)
			return
		}

		_, _ = s.rawQuery.WriteString(" AND " + dbTag + "=" + s.getBindVar())
		s.args = append(s.args, args)
		return
	}

	if strings.Contains(paramTag, "__nin") {
		_, _ = s.rawQuery.WriteString(" AND " + dbTag + " NOT IN (" + s.getBindVar() + ")")
		s.args = append(s.args, args)
		return
	}

	// __ in or unstated will result IN
	_, _ = s.rawQuery.WriteString(" AND " + dbTag + " IN (" + s.getBindVar() + ")")
	s.args = append(s.args, args)
}

func (s *sqlClausebuilder) getBindVar() string {
	return "?"
}

package query

import (
	"bytes"
	"testing"
	"time"

	"github.com/alpardfm/go-toolkit/log"
	"github.com/alpardfm/go-toolkit/sql"
	"github.com/stretchr/testify/assert"
)

type TestParam struct {
	ID      int         `db:"id" param:"id" cursorField:"id"`
	Name    string      `db:"name" param:"name" cursorField:"name"`
	Weight  float64     `db:"weight" param:"weight" cursorField:"weight"`
	Height  float64     `db:"height" param:"height" cursorField:"height"`
	Details interface{} `db:"details" param:"details" cursorField:"details"`
}

type TestParamWithArray struct {
	ID     []int64   `param:"id" db:"id"`
	Name   []string  `param:"name" db:"name"`
	Length []float64 `param:"length" db:"length"`
}

type TestParamWithStringWildcard struct {
	Name string `param:"name" db:"name"`
}

type TestParamLimitAndPage struct {
	ID    int64    `param:"id" db:"id"`
	Names []string `param:"name" db:"name"`
	Limit int64    `param:"limit" db:"limit"`
	Page  int64    `param:"page" db:"page"`
}

type TestParamSortBy struct {
	ID     int64    `param:"id" db:"id"`
	SortBy []string `param:"sort_by" db:"sort_by"`
	Limit  int64    `param:"limit" db:"limit"`
	Page   int64    `param:"page" db:"page"`
}

type TestNullType struct {
	ID        sql.NullInt64   `param:"id" db:"id"`
	Name      sql.NullString  `param:"name" db:"name"`
	Length    sql.NullFloat64 `param:"length" db:"length"`
	Active    sql.NullBool    `param:"active" db:"active"`
	Birthday  sql.NullDate    `param:"birthday" db:"birthday"`
	CreatedAt sql.NullTime    `param:"created_at" db:"created_at"`
}

func TestNewSQLQueryBuilder(t *testing.T) {
	type args struct {
		db       sql.Interface
		paramTag string
		dbTag    string
		options  []sqlQueryBuilderOption
	}
	tests := []struct {
		name    string
		args    args
		want    *sqlClausebuilder
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				db:       nil,
				paramTag: "param",
				dbTag:    "db",
			},
			want: &sqlClausebuilder{
				fieldTag: "cursorField",
				paramTag: "param",
				dbTag:    "db",
				rawQuery: bytes.NewBufferString(" WHERE 1=1"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewSQLQueryBuilder(tt.args.db, tt.args.paramTag, tt.args.dbTag, tt.args.options...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewSQLQueryBuilder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_sqlClausebuilder_Build(t *testing.T) {
	t.SkipNow() // remove this if you want to run the tests
	type args struct {
		param interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    string
		want1   []interface{}
		want2   string
		want3   []interface{}
		wantErr bool
	}{
		{
			name: "all primitive param used",
			args: args{param: &TestParam{
				ID:      1,
				Name:    "test",
				Weight:  float64(123.45),
				Height:  float64(67.89),
				Details: "simple string",
			}},
			want:    " WHERE 1=1 AND id=? AND name=? AND weight=? AND height=? AND details=?;",
			want1:   []interface{}{1, "test", float64(123.45), float64(67.89), "simple string"},
			want2:   " WHERE 1=1 AND id=? AND name=? AND weight=? AND height=? AND details=?;",
			want3:   []interface{}{1, "test", float64(123.45), float64(67.89), "simple string"},
			wantErr: false,
		},
		{
			name: "all primitive param zero value",
			args: args{param: &TestParam{
				ID:      0,
				Name:    "",
				Weight:  float64(0),
				Details: "",
			}},
			want:    " WHERE 1=1;",
			want1:   nil,
			want2:   " WHERE 1=1;",
			want3:   nil,
			wantErr: false,
		},
		{
			name: "test array with IN",
			args: args{param: &TestParamWithArray{
				ID:     []int64{1, 2},
				Name:   []string{"jack", "garland"},
				Length: []float64{1.23},
			}},
			want:    " WHERE 1=1 AND id IN (?, ?) AND name IN (?, ?) AND length IN (?);",
			want1:   []interface{}{int64(1), int64(2), "jack", "garland", float64(1.23)},
			want2:   " WHERE 1=1 AND id IN (?, ?) AND name IN (?, ?) AND length IN (?);",
			want3:   []interface{}{int64(1), int64(2), "jack", "garland", float64(1.23)},
			wantErr: false,
		},
		{
			name: "test empty array with IN",
			args: args{param: &TestParamWithArray{
				ID:     []int64{},
				Name:   []string{},
				Length: []float64{},
			}},
			want:    " WHERE 1=1;",
			want1:   nil,
			want2:   " WHERE 1=1;",
			want3:   nil,
			wantErr: false,
		},
		{
			name: "test string LIKE",
			args: args{param: &TestParamWithStringWildcard{
				Name: "garland%",
			}},
			want:    " WHERE 1=1 AND name LIKE ?;",
			want1:   []interface{}{"garland%"},
			want2:   " WHERE 1=1 AND name LIKE ?;",
			want3:   []interface{}{"garland%"},
			wantErr: false,
		},
		{
			name: "test limit and page",
			args: args{param: &TestParamLimitAndPage{
				ID:    1,
				Names: []string{"jack", "garland"},
				Limit: 10,
				Page:  1,
			}},
			want:    " WHERE 1=1 AND id=? AND name IN (?, ?) LIMIT 0, 10;",
			want1:   []interface{}{int64(1), "jack", "garland"},
			want2:   " WHERE 1=1 AND id=? AND name IN (?, ?);",
			want3:   []interface{}{int64(1), "jack", "garland"},
			wantErr: false,
		},
		{
			name: "test sort by",
			args: args{param: &TestParamSortBy{
				ID:     1,
				SortBy: []string{"id"},
				Limit:  10,
				Page:   1,
			}},
			want:    " WHERE 1=1 AND id=? ORDER BY id ASC LIMIT 0, 10;",
			want1:   []interface{}{int64(1)},
			want2:   " WHERE 1=1 AND id=?;",
			want3:   []interface{}{int64(1)},
			wantErr: false,
		},
		{
			name: "test null type",
			args: args{param: &TestNullType{
				ID:        sql.NullInt64{Valid: true},
				Name:      sql.NullString{Valid: true},
				Length:    sql.NullFloat64{Valid: true},
				Active:    sql.NullBool{Valid: true},
				Birthday:  sql.NullDate{Valid: true},
				CreatedAt: sql.NullTime{Valid: true},
			}},
			want:    " WHERE 1=1 AND id=? AND name=? AND length=? AND active=? AND birthday=? AND created_at=?;",
			want1:   []interface{}{int64(0), "", float64(0), false, time.Time{}, time.Time{}},
			want2:   " WHERE 1=1 AND id=? AND name=? AND length=? AND active=? AND birthday=? AND created_at=?;",
			want3:   []interface{}{int64(0), "", float64(0), false, time.Time{}, time.Time{}},
			wantErr: false,
		},
		{
			name: "test null type invalid",
			args: args{param: &TestNullType{
				ID:        sql.NullInt64{},
				Name:      sql.NullString{},
				Length:    sql.NullFloat64{},
				Active:    sql.NullBool{},
				Birthday:  sql.NullDate{},
				CreatedAt: sql.NullTime{},
			}},
			want:    " WHERE 1=1;",
			want1:   nil,
			want2:   " WHERE 1=1;",
			want3:   nil,
			wantErr: false,
		},
		// TODO: add test case for partial param
		// TODO: add test case for OR clause
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := sql.Init(sql.Config{
				Driver: "mysql",
				Leader: sql.ConnConfig{
					Host:     "localhost",
					Port:     3306,
					DB:       "",
					User:     "root",
					Password: "password",
				},
				Follower: sql.ConnConfig{
					Host:     "localhost",
					Port:     3306,
					DB:       "",
					User:     "root",
					Password: "password",
				},
			}, log.Init(log.Config{Level: "debug"}))

			qBuilder, err := NewSQLQueryBuilder(db, "param", "db")
			if err != nil {
				t.Error(err)
			}

			got, got1, got2, got3, err := qBuilder.Build(tt.args.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("sqlClausebuilder.Build() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.want1, got1)
			assert.Equal(t, tt.want2, got2)
			assert.Equal(t, tt.want3, got3)
		})
	}
}

func Test_sqlClausebuilder_AddPrefixQuery(t *testing.T) {
	t.SkipNow() // remove this if you want to run the tests
	type args struct {
		param       interface{}
		prefixquery string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		want1   []interface{}
		want2   string
		want3   []interface{}
		wantErr bool
	}{
		{
			name: "test prefix query",
			args: args{param: &TestParamLimitAndPage{
				ID:    1,
				Names: []string{"jack", "garland"},
				Limit: 10,
				Page:  1,
			}, prefixquery: "active = 1"},
			want:    " WHERE 1=1 AND active = 1 AND id=? AND name IN (?, ?) LIMIT 0, 10;",
			want1:   []interface{}{int64(1), "jack", "garland"},
			want2:   " WHERE 1=1 AND active = 1 AND id=? AND name IN (?, ?);",
			want3:   []interface{}{int64(1), "jack", "garland"},
			wantErr: false,
		},
		{
			name: "test prefix query empty",
			args: args{param: &TestParamLimitAndPage{
				ID:    1,
				Names: []string{"jack", "garland"},
				Limit: 10,
				Page:  1,
			}, prefixquery: ""},
			want:    " WHERE 1=1 AND id=? AND name IN (?, ?) LIMIT 0, 10;",
			want1:   []interface{}{int64(1), "jack", "garland"},
			want2:   " WHERE 1=1 AND id=? AND name IN (?, ?);",
			want3:   []interface{}{int64(1), "jack", "garland"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := sql.Init(sql.Config{
				Driver: "mysql",
				Leader: sql.ConnConfig{
					Host:     "localhost",
					Port:     3306,
					DB:       "",
					User:     "root",
					Password: "password",
				},
				Follower: sql.ConnConfig{
					Host:     "localhost",
					Port:     3306,
					DB:       "",
					User:     "root",
					Password: "password",
				},
			}, log.Init(log.Config{Level: "debug"}))

			qBuilder, err := NewSQLQueryBuilder(db, "param", "db")
			if err != nil {
				t.Error(err)
			}

			qBuilder.AddPrefixQuery(tt.args.prefixquery)

			got, got1, got2, got3, err := qBuilder.Build(tt.args.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("sqlClausebuilder.AddPrefixQuery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.want1, got1)
			assert.Equal(t, tt.want2, got2)
			assert.Equal(t, tt.want3, got3)
		})
	}
}

func Test_sqlClausebuilder_AddAliasPrefix(t *testing.T) {
	t.SkipNow() // remove this if you want to run the tests
	type args struct {
		param       interface{}
		aliasprefix string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		want1   []interface{}
		want2   string
		want3   []interface{}
		wantErr bool
	}{
		{
			name: "test alias prefix",
			args: args{param: &TestParamLimitAndPage{
				ID:    1,
				Names: []string{"jack", "garland"},
				Limit: 10,
				Page:  1,
			}, aliasprefix: "idk"},
			want:    " WHERE 1=1 AND idk.id=? AND idk.name IN (?, ?) LIMIT 0, 10;",
			want1:   []interface{}{int64(1), "jack", "garland"},
			want2:   " WHERE 1=1 AND idk.id=? AND idk.name IN (?, ?);",
			want3:   []interface{}{int64(1), "jack", "garland"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := sql.Init(sql.Config{
				Driver: "mysql",
				Leader: sql.ConnConfig{
					Host:     "localhost",
					Port:     3306,
					DB:       "",
					User:     "root",
					Password: "password",
				},
				Follower: sql.ConnConfig{
					Host:     "localhost",
					Port:     3306,
					DB:       "",
					User:     "root",
					Password: "password",
				},
			}, log.Init(log.Config{Level: "debug"}))

			qBuilder, err := NewSQLQueryBuilder(db, "param", "db")
			if err != nil {
				t.Error(err)
			}

			qBuilder.AddAliasPrefix(tt.args.aliasprefix, tt.args.param)

			got, got1, got2, got3, err := qBuilder.Build(tt.args.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("sqlClausebuilder.AddAliasPrefix() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.want1, got1)
			assert.Equal(t, tt.want2, got2)
			assert.Equal(t, tt.want3, got3)
		})
	}
}

func Test_sqlClausebuilder_AddSuffixQuery(t *testing.T) {
	t.SkipNow() // remove this if you want to run the tests
	type args struct {
		param       interface{}
		suffixquery string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		want1   []interface{}
		want2   string
		want3   []interface{}
		wantErr bool
	}{
		{
			name: "test suffix query",
			args: args{param: &TestParamLimitAndPage{
				ID:    1,
				Names: []string{"jack", "garland"},
				Limit: 10,
				Page:  1,
			}, suffixquery: "GROUP BY name"},
			want:    " WHERE 1=1 AND id=? AND name IN (?, ?) LIMIT 0, 10 GROUP BY name;",
			want1:   []interface{}{int64(1), "jack", "garland"},
			want2:   " WHERE 1=1 AND id=? AND name IN (?, ?) GROUP BY name;",
			want3:   []interface{}{int64(1), "jack", "garland"},
			wantErr: false,
		},
		{
			name: "test suffix query empty",
			args: args{param: &TestParamLimitAndPage{
				ID:    1,
				Names: []string{"jack", "garland"},
				Limit: 10,
				Page:  1,
			}, suffixquery: ""},
			want:    " WHERE 1=1 AND id=? AND name IN (?, ?) LIMIT 0, 10;",
			want1:   []interface{}{int64(1), "jack", "garland"},
			want2:   " WHERE 1=1 AND id=? AND name IN (?, ?);",
			want3:   []interface{}{int64(1), "jack", "garland"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := sql.Init(sql.Config{
				Driver: "mysql",
				Leader: sql.ConnConfig{
					Host:     "localhost",
					Port:     3306,
					DB:       "",
					User:     "root",
					Password: "password",
				},
				Follower: sql.ConnConfig{
					Host:     "localhost",
					Port:     3306,
					DB:       "",
					User:     "root",
					Password: "password",
				},
			}, log.Init(log.Config{Level: "debug"}))

			qBuilder, err := NewSQLQueryBuilder(db, "param", "db")
			if err != nil {
				t.Error(err)
			}

			qBuilder.AddSuffixQuery(tt.args.suffixquery)

			got, got1, got2, got3, err := qBuilder.Build(tt.args.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("sqlClausebuilder.AddSuffixQuery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.want1, got1)
			assert.Equal(t, tt.want2, got2)
			assert.Equal(t, tt.want3, got3)
		})
	}
}

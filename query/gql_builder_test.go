package query

import (
	"errors"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/alpardfm/go-toolkit/format"
)

func TestGQLBuilder(t *testing.T) {
	type Params struct {
		ID        int64    `gql:"id"`
		Name      string   `gql:"name"`
		Hobbies   []string `gql:"hobbies"`
		ChildsAge []int64  `gql:"childs_age"`
		Money     float64  `gql:"money"`
		IDK       string
		Birthdays time.Time `gql:"birthdays"`
		IsMarried bool      `gql:"isMarried"`
	}

	tests := []struct {
		name           string
		params         any
		wantErr        string
		wantResult     string
		exampleOptions []func(*gqlParamsBuilder) error
		fieldOptions   map[string]func(v reflect.Value) (reflect.Value, error)
		typeOptions    map[string]func(v reflect.Value) (reflect.Value, error)
	}{
		{
			name: "Test with all zero value",
			params: Params{
				ID:        0,
				Name:      "",
				Hobbies:   []string{},
				ChildsAge: []int64{},
				Money:     0.0,
			},
			wantResult: "",
		},
		{
			name:    "Test with error option",
			params:  nil,
			wantErr: "this option should be error",
			exampleOptions: []func(*gqlParamsBuilder) error{
				func(gpb *gqlParamsBuilder) error {
					return errors.New("this option should be error")
				},
			},
		},
		{
			name:    "Test with nil build parameter",
			params:  nil,
			wantErr: "params: must be a struct",
		},
		{
			name:    "Test with not struct build parameter",
			params:  1,
			wantErr: "params: must be a struct",
		},
		{
			name:    "Test with pointer struct build parameter",
			params:  &Params{},
			wantErr: "params: cannot be a pointer",
		},
		{
			name: "Test success build",
			params: Params{
				ID:        100,
				Name:      "My Name",
				Hobbies:   []string{"Coding", "Eating"},
				ChildsAge: []int64{19, 20, 21},
				Money:     10000.0,
				IDK:       "haha i don't know",
			},
			wantResult: ` id: 100  name: "My Name"  hobbies: ["Coding" "Eating"]  childs_age: [19 20 21]  money: 10000 `,
		},
		{
			name: "Test success build but not retrieve field if field doesn't have `gql` tag",
			params: Params{
				IDK: "haha i don't know",
			},
			wantResult: ``,
		},
		{
			name: "Test test field options",
			params: Params{
				Birthdays: func() time.Time { t, _ := time.Parse(format.DayMonthYearHourMinSec, "01/01/2002 01:00:00"); return t }(),
			},
			fieldOptions: map[string]func(v reflect.Value) (reflect.Value, error){
				"birthdays": func(v reflect.Value) (reflect.Value, error) {
					return reflect.ValueOf(v.Interface().(time.Time).Format(format.DayMonthYearHourMinSec)), nil
				},
			},
			wantResult: ` birthdays: "01/01/2002 01:00:00" `,
		},
		{
			name: "Test test field options with other fields",
			params: Params{
				Name:      "Irda Islakhu Afa",
				Birthdays: func() time.Time { t, _ := time.Parse(format.DayMonthYearHourMinSec, "01/01/2002 01:00:00"); return t }(),
			},
			fieldOptions: map[string]func(v reflect.Value) (reflect.Value, error){
				"birthdays": func(v reflect.Value) (reflect.Value, error) {
					return reflect.ValueOf(v.Interface().(time.Time).Format(format.DayMonthYearHourMinSec)), nil
				},
			},
			wantResult: ` name: "Irda Islakhu Afa"  birthdays: "01/01/2002 01:00:00" `,
		},
		{
			name: "Test test type options",
			params: Params{
				Birthdays: func() time.Time { t, _ := time.Parse(format.DayMonthYearHourMinSec, "01/01/2002 01:00:00"); return t }(),
			},
			typeOptions: map[string]func(v reflect.Value) (reflect.Value, error){
				"time.Time": func(v reflect.Value) (reflect.Value, error) {
					return reflect.ValueOf(v.Interface().(time.Time).Format(format.DayMonthYearHourMinSec)), nil
				},
			},
			wantResult: ` birthdays: "01/01/2002 01:00:00" `,
		},
		{
			name: "Test test type options with other fields",
			params: Params{
				Name:      "Irda Islakhu Afa",
				Birthdays: func() time.Time { t, _ := time.Parse(format.DayMonthYearHourMinSec, "01/01/2002 01:00:00"); return t }(),
			},
			typeOptions: map[string]func(v reflect.Value) (reflect.Value, error){
				"time.Time": func(v reflect.Value) (reflect.Value, error) {
					return reflect.ValueOf(v.Interface().(time.Time).Format(format.DayMonthYearHourMinSec)), nil
				},
			},
			wantResult: ` name: "Irda Islakhu Afa"  birthdays: "01/01/2002 01:00:00" `,
		},
		{
			name: "Test boolean data type true",
			params: Params{
				IsMarried: true,
			},
			wantErr:    "",
			wantResult: ` isMarried: true `,
		},
		{
			name: "Test boolean data type false",
			params: Params{
				IsMarried: false,
			},
			wantErr:    "",
			wantResult: ``,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gqlBuilder, err := NewGQLParamsBuilder("gql", tt.exampleOptions...)
			if err != nil {
				if !strings.EqualFold(tt.wantErr, err.Error()) {
					t.Errorf("NewGQLParamsBuilder() error: '%v', want error: '%v'", err.Error(), tt.wantErr)
				}
			} else {

				// assign field options
				gqlBuilder.AddFieldOptions(tt.fieldOptions)

				// assign type options
				gqlBuilder.AddTypeOptions(tt.typeOptions)

				result, err := gqlBuilder.Build(tt.params)
				if err != nil && !strings.EqualFold(tt.wantErr, err.Error()) {
					t.Errorf("NewGQLParamsBuilder() error: '%v', want error: '%v'", err.Error(), tt.wantErr)
				}

				if !strings.EqualFold(tt.wantResult, result) {
					t.Errorf("NewGQLParamsBuilder() result: '%v', want result: '%v'", result, tt.wantResult)
				}
			}
		})
	}
}

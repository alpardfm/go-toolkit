package convert

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/alpardfm/go-toolkit/format"
)

func TestToArrInt64(t *testing.T) {
	type args struct {
		i interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    []int64
		wantErr bool
	}{
		{
			name:    "error to []int64",
			args:    args{i: fmt.Errorf("huehue")},
			want:    []int64{},
			wantErr: true,
		},
		{
			name:    "[]error to []int64",
			args:    args{i: []error{fmt.Errorf("huehue")}},
			want:    []int64{},
			wantErr: true,
		},
		{
			name:    "int64 to []int64",
			args:    args{i: int64(5)},
			want:    []int64{5},
			wantErr: false,
		},
		{
			name:    "[]string to []int64",
			args:    args{i: []string{"1", "2", "3"}},
			want:    []int64{1, 2, 3},
			wantErr: false,
		},
		{
			name:    "[]float64 to []int64",
			args:    args{i: []float64{1.43, 2.3, 3.9}},
			want:    []int64{1, 2, 3},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToArrInt64(tt.args.i)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToArrInt64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToArrInt64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIntToChar(t *testing.T) {
	tests := []struct {
		name string
		args int
		want string
	}{
		{
			name: "1 to A",
			args: 1,
			want: "A",
		},
		{
			name: "3 to C",
			args: 3,
			want: "C",
		},
		{
			name: "25 to Y",
			args: 25,
			want: "Y",
		},
		{
			name: "26 to Z",
			args: 26,
			want: "Z",
		},
		{
			name: "0 to empty string",
			args: 0,
			want: "",
		},
		{
			name: "-1 to empty string",
			args: -1,
			want: "",
		},
		{
			name: "27 to AA",
			args: 27,
			want: "AA",
		},
		{
			name: "51 to AY",
			args: 51,
			want: "AY",
		},
		{
			name: "52 to AZ",
			args: 52,
			want: "AZ",
		},
		{
			name: "53 to BA",
			args: 53,
			want: "BA",
		},
		{
			name: "702 to ZZ",
			args: 702,
			want: "ZZ",
		},
		{
			name: "703 to AAA",
			args: 703,
			want: "AAA",
		},
		{
			name: "703 to BXQ",
			args: 1993,
			want: "BXQ",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IntToChar(tt.args); got != tt.want {
				t.Errorf("IntToChar() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPascalCaseToCamelCase(t *testing.T) {
	tests := []struct {
		name string
		args string
		want string
	}{
		{
			name: "Test empty string",
			args: "",
			want: "",
		},
		{
			name: "Pascal to Camel",
			args: "PasCal",
			want: "pasCal",
		},
		{
			name: "Camel to Camel",
			args: "caMel",
			want: "caMel",
		},
		{
			name: "Test single uppercase string",
			args: "A",
			want: "a",
		},
		{
			name: "Test single lowercase string",
			args: "a",
			want: "a",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PascalCaseToCamelCase(tt.args); got != tt.want {
				t.Errorf("PascalCaseToCamelCase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToFloat64(t *testing.T) {
	flo := float64(123.43)
	pointerFlo := &flo
	type args struct {
		i interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    float64
		wantErr bool
	}{
		{
			name:    "error",
			args:    args{i: fmt.Errorf("huehue")},
			want:    0,
			wantErr: true,
		},
		{
			name:    "[]error",
			args:    args{i: []error{fmt.Errorf("huehue")}},
			want:    1, // print the length
			wantErr: false,
		},
		{
			name:    "random string",
			args:    args{i: "huehue"},
			want:    0,
			wantErr: true,
		},
		{
			name:    "[]float64",
			args:    args{i: []float64{5, 6, 7}},
			want:    3, // print the length
			wantErr: false,
		},
		{
			name:    "int64",
			args:    args{i: int64(12)},
			want:    12,
			wantErr: false,
		},
		{
			name:    "pointer to float",
			args:    args{i: flo},
			want:    123.43,
			wantErr: false,
		},
		{
			name:    "float dereference",
			args:    args{i: *pointerFlo},
			want:    123.43,
			wantErr: false,
		},
		{
			name:    "string float",
			args:    args{i: "123.43"},
			want:    123.43,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToFloat64(tt.args.i)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToFloat64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ToFloat64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCopyStruct(t *testing.T) {
	type structA struct {
		V2 int64
		V1 *int
		V4 string
		V3 *string
		V6 float64
		V5 *float64
		V7 *bool
		V8 bool
		V9 struct {
			V10 string
		}
		V11 string
	}

	type structB struct {
		V1 int64
		V2 int64
		V3 string
		V4 string
		V5 float64
		V6 float64
		V7 bool
		V8 bool
		V9 struct {
			V10 string
		}
		V11 time.Time
	}

	a := structA{
		V1: func() *int { var i int = 100; return &i }(),
		V2: 100,
		V3: func() *string { var s string = "III"; return &s }(),
		V4: "III",
		V5: func() *float64 { var i float64 = 100.300; return &i }(),
		V6: 100.300,
		V7: func() *bool { var i bool = true; return &i }(),
		V8: true,
		V9: struct{ V10 string }{
			V10: "V10 mwehehehe",
		},
		V11: "21/04/2023 09:58:07.097123638",
	}

	timeToStringOption := func(value reflect.Value) (reflect.Value, error) {
		if !value.IsValid() {
			return reflect.ValueOf(time.Time{}), nil
		}

		timeValue, err := format.TimeParseWithDefaultFormat(value.String())
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(timeValue), nil
	}

	stringToTimeOption := func(value reflect.Value) (reflect.Value, error) {
		v := value.String()
		timeValue, err := format.TimeParseWithDefaultFormat(v)
		if err != nil {
			return reflect.Value{}, err
		}

		return reflect.ValueOf(timeValue), nil
	}

	intToStringOption := func(srcValue reflect.Value) (dstValue reflect.Value, err error) {
		return reflect.ValueOf(fmt.Sprintf("%v", srcValue.Int())), nil
	}

	tests := []struct {
		name         string
		srcArgs      any
		dstArgs      any
		fieldOptions map[string]func(srcValue reflect.Value) (dstValue reflect.Value, err error)
		typeOptions  map[string]func(srcValue reflect.Value) (dstValue reflect.Value, err error)
		wantMessage  string
	}{
		{
			name:        "src: error parameter String",
			srcArgs:     "a",
			dstArgs:     structB{},
			wantMessage: "src: parameter must be a struct but given type is String",
		},
		{
			name:        "src: error parameter Int",
			srcArgs:     100,
			dstArgs:     structB{},
			wantMessage: "src: parameter must be a struct but given type is Int",
		},
		{
			name:        "src: error parameter bool",
			srcArgs:     true,
			dstArgs:     structB{},
			wantMessage: "src: parameter must be a struct but given type is bool",
		},
		{
			name:        "src: error parameter Float64",
			srcArgs:     10.10,
			dstArgs:     structB{},
			wantMessage: "src: parameter must be a struct but given type is Float64",
		},
		{
			name:        "src: error parameter ptr",
			srcArgs:     func() *string { s := "s"; return &s }(),
			dstArgs:     structB{},
			wantMessage: "src: parameter must be a struct but given type is ptr",
		},
		{
			name:        "dst: error parameter String",
			srcArgs:     a,
			dstArgs:     "structB{}",
			wantMessage: "dst: parameter must be a pointer but given type is String",
		},
		{
			name:        "dst: error parameter Int",
			srcArgs:     a,
			dstArgs:     100,
			wantMessage: "dst: parameter must be a pointer but given type is Int",
		},
		{
			name:        "dst: error parameter bool",
			srcArgs:     a,
			dstArgs:     true,
			wantMessage: "dst: parameter must be a pointer but given type is bool",
		},
		{
			name:        "dst: error parameter Float64",
			srcArgs:     a,
			dstArgs:     10.10,
			wantMessage: "dst: parameter must be a pointer but given type is Float64",
		},
		{
			name:        "dst: error parameter pointer of string",
			srcArgs:     a,
			dstArgs:     func() *string { s := "s"; return &s }(),
			wantMessage: "dst: parameter must be a pointer but given type is pointer of string",
		},
		{
			name:        "dst: error parameter Struct",
			srcArgs:     a,
			dstArgs:     structB{},
			wantMessage: "dst: parameter must be a pointer but given type is Struct",
		},
		{
			name:    "CopyStruct() success with field options",
			srcArgs: a,
			dstArgs: &structB{},
			fieldOptions: map[string]func(value reflect.Value) (reflect.Value, error){
				"V11": stringToTimeOption,
				"V12": intToStringOption,
			},
		},
		{
			name:    "CopyStruct() success with different field type with type options",
			srcArgs: a,
			dstArgs: &structB{},
			typeOptions: map[string]func(value reflect.Value) (reflect.Value, error){
				"Time": timeToStringOption,
			},
			fieldOptions: map[string]func(srcValue reflect.Value) (dstValue reflect.Value, err error){
				"struct": func(value reflect.Value) (reflect.Value, error) {
					v := value.Interface().(struct{ V10 string })
					return reflect.ValueOf(v), nil
				},
			},
		},
		{
			name:    "CopyStruct() success with different field type with type options and field options",
			srcArgs: a,
			dstArgs: &structB{},
			typeOptions: map[string]func(value reflect.Value) (reflect.Value, error){
				"V9": func(value reflect.Value) (reflect.Value, error) {
					v := value.Interface().(struct{ V10 string })
					return reflect.ValueOf(v), nil
				},
			},
			fieldOptions: map[string]func(srcValue reflect.Value) (dstValue reflect.Value, err error){
				"V11": stringToTimeOption,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CopyStruct(tt.srcArgs, tt.dstArgs, tt.fieldOptions, tt.typeOptions)
			if err != nil && !strings.EqualFold(tt.wantMessage, err.Error()) {
				t.Errorf("CopyStruct() error message: '%v', want error message: '%v'", err.Error(), tt.wantMessage)
				return
			}

			if err == nil {
				// typeOfsrcArgs := reflect.TypeOf(tt.srcArgs)
				typeOfDstArgs := reflect.TypeOf(tt.dstArgs).Elem() // it's definitely a pointer from a struct
				valueOfSrcArgs := reflect.ValueOf(tt.srcArgs)
				valueOfToArgs := reflect.ValueOf(tt.dstArgs).Elem() // it's definitely a pointer from a struct
				for i := 0; i < typeOfDstArgs.NumField(); i++ {
					fieldValueDstArgs := valueOfToArgs.Field(i)
					if fieldValueDstArgs.Kind() == reflect.Ptr && !fieldValueDstArgs.IsNil() {
						fieldValueDstArgs = fieldValueDstArgs.Elem()
					}

					fieldValueSrcArgs := valueOfSrcArgs.FieldByName(typeOfDstArgs.Field(i).Name)
					if fieldValueSrcArgs.Kind() == reflect.Ptr && !fieldValueSrcArgs.IsNil() {
						fieldValueSrcArgs = fieldValueSrcArgs.Elem()
					}

					switch {
					// int
					case fieldValueDstArgs.CanInt():
						if fieldValueDstArgs.Int() != fieldValueSrcArgs.Int() {
							t.Errorf("CopyStruct() expected '%v', got '%v'", fieldValueSrcArgs.Int(), fieldValueDstArgs.Int())
						}

					// float
					case fieldValueDstArgs.CanFloat():
						if fieldValueDstArgs.Float() != fieldValueSrcArgs.Float() {
							t.Errorf("CopyStruct() expected '%v', got '%v'", fieldValueSrcArgs.Float(), fieldValueDstArgs.Float())
						}

					// string
					case fieldValueDstArgs.Kind() == reflect.String:
						if fieldValueDstArgs.String() != fieldValueSrcArgs.String() {
							t.Errorf("CopyStruct() expected '%v', got '%v'", fieldValueSrcArgs.String(), fieldValueDstArgs.String())
						}

					// bool
					case fieldValueDstArgs.Kind() == reflect.Bool:
						if fieldValueDstArgs.Bool() != fieldValueSrcArgs.Bool() {
							t.Errorf("CopyStruct() expected '%v', got '%v'", fieldValueSrcArgs.Bool(), fieldValueDstArgs.Bool())
						}

					// struct: nested struct currently not supported

					// time
					case fieldValueDstArgs.Type().Name() == "Time":
						{
							timeString := fieldValueSrcArgs.String()
							timeField := fieldValueDstArgs.Interface().(time.Time).Format(format.DayMonthYearHourMinSecMilisec)
							if timeField != timeString {
								t.Errorf("CopyStruct() expected '%v', got '%v'", timeField, timeString)
							}
						}
					}
				}
			}
		})
	}

}

func Test_ToSafeValue(t *testing.T) {
	type DataType int
	const (
		INT = DataType(iota + 1)
		FLOAT
		STRING
		STRUCT
	)

	type (
		exampleStruct struct {
			Anything string
		}

		params struct {
			dataType DataType
			value    any
		}

		want struct {
			value any
		}

		test struct {
			name   string
			params params
			want   want
		}
	)
	tests := []test{
		{
			name:   "int nil ptr",
			params: params{dataType: INT, value: nil},
			want:   want{value: 0},
		},
		{
			name:   "int ptr value 10",
			params: params{dataType: INT, value: ToPtr(10)},
			want:   want{value: 10},
		},
		{
			name:   "int ptr ptr value 10",
			params: params{dataType: INT, value: ToPtr(ToPtr(10))},
			want:   want{value: 10},
		},
		{
			name:   "string ptr",
			params: params{dataType: STRING, value: nil},
			want:   want{value: ""},
		},
		{
			name:   "string ptr value 'str'",
			params: params{dataType: STRING, value: ToPtr("str")},
			want:   want{value: "str"},
		},
		{
			name:   "struct ptr",
			params: params{dataType: STRUCT, value: nil},
			want:   want{value: exampleStruct{}},
		},
		{
			name:   "struct ptr value",
			params: params{dataType: STRUCT, value: ToPtr(exampleStruct{Anything: "haha"})},
			want:   want{value: exampleStruct{Anything: "haha"}},
		},
		{
			name:   "struct ptr value empty str",
			params: params{dataType: STRUCT, value: ToPtr("")},
			want:   want{value: exampleStruct{}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result any
			switch tt.params.dataType {
			case INT:
				result = ToSafeValue[int](tt.params.value)
			case STRING:
				result = ToSafeValue[string](tt.params.value)
			case FLOAT:
				result = ToSafeValue[float64](tt.params.value)
			case STRUCT:
				result = ToSafeValue[exampleStruct](tt.params.value)
			default:
				t.Fatalf("data type not supported!")
			}

			if result != tt.want.value {
				t.Fatalf("want result is '%v' but got '%v'", tt.want.value, result)
			}
		})
	}
}

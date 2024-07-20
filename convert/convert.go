package convert

import (
	"reflect"
	"strings"

	"github.com/alpardfm/go-toolkit/codes"
	"github.com/alpardfm/go-toolkit/errors"
	"github.com/cstockton/go-conv"
)

func ToInt64(i interface{}) (int64, error) {
	val, err := conv.Int64(i)
	if err != nil {
		return val, errors.NewWithCode(codes.CodeInvalidValue, "failed to parse value: %v. with err: %v", i, err)
	}

	return val, nil
}

func ToArrInt64(i interface{}) ([]int64, error) {
	val := make([]int64, 0)
	// if the input is not a slice, convert the input to slice with one member
	if reflect.TypeOf(i).Kind() != reflect.Slice {
		singleVal, err := conv.Int64(i)
		if err != nil {
			return []int64{}, errors.NewWithCode(codes.CodeInvalidValue, "failed to parse value: %v. with err: %v", i, err)
		}

		return append(val, singleVal), nil
	}

	s := reflect.ValueOf(i)
	for i := 0; i < s.Len(); i++ {
		singleVal, err := conv.Int64(s.Index(i).Interface())
		if err != nil {
			return []int64{}, errors.NewWithCode(codes.CodeInvalidValue, "failed to parse value: %v. with err: %v", i, err)
		}

		val = append(val, singleVal)
	}

	return val, nil
}

func ToFloat64(i interface{}) (float64, error) {
	v := reflect.ValueOf(i)
	v = reflect.Indirect(v)

	if !v.IsValid() {
		return 0, errors.NewWithCode(codes.CodeInvalidValue, "failed to parse value: %v. Invalid value", i)
	}

	floatType := reflect.TypeOf(float64(0))
	if v.Type().ConvertibleTo(floatType) {
		fv := v.Convert(floatType)
		return fv.Float(), nil
	}

	// if value is not convertible in Go vanilla, convert using library
	val, err := conv.Float64(i)
	if err != nil {
		return val, errors.NewWithCode(codes.CodeInvalidValue, "failed to parse value: %v. with err: %v", i, err)
	}

	return val, nil
}

func ToString(i interface{}) (string, error) {
	val, err := conv.String(i)
	if err != nil {
		return val, errors.NewWithCode(codes.CodeInvalidValue, "failed to parse value: %v. with err: %v", i, err)
	}

	return val, nil
}

// Convert integer to its character representation (eg. 1=A 3=C).
// Support multiple character for number greater than 26 (eg. 27=AA 703=AAA)
func IntToChar(i int) string {
	char := ""
	i--

	if i < 0 {
		return char
	}

	if firstcharint := i / 26; firstcharint > 0 {
		char += IntToChar(firstcharint)
		char += string(rune('A' + (i % 26)))
	} else {
		char += string(rune('A' + i))
	}

	return char
}

// Convert Pascal Case string (PascalCase) to Camel Case string (camelCase)
func PascalCaseToCamelCase(s string) string {
	if len(s) < 1 {
		return ""
	}

	first := strings.ToLower(s[:1])
	rest := s[1:]
	return first + rest
}

// Copy value between 2 struct. Here you can use 2 options:
//
//	fieldOpts : you can create your own logic how to copy value between 2 fields with different type (ex. string to time.Time) by field name
//	 - as example if field of entity is "CreatedAt" you can use it as key on options
//
//	typeOpts : you can create your own logic how to copy value between 2 fields with different type (ex. string to time.Time) by type name of destination
//	 - as example if package type is "time.Time" you can use "Time" as key on options
func CopyStruct(src, dst any, fieldOpts, typeOpts map[string](func(srcValue reflect.Value) (dstValue reflect.Value, err error))) error {
	typeOfSrc := reflect.TypeOf(src)
	valueOfSrc := reflect.ValueOf(src)

	if typeOfSrc.Kind() != reflect.Struct {
		return errors.NewWithCode(codes.CodeInvalidValue, "src: parameter must be a struct but given type is "+typeOfSrc.Kind().String())
	}

	typeOfDst := reflect.TypeOf(dst)
	valueOfDst := reflect.ValueOf(dst)

	if typeOfDst.Kind() != reflect.Pointer {
		return errors.NewWithCode(codes.CodeInvalidValue, "dst: parameter must be a pointer but given type is "+typeOfDst.Kind().String())
	}

	if typeOfDst.Elem().Kind() != reflect.Struct {
		return errors.NewWithCode(codes.CodeInvalidValue, "dst: parameter must be a pointer but given type is pointer of "+typeOfDst.Elem().Kind().String())
	}

	for i := 0; i < typeOfSrc.NumField(); i++ {
		var srcFieldValue reflect.Value
		if valueOfSrc.Field(i).Kind() == reflect.Pointer {
			srcFieldValue = valueOfSrc.Field(i).Elem()
		} else {
			srcFieldValue = valueOfSrc.Field(i)
		}

		var dstFieldValue reflect.Value = valueOfDst.Elem().FieldByName(typeOfSrc.Field(i).Name)
		if dstFieldValue.Kind() == reflect.Pointer {
			dstFieldValue = dstFieldValue.Elem()
		}

		if !dstFieldValue.IsValid() {
			continue
		}

		fieldDst, isExists := typeOfDst.Elem().FieldByName(typeOfSrc.Field(i).Name)
		if !isExists {
			continue
		}

		if opts, isExists := fieldOpts[fieldDst.Name]; isExists {
			value, err := opts(srcFieldValue)
			if err != nil {
				return errors.NewWithCode(codes.CodeInvalidValue, err.Error())
			}
			dstFieldValue.Set(value)
			delete(fieldOpts, fieldDst.Name)
			continue
		}

		if opts, isExists := typeOpts[fieldDst.Type.Name()]; isExists {
			value, err := opts(srcFieldValue)
			if err != nil {
				return errors.NewWithCode(codes.CodeInvalidValue, err.Error())
			}
			dstFieldValue.Set(value)
			continue
		}

		switch {
		case srcFieldValue.CanInt():
			{
				dstFieldValue.SetInt(srcFieldValue.Int())
			}

		case srcFieldValue.CanFloat():
			{
				dstFieldValue.SetFloat(srcFieldValue.Float())
			}

		case srcFieldValue.CanUint():
			{
				dstFieldValue.SetUint(srcFieldValue.Uint())
			}

		case srcFieldValue.Kind() == reflect.String, srcFieldValue.Kind() == reflect.Bool:
			{
				dstFieldValue.Set(srcFieldValue)
			}

		case srcFieldValue.Kind() == reflect.Struct:
			{
				// TODO: set value on nested struct
				// TODO: set value on destination (dst) struct on pointer field
			}
		}
	}

	return nil
}

func ToPtr[T any](value T) *T {
	return &value
}

// Convert any value to `T` data type and if value is not valid or `nil`, will use default value of `T` (null/nil safety) as return value
func ToSafeValue[T any](value any) T {
	rv := reflect.ValueOf(value)
	for rv.Kind() == reflect.Pointer {
		if !rv.IsNil() {
			rv = rv.Elem()
		} else {
			break
		}
	}

	if !rv.IsValid() {
		return ToSafeValue[T](new(T))
	}

	safeValue, isOk := rv.Interface().(T)
	if !isOk {
		return ToSafeValue[T](new(T))
	}
	return safeValue
}

// TODO: convert graphql response to struct

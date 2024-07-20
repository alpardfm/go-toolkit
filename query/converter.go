package query

import (
	"fmt"
	"reflect"
	"time"

	"github.com/alpardfm/go-toolkit/sql"
)

type builderFunction func(primitiveType int8, isLike, isMany bool, fieldName, paramTag, dbTag string, args interface{})

const (
	Int int8 = iota
	IntArr
	Int64
	Int64Arr
	Int32
	Int32Arr
	Int16
	Int16Arr
	Int8
	Int8Arr
	Uint
	UintArr
	Uint64
	Uint64Arr
	Uint32
	Uint32Arr
	Uint16
	Uint16Arr
	Uint8
	Uint8Arr
	String
	StringArr
	Float32
	Float32Arr
	Float64
	Float64Arr
	Bool
	Time
	TimeArr
)

func traverseOnParam(paramTagName, dbTagName, fieldTagName, fieldName, paramTagValue, dbTagValue string, aliasMap map[string]string, p reflect.Value, builderFunc builderFunction) {
	switch p.Kind() {

	// on pointer/ interface
	case reflect.Ptr, reflect.Interface:

		if !p.Elem().IsValid() || p.IsNil() {
			return
		}

		// handle if is not time type, null type, and struct
		// continue to traverse
		if !isTimeType(p.Elem()) && !isNullType(p.Elem()) && p.Elem().Kind() == reflect.Struct {
			traverseOnParam(paramTagName, dbTagName, fieldTagName, fieldName+"."+p.Elem().Type().Name(), paramTagValue, dbTagValue, aliasMap, p.Elem(), builderFunc)
		}

		// else convert on types
		convertOnTypes(paramTagValue, dbTagValue, fieldName, p.Elem(), builderFunc)
		return

	// on struct
	case reflect.Struct:
		if isTimeType(p) {
			convertOnTypes(paramTagValue, dbTagValue, fieldName+"."+getNameFromStructTagOrOriginalName(fieldTagName, p, 0), p, builderFunc)
			return
		}

		for i := 0; i < p.NumField(); i++ {
			// only exported struct that can be Traversed
			if p.Field(i).CanSet() {
				paramTagValue = p.Type().Field(i).Tag.Get(paramTagName)
				dbTagValue = p.Type().Field(i).Tag.Get(dbTagName)

				if dbTagValue == "-" {
					continue
				}
				var address string
				if p.CanAddr() {
					address = fmt.Sprint(p.Addr().Pointer())
				}
				alias := aliasMap[address]
				if alias != "" && address != "" {
					dbTagValue = alias + "." + dbTagValue
				}

				if isNullType(p.Field(i)) {
					convertOnTypes(paramTagValue, dbTagValue, fieldName+"."+getNameFromStructTagOrOriginalName(fieldTagName, p, i), p.Field(i), builderFunc)
					continue
				}

				traverseOnParam(paramTagName, dbTagName, fieldTagName, fieldName+"."+getNameFromStructTagOrOriginalName(fieldTagName, p, i), paramTagValue, dbTagValue, aliasMap, p.Field(i), builderFunc)
			}
		}

	default:
		convertOnTypes(paramTagValue, dbTagValue, fieldName, p, builderFunc)
		return
	}
}

func convertOnTypes(paramTagValue, dbTagValue, fieldName string, e reflect.Value, builderFunc builderFunction) {
	var (
		args           interface{}
		isMany, isLike bool
		primitiveType  int8
	)

	switch f := e.Interface().(type) {
	// Integer Fields
	case []int64,
		[]*int64,
		[]uint64,
		[]*uint64,
		[]sql.NullInt64,
		[]*sql.NullInt64,
		[]int32,
		[]*int32,
		[]uint32,
		[]*uint32,
		[]int16,
		[]*int16,
		[]uint16,
		[]*uint16,
		[]int8,
		[]*int8,
		[]uint8,
		[]*uint8,
		[]int,
		[]*int,
		[]uint,
		[]*uint,
		int64,
		uint64,
		sql.NullInt64,
		int32,
		uint32,
		int16,
		uint16,
		int8,
		uint8,
		int,
		uint:
		primitiveType, isMany, args = convertIntArgs(f)
		builderFunc(primitiveType, isLike, isMany, fieldName, paramTagValue, dbTagValue, args)
		return
	// Float fields
	case []float64,
		[]*float64,
		[]float32,
		[]*float32,
		[]sql.NullFloat64,
		[]*sql.NullFloat64,
		float64,
		float32,
		sql.NullFloat64:
		primitiveType, isMany, args = convertFloatArgs(f)
		builderFunc(primitiveType, isLike, isMany, fieldName, paramTagValue, dbTagValue, args)
		return
	// String fields
	case []string,
		[]*string,
		[]sql.NullString,
		[]*sql.NullString,
		string,
		sql.NullString:
		primitiveType, isMany, isLike, args = convertStringArgs(f)
		builderFunc(primitiveType, isLike, isMany, fieldName, paramTagValue, dbTagValue, args)
		return
	// Bool
	case []bool,
		[]*bool,
		[]sql.NullBool,
		[]*sql.NullBool,
		bool,
		sql.NullBool:
		primitiveType, isMany, args = convertBoolArgs(f)
		builderFunc(primitiveType, isLike, isMany, fieldName, paramTagValue, dbTagValue, args)
		return
	// time
	case []time.Time,
		[]*time.Time,
		[]sql.NullTime,
		[]*sql.NullTime,
		[]sql.NullDate,
		[]*sql.NullDate,
		time.Time,
		sql.NullTime,
		sql.NullDate:
		primitiveType, isMany, args = convertTimeArgs(f)
		builderFunc(primitiveType, isLike, isMany, fieldName, paramTagValue, dbTagValue, args)
		return
	}
}

func isTimeType(e reflect.Value) bool {
	return e.Kind() == reflect.Struct && e.Type().Name() == "Time"
}

func isNullType(e reflect.Value) bool {
	return e.Kind() == reflect.Struct &&
		(e.Type().Name() == "NullString" ||
			e.Type().Name() == "NullBool" ||
			e.Type().Name() == "NullFloat64" ||
			e.Type().Name() == "NullInt64" ||
			e.Type().Name() == "NullTime" ||
			e.Type().Name() == "NullDate")
}

func getNameFromStructTagOrOriginalName(fieldName string, v reflect.Value, i int) string {
	name := v.Type().Field(i).Tag.Get(fieldName)
	if len(name) > 0 {
		return name
	}
	return v.Type().Field(i).Name
}

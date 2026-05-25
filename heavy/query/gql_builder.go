package query

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"

	"github.com/alpardfm/go-toolkit/codes"
	"github.com/alpardfm/go-toolkit/errors"
)

type GQLParamsInterface interface {
	// Use this to build graphql parameters
	Build(params any) (string, error)

	// Use this if you want to create your own logic while build graphql parameters.
	//   - Define the options by type of field so you can define 1 options for all field with same type. As example the field type is "time.Time" but you want to format it to string in graphql parameters. So, use "time.Time" as key on options
	AddTypeOptions(options map[string]func(v reflect.Value) (reflect.Value, error))

	// Use this if you want to create your own logic while build graphql parameters.
	//   - Define the options by name from `gql` tag and you can define your own tag name when initialize this builder with `NewGQLParamsBuilder`. As example the field type is "time.Time" but you want to format it to string in graphql parameters. So, use "time.Time" as key on options
	AddFieldOptions(options map[string]func(v reflect.Value) (reflect.Value, error))
}

type gqlParamsBuilder struct {
	params       bytes.Buffer
	paramTag     string
	typeOptions  map[string]func(v reflect.Value) (reflect.Value, error)
	fieldOptions map[string]func(v reflect.Value) (reflect.Value, error)
}

func NewGQLParamsBuilder(paramTag string, options ...func(*gqlParamsBuilder) error) (GQLParamsInterface, error) {
	gql := gqlParamsBuilder{
		paramTag: paramTag,
		params:   *bytes.NewBuffer([]byte{}),
	}

	for _, opt := range options {
		if err := opt(&gql); err != nil {
			return nil, err
		}
	}

	return &gql, nil
}

func (g *gqlParamsBuilder) Build(params any) (string, error) {
	reflectValue := reflect.ValueOf(params)
	reflectType := reflect.TypeOf(params)

	if reflectValue.Kind() == reflect.Pointer {
		return "", errors.NewWithCode(codes.CodeInvalidValue, "params: cannot be a pointer")
	}

	if reflectValue.Kind() != reflect.Struct {
		return "", errors.NewWithCode(codes.CodeInvalidValue, "params: must be a struct")
	}

	for i := 0; i < reflectType.NumField(); i++ {
		fieldTagName, isExists := reflectType.Field(i).Tag.Lookup(g.paramTag)
		if !isExists {
			continue
		}

		fieldType := reflectType.Field(i)
		fieldValue := reflectValue.Field(i)

		if fieldValue.IsZero() || !fieldValue.IsValid() {
			continue
		}
		// run field options
		if opt, isExists := g.fieldOptions[fieldTagName]; isExists {
			result, err := opt(fieldValue)
			if err != nil {
				return "", err
			}
			fieldValue = result
			delete(g.fieldOptions, fieldTagName)
		}

		// run type options
		if opt, isExists := g.typeOptions[fmt.Sprintf("%v", fieldType.Type)]; isExists {
			result, err := opt(fieldValue)
			if err != nil {
				return "", err
			}
			fieldValue = result
		}

		switch {
		// int
		case fieldValue.Kind() == reflect.Int, fieldValue.CanInt():
			g.params.WriteString(fmt.Sprintf(` %s: %d `, fieldTagName, fieldValue.Int()))

		// string
		case fieldValue.Kind() == reflect.String:
			g.params.WriteString(fmt.Sprintf(` %s: "%v" `, fieldTagName, fieldValue.String()))

		// float
		case fieldValue.CanFloat():
			g.params.WriteString(fmt.Sprintf(` %s: %v `, fieldTagName, fieldValue.Float()))

		// slice/array
		case fieldValue.Kind() == reflect.Slice, fieldValue.Kind() == reflect.Array:
			sType := fmt.Sprintf("%v", fieldValue.Type())
			var paramFormat *bytes.Buffer = bytes.NewBuffer([]byte{})
			switch {
			case strings.Contains(sType, "]int"), strings.Contains(sType, "]float"):
				paramFormat.WriteString(` %s: %v `)
			case strings.Contains(sType, "]string"):
				paramFormat.WriteString(` %s: %q `)
			}

			if fmt.Sprintf("%v", fieldValue.Interface()) == "[]" {
				continue
			}

			g.params.WriteString(fmt.Sprintf(paramFormat.String(), fieldTagName, fieldValue.Interface()))

		// boolean
		case fieldValue.Kind() == reflect.Bool:
			if fieldValue.Bool() {
				g.params.WriteString(fmt.Sprintf(" %s: %v ", fieldTagName, fieldValue.Bool()))
			}

		}

	}

	return g.params.String(), nil
}

func (g *gqlParamsBuilder) AddTypeOptions(options map[string]func(v reflect.Value) (reflect.Value, error)) {
	g.typeOptions = options
}

func (g *gqlParamsBuilder) AddFieldOptions(options map[string]func(v reflect.Value) (reflect.Value, error)) {
	g.fieldOptions = options
}

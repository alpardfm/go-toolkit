package parser

import (
	"github.com/alpardfm/go-toolkit/codes"
	"github.com/alpardfm/go-toolkit/errors"
	jsoniter "github.com/json-iterator/go"
)

type jsonConfig string

const (
	vanillaCompatible jsonConfig = "standard"
	defaultConfig     jsonConfig = "default"
	fastestConfig     jsonConfig = "fastest"
	customConfig      jsonConfig = "custom"
)

// JSONOptions configures the JSON parser behavior including encoding settings.
type JSONOptions struct {
	Config                        jsonConfig
	IndentionStep                 int
	MarshalFloatWith6Digits       bool
	EscapeHTML                    bool
	SortMapKeys                   bool
	UseNumber                     bool
	DisallowUnknownFields         bool
	TagKey                        string
	OnlyTaggedField               bool
	ValidateJSONRawMessage        bool
	ObjectFieldMustBeSimpleString bool
	CaseSensitive                 bool
}

// JSONInterface defines the contract for JSON marshaling and unmarshaling operations.
type JSONInterface interface {
	Marshal(orig any) ([]byte, error)
	Unmarshal(blob []byte, dest any) error
	// TODO: add schema validation
}

type jsonParser struct {
	API jsoniter.API
}

func initJSON(opt JSONOptions) JSONInterface {
	var jsonAPI jsoniter.API
	switch opt.Config {
	case defaultConfig:
		jsonAPI = jsoniter.ConfigDefault
	case fastestConfig:
		jsonAPI = jsoniter.ConfigFastest
	case customConfig:
		jsonAPI = jsoniter.Config{
			IndentionStep:                 opt.IndentionStep,
			MarshalFloatWith6Digits:       opt.MarshalFloatWith6Digits,
			EscapeHTML:                    opt.EscapeHTML,
			SortMapKeys:                   opt.SortMapKeys,
			UseNumber:                     opt.UseNumber,
			DisallowUnknownFields:         opt.DisallowUnknownFields,
			TagKey:                        opt.TagKey,
			OnlyTaggedField:               opt.OnlyTaggedField,
			ValidateJsonRawMessage:        opt.ValidateJSONRawMessage,
			ObjectFieldMustBeSimpleString: opt.ObjectFieldMustBeSimpleString,
			CaseSensitive:                 opt.CaseSensitive,
		}.Froze()
	default:
		jsonAPI = jsoniter.ConfigCompatibleWithStandardLibrary
	}

	p := &jsonParser{
		API: jsonAPI,
	}

	return p
}

func (p *jsonParser) Marshal(orig any) ([]byte, error) {
	stream := p.API.BorrowStream(nil)
	defer p.API.ReturnStream(stream)
	stream.WriteVal(orig)
	result := make([]byte, stream.Buffered())
	if stream.Error != nil {
		return nil, errors.NewWithCode(codes.CodeJSONMarshalError, stream.Error.Error())
	}
	copy(result, stream.Buffer())
	return result, nil
}

func (p *jsonParser) Unmarshal(blob []byte, dest any) error {
	iter := p.API.BorrowIterator(blob)
	defer p.API.ReturnIterator(iter)
	iter.ReadVal(dest)
	if iter.Error != nil {
		return errors.NewWithCode(codes.CodeJSONUnmarshalError, iter.Error.Error())
	}
	return nil
}

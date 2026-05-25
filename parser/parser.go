package parser

// Parser provides access to different parsing implementations.
type Parser interface {
	JSONParser() JSONInterface
}

// Options configures the parser with specific settings for each format.
type Options struct {
	JSONOptions JSONOptions
}

type parser struct {
	json JSONInterface
}

// InitParser creates a new Parser instance configured with the given options.
func InitParser(opt Options) Parser {
	return &parser{
		json: initJSON(opt.JSONOptions),
	}
}

func (p *parser) JSONParser() JSONInterface {
	return p.json
}

package parser

import "github.com/alpardfm/go-toolkit/log"

type Parser interface {
	JSONParser() JSONInterface
}

type Options struct {
	JSONOptions JSONOptions
}

type parser struct {
	json JSONInterface
}

func InitParser(log log.Interface, opt Options) Parser {
	return &parser{
		json: initJSON(opt.JSONOptions),
	}
}

func (p *parser) JSONParser() JSONInterface {
	return p.json
}

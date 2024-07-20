package format

import (
	"bytes"
	"text/template"

	"github.com/alpardfm/go-toolkit/codes"
	"github.com/alpardfm/go-toolkit/errors"
)

// Format string with golang template
func StringParseTmpl(strFmt string, values any) (string, error) {
	tmpl, err := template.New("").Parse(strFmt)
	if err != nil {
		return "", errors.NewWithCode(codes.CodeStrTemplateInvalidFormat, "cannot parse str format, %v", err)
	}

	buff := bytes.Buffer{}
	if err := tmpl.Execute(&buff, values); err != nil {
		return "", errors.NewWithCode(codes.CodeStrTemplateExecuteErr, "cannot execute template, %v", err)
	}

	return buff.String(), nil
}

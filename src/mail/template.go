package mail

import (
	"bytes"
	"text/template"

	"go.uber.org/zap"
)

// parsing email template function
func EmailTemplate(tmplPath string, data map[string]interface{}) string {
	// parsing template file
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		zap.S().Error(err)
		return ""
	}
	// creating new buffer as io writer
	buf := new(bytes.Buffer)
	// pasing above template with data and result data in buffer
	err = tmpl.Execute(buf, data)
	if err != nil {
		zap.S().Error(err)
		return ""
	}
	// return buffer in string
	return buf.String()
}

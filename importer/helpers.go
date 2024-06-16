package importer

import (
    "text/template"
    "strings"
    "fmt"
)

var (
	valTrue = true
	tmplFuncMap = template.FuncMap{
		"printStringSlice": func(slice []string) string {return fmt.Sprintf("[%s]",fmt.Sprintf("%s%s%s", "\"", strings.Join(slice, "\", \""), "\""))},
	}
)
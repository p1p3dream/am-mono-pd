package template

import (
	"bytes"
	"fmt"
	html_template "html/template"
	"io"
	"path/filepath"
	text_template "text/template"

	"github.com/Masterminds/sprig/v3"
)

func baseFuncMap() map[string]any {
	m := sprig.GenericFuncMap()

	m["osAbs"] = filepath.Abs

	return m
}

func HtmlFuncMap() html_template.FuncMap {
	return html_template.FuncMap(baseFuncMap())
}

func TextFuncMap() text_template.FuncMap {
	return text_template.FuncMap(baseFuncMap())
}

func ParseHtml(name string, body string) (io.Reader, error) {
	t, err := html_template.New(name).Funcs(HtmlFuncMap()).Parse(body)
	if err != nil {
		return nil, fmt.Errorf("couldn't parse config file: %w", err)
	}

	buf := new(bytes.Buffer)

	if err := t.Execute(buf, nil); err != nil {
		return nil, fmt.Errorf("couldn't execute config file template: %w", err)
	}

	return buf, nil
}

func ParseText(name string, body string) (io.Reader, error) {
	t, err := text_template.New(name).Funcs(TextFuncMap()).Parse(body)
	if err != nil {
		return nil, fmt.Errorf("couldn't parse config file: %w", err)
	}

	buf := new(bytes.Buffer)

	if err := t.Execute(buf, nil); err != nil {
		return nil, fmt.Errorf("couldn't execute config file template: %w", err)
	}

	return buf, nil
}

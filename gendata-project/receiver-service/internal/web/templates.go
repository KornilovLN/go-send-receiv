package web

import (
	"embed"
	"html/template"
	"strings"
)

//go:embed templates/*.html
var templateFS embed.FS

var templates *template.Template

func init() {
	var err error
	templates, err = template.ParseFS(templateFS, "templates/*.html")
	if err != nil {
		panic(err)
	}
}

// GetIndexTemplate возвращает HTML шаблон главной страницы
func GetIndexTemplate() string {
	var buf strings.Builder
	err := templates.ExecuteTemplate(&buf, "index.html", nil)
	if err != nil {
		return "Error loading template"
	}
	return buf.String()
}

// GetTableTemplate возвращает HTML шаблон для табличного отображения данных
func GetTableTemplate() string {
	var buf strings.Builder
	err := templates.ExecuteTemplate(&buf, "table.html", nil) // создать table.html
	if err != nil {
		return "Error loading template"
	}
	return buf.String()
}

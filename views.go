package main 

import (
	"html/template"
	"net/http"
	"path/filepath"
)

func index(writer http.ResponseWriter, request *http.Request) {
	tmpl , _ := template.ParseFiles(filepath.Join("staticfiles","index.html"))
	tmpl.Execute(writer,nil)
}
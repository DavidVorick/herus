package main

// index.go formats and serves the index page for the herus website. It also
// contains all of the code that manages the headers and footers for the
// website.

import (
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
)

const (
	indexPage = "/index.go"

	indexTitle = "Herus - Learning Evolved"

	cssPrefix = "/css/"
)

var (
	footerTpl = filepath.Join(dirTemplates, "footer.tpl")
	headerTpl = filepath.Join(dirTemplates, "header.tpl")
	indexTpl  = filepath.Join(dirTemplates, "index.tpl")
)

// HeaderTemplateData defines the data which is used to fill out the header
// template file.
type HeaderTemplateData struct {
	Title      string
	CSSImports []string // example: 'css/pure-min.css'
}

// IndexTemplateData defines the data which is used to fill out the index
// template file.
type IndexTemplateData struct {
	// Empty right now.
}

// cssHandler feeds all of the css files to anyone requesting them.
func cssHandler(w http.ResponseWriter, r *http.Request) {
	cssFile := strings.TrimPrefix(r.URL.Path, cssPrefix)
	css, err := ioutil.ReadFile(filepath.Join(dirCSS, cssFile))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_, err = w.Write(css)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

// executeHeader builds the header portion of the page from the header
// template.
func executeHeader(w io.Writer, htd HeaderTemplateData) error {
	t, err := template.ParseFiles(headerTpl)
	if err != nil {
		return err
	}
	return t.Execute(w, htd)
}

// executeFooter builds the footer portion of the page.
func executeFooter(w io.Writer) error {
	t, err := template.ParseFiles(footerTpl)
	if err != nil {
		return err
	}
	return t.Execute(w, nil)
}

// executeIndex builds the body portion of the index page.
func executeIndexBody(w io.Writer, itd IndexTemplateData) error {
	t, err := template.ParseFiles(indexTpl)
	if err != nil {
		return err
	}
	return t.Execute(w, itd)
}

// indexHandler will handle any requests coming to the index page.
func indexHandler(w http.ResponseWriter, r *http.Request) {
	err := executeHeader(w, HeaderTemplateData{Title: indexTitle})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = executeIndexBody(w, IndexTemplateData{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = executeFooter(w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

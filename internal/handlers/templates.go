package handlers

import (
	"github.com/neophenix/lxdepot/internal/utils"
	"html/template"
	"log"
)

// some settings from our main func
var WebRoot string
var CacheTemplates bool

// template cache
var templates = make(map[string]*template.Template)

// readTemplate is used by the various handlers to read the template file off disk, or return
// the template from cache if we already did that.  -cache_templates=false can be passed on the
// command line to always read off disk, useful for developing
func readTemplate(filename string) *template.Template {
	if CacheTemplates {
		if tmpl, ok := templates[filename]; ok {
			return tmpl
		}
	}

	// Until I find this is bad, I'm just going to always pass these functions into the template to simplify code.
	funcs := template.FuncMap{
		"MakeBytesMoreHuman":    utils.MakeBytesMoreHuman,
		"MakeIntBytesMoreHuman": utils.MakeIntBytesMoreHuman,
	}

	// web templates always have the base.tmpl that provides the overall layout, and then the requested template
	// provides all the content
	t, err := template.New(filename).Funcs(funcs).ParseFiles(WebRoot+"/templates/base.tmpl", WebRoot+"/templates/"+filename)
	if err != nil {
		log.Fatal("Could not open template: " + WebRoot + "/" + filename + " : " + err.Error())
	}

	// drop the template in cache for later
	templates[filename] = t

	return t
}

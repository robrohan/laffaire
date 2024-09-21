package handlers

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/robrohan/go-web-template/internals/models"
)

type pageData struct {
	Title       string
	CompanyName string
	User        *models.User
}

func TemplateInit() *template.Template {
	t, err := template.ParseGlob("./templates/*")
	if err != nil {
		log.Println("Cannot parse templates: ", err)
		os.Exit(-1)
	}

	return t
}

func ServePage(env *models.Env, t *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		routeMatch, _ := regexp.Compile(`\/(\w+)`)
		pd := pageData{
			"Laffaire Home",
			"Laffaire",
			env.User,
		}

		matches := routeMatch.FindStringSubmatch(r.URL.Path)

		env.Log.Debug("request", "path", r.URL.Path)
		env.Log.Debug("request", "match", matches)

		if len(matches) >= 1 {
			page := matches[1] + ".html"
			if t.Lookup(page) != nil {
				w.WriteHeader(200)
				t.ExecuteTemplate(w, page, pd)
				return
			}
		} else if r.URL.Path == "/" {
			w.WriteHeader(200)
			t.ExecuteTemplate(w, "index.html", pd)
			return
		}

		w.WriteHeader(404)
		w.Write([]byte("Not Found"))
	}
}

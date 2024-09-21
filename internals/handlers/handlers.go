package handlers

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/google/uuid"
	"github.com/robrohan/go-web-template/internals/env"
	"github.com/robrohan/go-web-template/internals/models"
)

type pageData struct {
	Title       string
	CompanyName string
	User        *models.User
}

type eventListPageData struct {
	pageData
	Events *[]models.Event
}

type eventPageData struct {
	pageData
	Event *models.Event
}

func TemplateInit() *template.Template {
	t, err := template.ParseGlob("./templates/*")
	if err != nil {
		log.Println("Cannot parse templates: ", err)
		os.Exit(-1)
	}

	return t
}

func EventsPage(env *env.Env, t *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		page := "events.html"

		userId, _ := uuid.Parse(env.User.UUID)
		events, err := env.Repo.GetEventsByUserId(userId)
		if err != nil {
			log.Println("events query errored ", err)
			return
		}

		pd := eventListPageData{
			pageData{
				"Laffaire Home",
				"Laffaire",
				env.User,
			},
			events,
		}
		if t.Lookup(page) != nil {
			t.ExecuteTemplate(w, page, pd)
		}
	}
}

func EventPage(env *env.Env, t *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		page := "event.html"

		pd := eventPageData{
			pageData{
				"Laffaire Home",
				"Laffaire",
				env.User,
			},
			nil,
		}

		switch r.Method {
		case "POST":
			r.ParseForm()
			eventUuid := r.FormValue("event_uuid")
			log.Println("event uuid from form is: ", eventUuid)

			if eventUuid != "" {
				event := models.Event{
					UUID:        eventUuid,
					UserId:      env.User.UUID,
					Title:       r.FormValue("title"),
					Description: r.FormValue("description"),
				}
				env.Repo.UpsertEvent(&event)

				http.Redirect(w, r, "/-/events", http.StatusFound)
			}
		case "GET":
			eventUuid := r.URL.Query().Get("event")
			if eventUuid != "" {
				eventId, err := uuid.Parse(eventUuid)
				if err != nil {
					log.Println("Cannot parse event uuid ", err)
					return
				}
				event, err := env.Repo.GetEventById(eventId)
				if err != nil {
					log.Println("Cannot get the event from the db ", err)
					return
				}
				pd.Event = event
			} else {
				event := models.Event{
					UUID: uuid.New().String(),
				}
				pd.Event = &event
			}
		}

		if t.Lookup(page) != nil {
			t.ExecuteTemplate(w, page, pd)
		}
	}
}

func ServePage(env *env.Env, t *template.Template) http.HandlerFunc {
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

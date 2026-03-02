package handlers

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/robrohan/laffaire/internals/env"
	"github.com/robrohan/laffaire/internals/ical"
	"github.com/robrohan/laffaire/internals/models"
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

type entriesListPageData struct {
	pageData
	EventUUID *uuid.UUID
	Entries   *[]models.Entry
}

type entryPageData struct {
	pageData
	Entry *models.Entry
}

func TemplateInit() *template.Template {
	t, err := template.ParseGlob("./templates/*")
	if err != nil {
		log.Println("Cannot parse templates: ", err)
		os.Exit(-1)
	}

	return t
}

func createDateTime(date string, time string) string {
	if date == "" {
		return ""
	}
	date = strings.Replace(date, "-", "", -1)
	if time != "" {
		time = strings.Replace(time, ":", "", -1)
		date = date + "T" + time + "00"
	} else {
		date = date + "T000000"
	}
	return date
}

func IcalPage(env *env.Env, t *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		log.Println("vars", vars["id"])

		eventUuid, _ := uuid.Parse(vars["id"])
		entries, err := env.Repo.GetEntriesByEventId(eventUuid)
		if err != nil {
			log.Println("cannot get the event from the db ", err)
			return
		}

		log.Println("entries", entries)

		calendarName := "Plan"
		var ics bytes.Buffer
		log.Printf("creating prolog")
		ical.Prolog(&ics, calendarName, "//Rob Rohan//Made up go code//EN", "NZDT")
		for i := 0; i < len(*entries); i++ {
			e := (*entries)[i]

			start := createDateTime(e.StartDate, e.StartTime)
			end := createDateTime(e.EndDate, e.EndTime)
			calid := strings.Split(e.UUID, "-")[0]
			timestamp := int32(time.Now().Unix())

			if end == "" {
				end = start
			}

			if start != "" {
				ics.WriteString("BEGIN:VEVENT\r\n")
				fmt.Fprintf(&ics, "DTSTAMP:%v\r\n", start)
				fmt.Fprintf(&ics, "UID:R-%v-%v\r\n", calid, timestamp)
				fmt.Fprintf(&ics, "DTSTART;VALUE=DATE:%v\r\n", start)
				fmt.Fprintf(&ics, "DTEND;VALUE=DATE:%v\r\n", end)
				fmt.Fprintf(&ics, "SUMMARY:%v\r\n", e.Subject)
				fmt.Fprintf(&ics, "DESCRIPTION:%v\r\n", e.Description)
				fmt.Fprintf(&ics, "CATEGORIES:%v\r\n", calendarName)
				ics.WriteString("END:VEVENT\r\n")
			}
		}
		log.Printf("writing epilog")
		ical.Epilog(&ics)

		w.Header().Set("Content-Type", "text/calendar")
		w.Write(ics.Bytes())
	}
}

func EntriesPage(env *env.Env, t *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		page := "entries.html"
		pd := entriesListPageData{
			pageData{
				"Laffaire Home",
				"Laffaire",
				env.User,
			},
			nil,
			nil,
		}

		eventId := r.URL.Query().Get("event")
		eventUuid, _ := uuid.Parse(eventId)
		entries, err := env.Repo.GetEntriesByEventId(eventUuid)
		if err != nil {
			log.Println("cannot get the event from the db ", err)
			return
		}

		pd.Entries = entries
		pd.EventUUID = &eventUuid

		if t.Lookup(page) != nil {
			t.ExecuteTemplate(w, page, pd)
		}
	}
}

func EntryPage(env *env.Env, t *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		page := "entry.html"
		pd := entryPageData{
			pageData{
				"Laffaire Home",
				"Laffaire",
				env.User,
			},
			nil,
		}

		// We should always have an event id
		eventUuid := r.FormValue("event_uuid")
		if eventUuid == "" {
			eventUuid = r.URL.Query().Get("event")
		}
		log.Println("have event uuid", eventUuid)

		switch r.Method {
		case "POST":
			r.ParseForm()
			entryUuid := r.FormValue("entry_uuid")

			log.Println("event uuid from form is:", entryUuid)

			allday := r.FormValue("all_day_event")
			private := r.FormValue("private")

			if entryUuid != "" {
				entry := models.Entry{
					UUID:        entryUuid,
					EventId:     eventUuid,
					Subject:     r.FormValue("subject"),
					StartDate:   r.FormValue("start_date"),
					StartTime:   r.FormValue("start_time"),
					EndDate:     r.FormValue("end_date"),
					EndTime:     r.FormValue("end_time"),
					AllDayEvent: (allday == "on"),
					Description: r.FormValue("description"),
					Location:    r.FormValue("location"),
					Private:     (private == "on"),
				}
				err := env.Repo.UpsertEntry(&entry)
				if err != nil {
					log.Println("upsert error", err)
					return
				}

				http.Redirect(w, r, "/-/entries?event="+eventUuid, http.StatusFound)
			}
		case "GET":
			entryUuid := r.URL.Query().Get("entry")
			if entryUuid != "" {
				entryId, err := uuid.Parse(entryUuid)
				if err != nil {
					log.Println("cannot parse event uuid ", err)
					return
				}
				entry, err := env.Repo.GetEntryById(entryId)
				if err != nil {
					log.Println("cannot get the event from the db ", err)
					return
				}
				pd.Entry = entry
			} else {
				entry := models.Entry{
					UUID:    uuid.New().String(),
					EventId: eventUuid,
				}
				pd.Entry = &entry
			}
		case "DELETE":
			log.Println("delete entry")
			entryUuid := r.URL.Query().Get("entry")
			eventUuid := r.URL.Query().Get("event")
			log.Println(entryUuid, eventUuid)

			err := env.Repo.DeleteEntry(entryUuid, eventUuid)
			if err != nil {
				log.Println("cannot delete the entry from the db ", err)
				return
			}
			http.Redirect(w, r, "/-/entries?event="+eventUuid, http.StatusTemporaryRedirect)
			return
		}

		if t.Lookup(page) != nil {
			t.ExecuteTemplate(w, page, pd)
		}
	}
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

type tokenListPageData struct {
	pageData
	Tokens *[]models.Token
}

type tokenPageData struct {
	pageData
	NewToken *models.Token // nil = show form; non-nil = show newly created token
}

func TokensPage(env *env.Env, t *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		page := "tokens.html"

		if r.Method == "DELETE" {
			tokenUuid := r.URL.Query().Get("token")
			if err := env.Repo.DeleteToken(tokenUuid, env.User.UUID); err != nil {
				log.Println("delete token error", err)
			}
			http.Redirect(w, r, "/-/tokens", http.StatusTemporaryRedirect)
			return
		}

		userId, _ := uuid.Parse(env.User.UUID)
		tokens, err := env.Repo.GetTokensByUserId(userId)
		if err != nil {
			log.Println("tokens query errored", err)
			return
		}

		pd := tokenListPageData{
			pageData{"Laffaire Tokens", "Laffaire", env.User},
			tokens,
		}
		if t.Lookup(page) != nil {
			t.ExecuteTemplate(w, page, pd)
		}
	}
}

func TokenPage(env *env.Env, t *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		page := "token.html"
		pd := tokenPageData{
			pageData{"Laffaire Token", "Laffaire", env.User},
			nil,
		}

		if r.Method == "POST" {
			r.ParseForm()
			name := r.FormValue("name")
			if name != "" {
				b := make([]byte, 32)
				if _, err := rand.Read(b); err != nil {
					log.Println("failed to generate token", err)
					return
				}
				token := models.Token{
					UUID:      uuid.New().String(),
					UserId:    env.User.UUID,
					Name:      name,
					Token:     fmt.Sprintf("%x", b),
					CreatedAt: time.Now().UTC().Format(time.RFC3339),
				}
				if err := env.Repo.CreateToken(&token); err != nil {
					log.Println("create token error", err)
					return
				}
				pd.NewToken = &token
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

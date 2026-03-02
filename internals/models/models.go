package models

import "github.com/google/uuid"

// UserInfo is the data we get back from the auth service
type UserInfo struct {
	Id            string `json:"id"`
	Email         string `json:"email"`
	Picture       string `json:"picture"`
	VerifiedEmail bool   `json:"verified_email"`
}

// User is an example model in the application (saved in the db)
type User struct {
	UUID    string  `db:"uuid"     json:"id"`
	Email   string  `db:"email"    json:"email"`
	Name    *string `db:"username" json:"name,omitempty"`
	Picture *string `db:"picture"  json:"picture,omitempty"`
	AuthId  string  `db:"authid"   json:"-"`
	Salt    *string `db:"salt"     json:"-"`
}

/////////////////////////////////////////////////////////

// Entry is a single item on a calendar
type Entry struct {
	UUID        string `db:"uuid"          json:"id"`
	EventId     string `db:"event_uuid"    json:"event_id"`
	Subject     string `db:"subject"       json:"subject"`
	StartDate   string `db:"start_date"    json:"start_date"`
	StartTime   string `db:"start_time"    json:"start_time"`
	EndDate     string `db:"end_date"      json:"end_date"`
	EndTime     string `db:"end_time"      json:"end_time"`
	AllDayEvent bool   `db:"all_day_event" json:"all_day_event"`
	Description string `db:"description"   json:"description"`
	Location    string `db:"location"      json:"location"`
	Private     bool   `db:"private"       json:"private"`
}

// Event is a group of Entries - for example Training, Doctors Appointment, etc
type Event struct {
	UUID        string `db:"uuid"        json:"id"`
	UserId      string `db:"user_uuid"   json:"-"`
	Title       string `db:"title"       json:"title"`
	Description string `db:"description" json:"description"`
}

//////////////////////////////////

func NewUser(authid string, email string, picture string) *User {
	id := uuid.New()
	a := User{
		UUID:    id.String(),
		AuthId:  authid,
		Email:   email,
		Picture: &picture,
	}
	return &a
}

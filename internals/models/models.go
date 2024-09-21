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
	UUID    string  `db:"uuid"`
	Email   string  `db:"email"`
	Name    *string `db:"username"`
	Picture *string `db:"picture"`
	AuthId  string  `db:"authid"`
	Salt    *string `db:"salt"`
}

/////////////////////////////////////////////////////////

// Entry is a single item on a calendar
type Entry struct {
	UUID        string `db:"uuid"`
	EventId     string `db:"event_uuid"`
	Subject     string `db:"subject"`
	StartDate   string `db:"start_date"`
	StartTime   string `db:"start_time"`
	EndDate     string `db:"end_date"`
	EndTime     string `db:"end_time"`
	AllDayEvent bool   `db:"all_day_event"`
	Description string `db:"description"`
	Location    string `db:"location"`
	Private     bool   `db:"private"`
}

// Event is a group of Entries - for example Training, Doctors Appointment, etc
type Event struct {
	UUID        string `db:"uuid"`
	UserId      string `db:"user_uuid"`
	Title       string `db:"title"`
	Description string `db:"description"`
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

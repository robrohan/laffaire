package repository

import (
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/robrohan/go-web-template/internals/models"
)

type DataRepository struct {
	Db                       *sqlx.DB
	upsertUserQuery          *sqlx.Stmt
	getUserByEmailQuery      *sqlx.Stmt
	getUserByIdQuery         *sqlx.Stmt
	upsertEventQuery         *sqlx.Stmt
	upsertEntryQuery         *sqlx.Stmt
	getEventsByUserIdQuery   *sqlx.Stmt
	getEventByIdQuery        *sqlx.Stmt
	getEntriesByEventIdQuery *sqlx.Stmt
}

func prepareQuery(query string, db *sqlx.DB) *sqlx.Stmt {
	stmt, err := db.Preparex(query)
	if err != nil {
		log.Fatal(err)
	}
	return stmt
}

// Attach creates a new repository and sets up prepared statements
func Attach(schema string, db *sqlx.DB, driver string) *DataRepository {
	a := DataRepository{
		Db: db,
	}

	a.upsertUserQuery = prepareQuery(`
		INSERT INTO users (uuid, authid, email, picture, salt)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (email) DO UPDATE
			SET picture = $4,
			salt = $5;
	`, db)

	a.getUserByEmailQuery = prepareQuery(`
		SELECT uuid, email, username, picture, authid, salt
		FROM users
		WHERE email = $1
	`, db)

	a.getUserByIdQuery = prepareQuery(`
		SELECT uuid, email, username, picture, authid, salt
		FROM users
		WHERE uuid = $1
	`, db)

	a.upsertEventQuery = prepareQuery(`
		INSERT INTO event (uuid, user_uuid, title, description)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (uuid) DO UPDATE
			SET title = $3,
			description = $4;
	`, db)

	a.upsertEntryQuery = prepareQuery(`
		INSERT INTO entry (
			uuid, event_uuid,
			subject,
			start_date, start_time,
			end_date, end_time,
			all_day_event, description, location, private)
		VALUES (
			$1, $2,
			$3,
			$4, $5,
			$6, $7,
			$8,	$9, $10, $11)
		ON CONFLICT (uuid) DO UPDATE
			SET subject = $3,
			start_date = $4,
			start_time = $5,
			end_date = $6,
			end_time = $7,
			all_day_event = $8,
			description = $9,
			location = $10,
			private = $11;
	`, db)

	a.getEventsByUserIdQuery = prepareQuery(`
		SELECT *
		FROM event
		WHERE user_uuid = $1
		LIMIT $2
		OFFSET $3
	`, db)

	a.getEntriesByEventIdQuery = prepareQuery(`
		SELECT *
		FROM entry
		WHERE event_uuid = $1
		LIMIT $2
		OFFSET $3
	`, db)

	a.getEventByIdQuery = prepareQuery(`
		SELECT *
		FROM event
		WHERE uuid = $1
	`, db)

	return &a
}

func (r *DataRepository) Begin() (*sqlx.Tx, error) {
	return r.Db.Beginx()
}

func (r *DataRepository) UpsertUser(user *models.User, salt string) error {
	_, err := r.upsertUserQuery.Exec(
		user.UUID, user.AuthId, user.Email, user.Picture, salt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *DataRepository) UpsertEvent(event *models.Event) error {
	_, err := r.upsertEventQuery.Exec(
		event.UUID, event.UserId, event.Title, event.Description,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *DataRepository) UpsertEntry(entry *models.Entry) error {
	_, err := r.upsertEntryQuery.Exec(
		entry.UUID, entry.EventId, entry.Subject,
		entry.StartDate, entry.StartTime,
		entry.EndDate, entry.EndTime,
		entry.AllDayEvent, entry.Description, entry.Private,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *DataRepository) GetUserById(uuid uuid.UUID) (*models.User, error) {
	rows, err := r.getUserByIdQuery.Queryx(uuid)
	if err != nil {
		return nil, err
	}
	if rows == nil {
		return nil, errors.New("no rows")
	}

	user := models.User{}
	for rows.Next() {
		err = rows.StructScan(&user)
		if err != nil {
			return nil, err
		}
	}

	return &user, nil
}

func (r *DataRepository) GetEventById(event_uuid uuid.UUID) (*models.Event, error) {
	rows, err := r.getEventByIdQuery.Queryx(event_uuid)
	if err != nil {
		return nil, err
	}
	if rows == nil {
		return nil, errors.New("no rows")
	}

	event := models.Event{}
	for rows.Next() {
		err = rows.StructScan(&event)
		if err != nil {
			return nil, err
		}
	}
	return &event, nil
}

func (r *DataRepository) GetEventsByUserId(user_uuid uuid.UUID) (*[]models.Event, error) {
	rows, err := r.getEventsByUserIdQuery.Queryx(user_uuid, 100, 0)
	if err != nil {
		return nil, err
	}
	if rows == nil {
		return nil, errors.New("no rows")
	}

	events := make([]models.Event, 0)
	for rows.Next() {
		event := models.Event{}
		err = rows.StructScan(&event)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return &events, nil
}

func (r *DataRepository) GetEntriesByEventId(event_uuid uuid.UUID) (*[]models.Entry, error) {
	rows, err := r.getEntriesByEventIdQuery.Queryx(event_uuid, 100, 0)
	if err != nil {
		return nil, err
	}
	if rows == nil {
		return nil, errors.New("no rows")
	}

	entries := make([]models.Entry, 0)
	for rows.Next() {
		entry := models.Entry{}
		err = rows.StructScan(&entry)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	return &entries, nil
}

func (r *DataRepository) GetUser(email string) (*models.User, error) {
	rows, err := r.getUserByEmailQuery.Queryx(email)
	if err != nil {
		return nil, err
	}
	if rows == nil {
		return nil, errors.New("no rows")
	}

	user := models.User{}
	for rows.Next() {
		err = rows.StructScan(&user)
		if err != nil {
			return nil, err
		}
	}

	return &user, nil
}

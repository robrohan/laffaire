package repository

import (
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/robrohan/laffaire/internals/models"
)

type DataRepository struct {
	Db                       *sqlx.DB
	upsertUserQuery          *sqlx.Stmt
	getUserByEmailQuery      *sqlx.Stmt
	getUserByIdQuery         *sqlx.Stmt
	upsertEventQuery         *sqlx.Stmt
	deleteEventQuery         *sqlx.Stmt
	upsertEntryQuery         *sqlx.Stmt
	deleteEntryQuery         *sqlx.Stmt
	getEventsByUserIdQuery   *sqlx.Stmt
	getEventByIdQuery        *sqlx.Stmt
	getEntriesByEventIdQuery *sqlx.Stmt
	getEntryByIdQuery        *sqlx.Stmt
	createTokenQuery         *sqlx.Stmt
	getTokensByUserIdQuery   *sqlx.Stmt
	deleteTokenQuery         *sqlx.Stmt
	getUserByTokenQuery      *sqlx.Stmt
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

	a.deleteEventQuery = prepareQuery(`
		DELETE FROM event
		WHERE uuid = $1
		AND user_uuid = $2
	`, db)

	a.deleteEntryQuery = prepareQuery(`
		DELETE FROM entry
		WHERE uuid = $1
		AND event_uuid = $2
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
		ORDER BY start_date
		LIMIT $2
		OFFSET $3
	`, db)

	a.getEventByIdQuery = prepareQuery(`
		SELECT *
		FROM event
		WHERE uuid = $1
	`, db)

	a.getEntryByIdQuery = prepareQuery(`
		SELECT *
		FROM entry
		WHERE uuid = $1
	`, db)

	a.createTokenQuery = prepareQuery(`
		INSERT INTO token (uuid, user_uuid, name, token, created_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (uuid) DO NOTHING
	`, db)

	a.getTokensByUserIdQuery = prepareQuery(`
		SELECT uuid, user_uuid, name, token, created_at
		FROM token
		WHERE user_uuid = $1
		ORDER BY created_at DESC
	`, db)

	a.deleteTokenQuery = prepareQuery(`
		DELETE FROM token
		WHERE uuid = $1
		AND user_uuid = $2
	`, db)

	a.getUserByTokenQuery = prepareQuery(`
		SELECT u.uuid, u.email, u.username, u.picture, u.authid, u.salt
		FROM users u
		JOIN token t ON u.uuid = t.user_uuid
		WHERE t.token = $1
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
		entry.AllDayEvent, entry.Description, entry.Location, entry.Private,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *DataRepository) DeleteEvent(eventUuid string, userUuid string) error {
	_, err := r.deleteEventQuery.Exec(eventUuid, userUuid)
	if err != nil {
		return err
	}
	return nil
}

func (r *DataRepository) DeleteEntry(entryUuid string, eventUuid string) error {
	_, err := r.deleteEntryQuery.Exec(entryUuid, eventUuid)
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

func (r *DataRepository) GetEntryById(entry_uuid uuid.UUID) (*models.Entry, error) {
	rows, err := r.getEntryByIdQuery.Queryx(entry_uuid)
	if err != nil {
		return nil, err
	}
	if rows == nil {
		return nil, errors.New("no rows")
	}

	entry := models.Entry{}
	for rows.Next() {
		err = rows.StructScan(&entry)
		if err != nil {
			return nil, err
		}
	}

	return &entry, nil
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

func (r *DataRepository) CreateToken(token *models.Token) error {
	_, err := r.createTokenQuery.Exec(
		token.UUID, token.UserId, token.Name, token.Token, token.CreatedAt,
	)
	return err
}

func (r *DataRepository) GetTokensByUserId(userUuid uuid.UUID) (*[]models.Token, error) {
	rows, err := r.getTokensByUserIdQuery.Queryx(userUuid)
	if err != nil {
		return nil, err
	}
	tokens := make([]models.Token, 0)
	for rows.Next() {
		t := models.Token{}
		if err = rows.StructScan(&t); err != nil {
			return nil, err
		}
		tokens = append(tokens, t)
	}
	return &tokens, nil
}

func (r *DataRepository) DeleteToken(tokenUuid string, userUuid string) error {
	_, err := r.deleteTokenQuery.Exec(tokenUuid, userUuid)
	return err
}

func (r *DataRepository) GetUserByToken(tokenValue string) (*models.User, error) {
	rows, err := r.getUserByTokenQuery.Queryx(tokenValue)
	if err != nil {
		return nil, err
	}
	user := models.User{}
	for rows.Next() {
		if err = rows.StructScan(&user); err != nil {
			return nil, err
		}
	}
	if user.UUID == "" {
		return nil, errors.New("token not found")
	}
	return &user, nil
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

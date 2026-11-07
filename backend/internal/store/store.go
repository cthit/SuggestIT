// Package store persists suggestions in PostgreSQL.
package store

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"time"

	"github.com/lib/pq"
)

// ErrNotFound is returned when a suggestion does not exist.
var ErrNotFound = errors.New("suggestion not found")

// AnonymousAuthor is used for suggestions without an author.
const AnonymousAuthor = "Anonymous"

// Suggestion is a stored suggestion.
type Suggestion struct {
	ID        string    `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	Title     string    `json:"title"`
	Text      string    `json:"text"`
	Author    string    `json:"author"`
}

// NewSuggestion is the validated input for creating a suggestion.
type NewSuggestion struct {
	Title  string
	Text   string
	Author string
}

// NewUUID returns a random (version 4, RFC 9562) UUID string.
func NewUUID() string {
	var b [16]byte
	_, _ = rand.Read(b[:]) // crypto/rand.Read never returns an error
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	var s [36]byte
	hex.Encode(s[0:8], b[0:4])
	s[8] = '-'
	hex.Encode(s[9:13], b[4:6])
	s[13] = '-'
	hex.Encode(s[14:18], b[6:8])
	s[18] = '-'
	hex.Encode(s[19:23], b[8:10])
	s[23] = '-'
	hex.Encode(s[24:], b[10:])
	return string(s[:])
}

// ValidUUID reports whether s is a UUID in canonical 8-4-4-4-12 hex form.
func ValidUUID(s string) bool {
	if len(s) != 36 {
		return false
	}
	for i := range len(s) {
		c := s[i]
		switch i {
		case 8, 13, 18, 23:
			if c != '-' {
				return false
			}
		default:
			if !('0' <= c && c <= '9' || 'a' <= c && c <= 'f' || 'A' <= c && c <= 'F') {
				return false
			}
		}
	}
	return true
}

// Postgres is the PostgreSQL-backed store.
type Postgres struct {
	db  *sql.DB
	now func() time.Time
}

// Open prepares a connection pool. It does not connect; use Ping.
func Open(dsn string) (*Postgres, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxIdleTime(5 * time.Minute)
	db.SetConnMaxLifetime(time.Hour)
	return &Postgres{db: db, now: time.Now}, nil
}

// Close closes the connection pool.
func (p *Postgres) Close() error { return p.db.Close() }

// Ping checks that the database is reachable.
func (p *Postgres) Ping(ctx context.Context) error { return p.db.PingContext(ctx) }

// Migrate creates the suggestion table if it does not exist. The schema is
// unchanged from earlier versions so existing data keeps working.
func (p *Postgres) Migrate(ctx context.Context) error {
	_, err := p.db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS suggestion (
			id        UUID      NOT NULL,
			timestamp TIMESTAMP NOT NULL,
			title     TEXT      NOT NULL,
			text      TEXT      NOT NULL,
			author    TEXT,
			PRIMARY KEY (id)
		)`)
	return err
}

const selectColumns = `SELECT id, timestamp, title, text, author FROM suggestion`

type scanner interface{ Scan(dest ...any) error }

func scanSuggestion(row scanner) (Suggestion, error) {
	var (
		s      Suggestion
		author sql.NullString
	)
	if err := row.Scan(&s.ID, &s.Timestamp, &s.Title, &s.Text, &author); err != nil {
		return Suggestion{}, err
	}
	// TIMESTAMP (without time zone) columns hold UTC wall-clock times.
	t := s.Timestamp
	s.Timestamp = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), time.UTC)
	s.Author = author.String
	if !author.Valid || author.String == "" {
		s.Author = AnonymousAuthor
	}
	return s, nil
}

// List returns all suggestions, newest first. It never returns nil.
func (p *Postgres) List(ctx context.Context) ([]Suggestion, error) {
	rows, err := p.db.QueryContext(ctx, selectColumns+` ORDER BY timestamp DESC, id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []Suggestion{}
	for rows.Next() {
		s, err := scanSuggestion(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

// Get returns one suggestion or ErrNotFound.
func (p *Postgres) Get(ctx context.Context, id string) (Suggestion, error) {
	if !ValidUUID(id) {
		return Suggestion{}, ErrNotFound
	}
	s, err := scanSuggestion(p.db.QueryRowContext(ctx, selectColumns+` WHERE id = $1`, id))
	if errors.Is(err, sql.ErrNoRows) {
		return Suggestion{}, ErrNotFound
	}
	return s, err
}

// Create inserts a suggestion with a new id and the current UTC time.
func (p *Postgres) Create(ctx context.Context, n NewSuggestion) (Suggestion, error) {
	s := Suggestion{
		ID:        NewUUID(),
		Timestamp: p.now().UTC().Truncate(time.Microsecond), // Postgres precision
		Title:     n.Title,
		Text:      n.Text,
		Author:    n.Author,
	}
	_, err := p.db.ExecContext(ctx,
		`INSERT INTO suggestion (id, timestamp, title, text, author) VALUES ($1, $2, $3, $4, $5)`,
		s.ID, s.Timestamp, s.Title, s.Text, s.Author)
	if err != nil {
		return Suggestion{}, err
	}
	return s, nil
}

// Delete removes one suggestion or returns ErrNotFound.
func (p *Postgres) Delete(ctx context.Context, id string) error {
	if !ValidUUID(id) {
		return ErrNotFound
	}
	res, err := p.db.ExecContext(ctx, `DELETE FROM suggestion WHERE id = $1`, id)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

// DeleteMany removes all suggestions with the given ids in one statement.
// Unknown ids are ignored. All ids must be valid UUIDs.
func (p *Postgres) DeleteMany(ctx context.Context, ids []string) error {
	if len(ids) == 0 {
		return nil
	}
	_, err := p.db.ExecContext(ctx, `DELETE FROM suggestion WHERE id = ANY($1::uuid[])`, pq.Array(ids))
	return err
}

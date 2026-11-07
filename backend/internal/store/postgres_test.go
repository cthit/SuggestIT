package store

import (
	"context"
	"os"
	"slices"
	"testing"
	"time"
)

// TestPostgres runs against a real database when SUGGESTIT_TEST_DATABASE_URL
// is set (e.g. postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable).
// It is skipped otherwise. It drops and recreates the suggestion table.
func TestPostgres(t *testing.T) {
	dsn := os.Getenv("SUGGESTIT_TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("SUGGESTIT_TEST_DATABASE_URL not set")
	}
	ctx := context.Background()
	p, err := Open(dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer p.Close()
	if err := p.Ping(ctx); err != nil {
		t.Fatal(err)
	}
	if _, err := p.db.ExecContext(ctx, `DROP TABLE IF EXISTS suggestion`); err != nil {
		t.Fatal(err)
	}
	for range 2 { // idempotent
		if err := p.Migrate(ctx); err != nil {
			t.Fatal(err)
		}
	}

	if list, err := p.List(ctx); err != nil || list == nil || len(list) != 0 {
		t.Fatalf("empty List = %v, %v", list, err)
	}

	// A legacy row with NULL author, written like the old backend did.
	legacyID := NewUUID()
	if _, err := p.db.ExecContext(ctx, `INSERT INTO suggestion VALUES ($1, '2020-01-02 03:04:05.123456', 'old', 'legacy', NULL)`, legacyID); err != nil {
		t.Fatal(err)
	}

	clock := time.Date(2026, 11, 6, 12, 0, 0, 123456789, time.FixedZone("CET", 3600))
	p.now = func() time.Time { clock = clock.Add(time.Second); return clock }
	a, err := p.Create(ctx, NewSuggestion{Title: "a", Text: "line1\nline2", Author: "Bob"})
	if err != nil {
		t.Fatal(err)
	}
	b, err := p.Create(ctx, NewSuggestion{Title: "b", Text: "x", Author: AnonymousAuthor})
	if err != nil {
		t.Fatal(err)
	}
	if a.Timestamp.Location() != time.UTC || a.Timestamp.Nanosecond()%1000 != 0 {
		t.Errorf("created timestamp %v not UTC/µs", a.Timestamp)
	}

	got, err := p.Get(ctx, a.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got != a {
		t.Errorf("Get = %+v, want %+v", got, a)
	}
	legacy, err := p.Get(ctx, legacyID)
	if err != nil {
		t.Fatal(err)
	}
	if legacy.Author != AnonymousAuthor || !legacy.Timestamp.Equal(time.Date(2020, 1, 2, 3, 4, 5, 123456000, time.UTC)) {
		t.Errorf("legacy row = %+v", legacy)
	}
	if _, err := p.Get(ctx, NewUUID()); err != ErrNotFound {
		t.Errorf("Get unknown = %v", err)
	}

	list, err := p.List(ctx)
	if err != nil {
		t.Fatal(err)
	}
	ids := func(l []Suggestion) []string {
		var out []string
		for _, s := range l {
			out = append(out, s.ID)
		}
		return out
	}
	if want := []string{b.ID, a.ID, legacyID}; !slices.Equal(ids(list), want) {
		t.Errorf("List order = %v, want %v", ids(list), want)
	}

	if err := p.Delete(ctx, b.ID); err != nil {
		t.Fatal(err)
	}
	if err := p.Delete(ctx, b.ID); err != ErrNotFound {
		t.Errorf("second Delete = %v", err)
	}
	if err := p.DeleteMany(ctx, []string{a.ID, legacyID, NewUUID()}); err != nil {
		t.Fatal(err)
	}
	if list, _ := p.List(ctx); len(list) != 0 {
		t.Errorf("after DeleteMany: %v", list)
	}
}

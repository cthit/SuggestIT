package store

import (
	"context"
	"encoding/hex"
	"strings"
	"testing"
)

func TestNewUUID(t *testing.T) {
	seen := map[string]bool{}
	for range 1000 {
		id := NewUUID()
		if !ValidUUID(id) {
			t.Fatalf("invalid uuid %q", id)
		}
		if id != strings.ToLower(id) {
			t.Fatalf("uuid not lowercase: %q", id)
		}
		if id[14] != '4' {
			t.Fatalf("version nibble not 4: %q", id)
		}
		if v := id[19]; v != '8' && v != '9' && v != 'a' && v != 'b' {
			t.Fatalf("variant not RFC 9562: %q", id)
		}
		raw, err := hex.DecodeString(strings.ReplaceAll(id, "-", ""))
		if err != nil || len(raw) != 16 {
			t.Fatalf("not 16 hex bytes: %q", id)
		}
		if seen[id] {
			t.Fatalf("duplicate uuid %q", id)
		}
		seen[id] = true
	}
}

func TestValidUUID(t *testing.T) {
	valid := []string{
		"123e4567-e89b-42d3-a456-426614174000",
		"123E4567-E89B-42D3-A456-426614174000",
		"00000000-0000-0000-0000-000000000000",
	}
	invalid := []string{
		"",
		"delete",
		"123e4567e89b42d3a456426614174000",
		"123e4567-e89b-42d3-a456-42661417400",
		"123e4567-e89b-42d3-a456-4266141740000",
		"123e4567-e89b-42d3-a456_426614174000",
		"g23e4567-e89b-42d3-a456-426614174000",
		"123e4567-e89b-42d3-a456-42661417400\x00",
		"' OR 1=1 --",
	}
	for _, s := range valid {
		if !ValidUUID(s) {
			t.Errorf("ValidUUID(%q) = false", s)
		}
	}
	for _, s := range invalid {
		if ValidUUID(s) {
			t.Errorf("ValidUUID(%q) = true", s)
		}
	}
}

// Invalid ids must short-circuit to ErrNotFound without touching the DB
// (the zero Postgres has no connection at all).
func TestInvalidIDIsNotFound(t *testing.T) {
	p, err := Open("postgres://u:p@127.0.0.1:1/db?sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}
	defer p.Close()
	if _, err := p.Get(context.Background(), "nope"); err != ErrNotFound {
		t.Errorf("Get = %v, want ErrNotFound", err)
	}
	if err := p.Delete(context.Background(), "nope"); err != ErrNotFound {
		t.Errorf("Delete = %v, want ErrNotFound", err)
	}
	if err := p.DeleteMany(context.Background(), nil); err != nil {
		t.Errorf("DeleteMany(nil) = %v, want nil", err)
	}
}

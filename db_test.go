package main

import (
	"database/sql"
	"testing"
)

func testDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := openDB(":memory:")
	if err != nil {
		t.Fatalf("openDB: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	return db
}

func TestInsertAndGet(t *testing.T) {
	db := testDB(t)
	id, err := insertLink(db, "https://example.com")
	if err != nil {
		t.Fatalf("insertLink: %v", err)
	}
	url, err := getURL(db, id)
	if err != nil {
		t.Fatalf("getURL: %v", err)
	}
	if url != "https://example.com" {
		t.Errorf("getURL = %q, want %q", url, "https://example.com")
	}
}

func TestDeleteThenGet(t *testing.T) {
	db := testDB(t)
	id, err := insertLink(db, "https://example.com")
	if err != nil {
		t.Fatalf("insertLink: %v", err)
	}
	if err := deleteLink(db, id); err != nil {
		t.Fatalf("deleteLink: %v", err)
	}
	if _, err := getURL(db, id); err == nil {
		t.Error("getURL after delete should fail")
	}
}

func TestListLinks(t *testing.T) {
	db := testDB(t)
	insertLink(db, "https://a.com")
	insertLink(db, "https://b.com")
	links, err := listLinks(db)
	if err != nil {
		t.Fatalf("listLinks: %v", err)
	}
	if len(links) != 2 {
		t.Fatalf("listLinks returned %d rows, want 2", len(links))
	}
	if links[0].URL != "https://a.com" || links[1].URL != "https://b.com" {
		t.Errorf("unexpected URLs: %v", links)
	}
}

func TestGetNotFound(t *testing.T) {
	db := testDB(t)
	if _, err := getURL(db, 999); err == nil {
		t.Error("getURL(999) should fail")
	}
}

func TestDeleteNotFound(t *testing.T) {
	db := testDB(t)
	if err := deleteLink(db, 999); err == nil {
		t.Error("deleteLink(999) should fail")
	}
}

package kv

import "testing"

func TestStore(t *testing.T) {
	store := NewStore()

	store.Put("book:B101", "Computer Networks")

	value, exists := store.Get("book:B101")

	if !exists {
		t.Fatal("book not found")
	}

	if value != "Computer Networks" {
		t.Fatal("wrong book value")
	}

	store.Delete("book:B101")

	_, exists = store.Get("book:B101")

	if exists {
		t.Fatal("book was not deleted")
	}
}

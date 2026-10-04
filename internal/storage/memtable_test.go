package storage

import "testing"

func TestMemTable(t *testing.T) {
	memtable := NewMemTable()

	memtable.Put("B101", "Computer Networks")

	value, ok := memtable.Get("B101")

	if !ok {
		t.Fatal("book not found")
	}

	if value != "Computer Networks" {
		t.Fatalf("unexpected value: %s", value)
	}

	memtable.Delete("B101")

	_, ok = memtable.Get("B101")

	if ok {
		t.Fatal("book should have been deleted")
	}
}

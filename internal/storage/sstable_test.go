package storage

import (
	"os"
	"testing"
)

func TestSSTable(t *testing.T) {
	filename := "test-sstable.db"

	memtable := NewMemTable()

	memtable.Put("B102", "Operating Systems")
	memtable.Put("B101", "Computer Networks")

	err := WriteSSTable(filename, memtable)
	if err != nil {
		t.Fatal(err)
	}

	value, found, err := ReadSSTable(filename, "B101")
	if err != nil {
		t.Fatal(err)
	}

	if !found {
		t.Fatal("book not found")
	}

	if value != "Computer Networks" {
		t.Fatalf("unexpected value: %s", value)
	}

	os.Remove(filename)
}

package storage

import (
	"os"
	"testing"
)

func TestWAL(t *testing.T) {
	filename := "test-wal.log"

	wal, err := OpenWAL(filename)
	if err != nil {
		t.Fatal(err)
	}

	err = wal.Write("PUT B101 Computer Networks")
	if err != nil {
		t.Fatal(err)
	}

	err = wal.Close()
	if err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}

	expected := "PUT B101 Computer Networks\n"

	if string(data) != expected {
		t.Fatalf("expected %q, got %q", expected, string(data))
	}

	os.Remove(filename)
}

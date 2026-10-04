package storage

import (
	"os"
	"testing"
)

func newTestLSM(t *testing.T) *LSMTree {
	t.Helper()

	dir := t.TempDir()

	lsm, err := NewLSMTree(dir, 3)
	if err != nil {
		t.Fatal(err)
	}

	return lsm
}

func TestLSMPutAndGet(t *testing.T) {
	lsm := newTestLSM(t)
	defer lsm.Close()

	err := lsm.Put("book1", "Computer Networks")
	if err != nil {
		t.Fatal(err)
	}

	value, err := lsm.Get("book1")
	if err != nil {
		t.Fatal(err)
	}

	if value != "Computer Networks" {
		t.Fatalf("expected Computer Networks, got %s", value)
	}
}

func TestLSMDelete(t *testing.T) {
	lsm := newTestLSM(t)
	defer lsm.Close()

	err := lsm.Put("book1", "Computer Networks")
	if err != nil {
		t.Fatal(err)
	}

	err = lsm.Delete("book1")
	if err != nil {
		t.Fatal(err)
	}

	_, err = lsm.Get("book1")

	if err != ErrKeyNotFound {
		t.Fatalf("expected key not found, got %v", err)
	}
}

func TestLSMSizeAndFlush(t *testing.T) {
	dir := t.TempDir()

	lsm, err := NewLSMTree(dir, 2)
	if err != nil {
		t.Fatal(err)
	}
	defer lsm.Close()

	lsm.Put("book1", "Computer Networks")

	if lsm.memTableSize() != 1 {
		t.Fatalf("expected size 1")
	}

	lsm.Put("book2", "Operating Systems")

	files, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}

	foundSSTable := false

	for _, file := range files {
		if len(file.Name()) >= 8 &&
			file.Name()[:8] == "sstable-" {
			foundSSTable = true
		}
	}

	if !foundSSTable {
		t.Fatal("SSTable was not created")
	}
}

func TestLSMManualFlush(t *testing.T) {
	dir := t.TempDir()

	lsm, err := NewLSMTree(dir, 100)
	if err != nil {
		t.Fatal(err)
	}
	defer lsm.Close()

	err = lsm.Put("book1", "Computer Networks")
	if err != nil {
		t.Fatal(err)
	}

	err = lsm.Flush()
	if err != nil {
		t.Fatal(err)
	}

	_, err = os.Stat(dir + "/sstable-1.db")
	if err != nil {
		t.Fatalf("SSTable was not created: %v", err)
	}
}

func TestLSMWALCreated(t *testing.T) {
	dir := t.TempDir()

	lsm, err := NewLSMTree(dir, 10)
	if err != nil {
		t.Fatal(err)
	}

	err = lsm.Put("book1", "Computer Networks")
	if err != nil {
		t.Fatal(err)
	}

	err = lsm.Close()
	if err != nil {
		t.Fatal(err)
	}

	_, err = os.Stat(dir + "/wal.log")
	if err != nil {
		t.Fatalf("WAL was not created: %v", err)
	}
}
func TestLSMGetFromSSTable(t *testing.T) {
	dir := t.TempDir()

	lsm, err := NewLSMTree(dir, 2)
	if err != nil {
		t.Fatal(err)
	}
	defer lsm.Close()

	err = lsm.Put("book1", "Computer Networks")
	if err != nil {
		t.Fatal(err)
	}

	err = lsm.Put("book2", "Operating Systems")
	if err != nil {
		t.Fatal(err)
	}

	// The MemTable should automatically flush
	// because the maximum size is 2.
	value, err := lsm.Get("book1")
	if err != nil {
		t.Fatal(err)
	}

	if value != "Computer Networks" {
		t.Fatalf(
			"expected Computer Networks, got %s",
			value,
		)
	}
}
func TestLSMWALRecovery(t *testing.T) {
	dir := t.TempDir()

	// First instance.
	lsm1, err := NewLSMTree(dir, 100)
	if err != nil {
		t.Fatal(err)
	}

	err = lsm1.Put("book1", "Computer Networks")
	if err != nil {
		t.Fatal(err)
	}

	err = lsm1.Close()
	if err != nil {
		t.Fatal(err)
	}

	// Second instance simulates restarting the program.
	lsm2, err := NewLSMTree(dir, 100)
	if err != nil {
		t.Fatal(err)
	}
	defer lsm2.Close()

	value, err := lsm2.Get("book1")
	if err != nil {
		t.Fatal(err)
	}

	if value != "Computer Networks" {
		t.Fatalf(
			"expected Computer Networks, got %s",
			value,
		)
	}
}
func TestLSMMultipleSSTables(t *testing.T) {
	dir := t.TempDir()

	lsm, err := NewLSMTree(dir, 2)
	if err != nil {
		t.Fatal(err)
	}
	defer lsm.Close()

	// First two entries create SSTable 1.
	if err := lsm.Put("book1", "Computer Networks"); err != nil {
		t.Fatal(err)
	}

	if err := lsm.Put("book2", "Operating Systems"); err != nil {
		t.Fatal(err)
	}

	// Next two entries create SSTable 2.
	if err := lsm.Put("book3", "Database Systems"); err != nil {
		t.Fatal(err)
	}

	if err := lsm.Put("book4", "Distributed Systems"); err != nil {
		t.Fatal(err)
	}

	// Check data from the first SSTable.
	value, err := lsm.Get("book1")
	if err != nil {
		t.Fatal(err)
	}

	if value != "Computer Networks" {
		t.Fatalf("expected Computer Networks, got %s", value)
	}

	// Check data from the second SSTable.
	value, err = lsm.Get("book4")
	if err != nil {
		t.Fatal(err)
	}

	if value != "Distributed Systems" {
		t.Fatalf("expected Distributed Systems, got %s", value)
	}
}
func TestLSMCompaction(t *testing.T) {
	dir := t.TempDir()

	lsm, err := NewLSMTree(dir, 2)
	if err != nil {
		t.Fatal(err)
	}
	defer lsm.Close()

	// Create SSTable 1.
	if err := lsm.Put("book1", "Computer Networks"); err != nil {
		t.Fatal(err)
	}

	if err := lsm.Put("book2", "Operating Systems"); err != nil {
		t.Fatal(err)
	}

	// Create SSTable 2.
	if err := lsm.Put("book3", "Database Systems"); err != nil {
		t.Fatal(err)
	}

	if err := lsm.Put("book4", "Distributed Systems"); err != nil {
		t.Fatal(err)
	}

	// Compact the SSTables.
	if err := lsm.Compact(); err != nil {
		t.Fatal(err)
	}

	// Check that data is still available.
	value, err := lsm.Get("book1")
	if err != nil {
		t.Fatal(err)
	}

	if value != "Computer Networks" {
		t.Fatalf(
			"expected Computer Networks, got %s",
			value,
		)
	}

	value, err = lsm.Get("book4")
	if err != nil {
		t.Fatal(err)
	}

	if value != "Distributed Systems" {
		t.Fatalf(
			"expected Distributed Systems, got %s",
			value,
		)
	}
}
func TestDeleteAndCompaction(t *testing.T) {
	dir := t.TempDir()

	lsm, err := NewLSMTree(dir, 2)
	if err != nil {
		t.Fatal(err)
	}
	defer lsm.Close()

	// Create first SSTable.
	if err := lsm.Put("book1", "Computer Networks"); err != nil {
		t.Fatal(err)
	}

	if err := lsm.Put("book2", "Operating Systems"); err != nil {
		t.Fatal(err)
	}

	// Delete book1.
	if err := lsm.Delete("book1"); err != nil {
		t.Fatal(err)
	}

	// Create another SSTable.
	if err := lsm.Put("book3", "Database Systems"); err != nil {
		t.Fatal(err)
	}

	if err := lsm.Put("book4", "Distributed Systems"); err != nil {
		t.Fatal(err)
	}

	// Compact.
	if err := lsm.Compact(); err != nil {
		t.Fatal(err)
	}

	// book1 must remain deleted.
	_, err = lsm.Get("book1")

	if err != ErrKeyNotFound {
		t.Fatalf("expected book1 to be deleted, got error: %v", err)
	}

	// Other books must still exist.
	value, err := lsm.Get("book2")
	if err != nil {
		t.Fatal(err)
	}

	if value != "Operating Systems" {
		t.Fatalf("unexpected value: %s", value)
	}
}

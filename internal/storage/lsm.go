package storage

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

var ErrKeyNotFound = errors.New("key not found")

type LSMTree struct {
	mu sync.RWMutex

	memTable *MemTable
	wal      *WAL

	dataDir string

	maxEntries int
	sstableID  int
}

func NewLSMTree(dataDir string, maxEntries int) (*LSMTree, error) {
	if maxEntries <= 0 {
		maxEntries = 100
	}

	err := os.MkdirAll(dataDir, 0755)
	if err != nil {
		return nil, err
	}

	walPath := filepath.Join(dataDir, "wal.log")

	wal, err := OpenWAL(walPath)
	if err != nil {
		return nil, err
	}

	nextID, err := findNextSSTableID(dataDir)
	if err != nil {
		wal.Close()
		return nil, err
	}

	lsm := &LSMTree{
		memTable:   NewMemTable(),
		wal:        wal,
		dataDir:    dataDir,
		maxEntries: maxEntries,
		sstableID:  nextID - 1,
	}

	err = lsm.recover()
	if err != nil {
		wal.Close()
		return nil, err
	}

	return lsm, nil
}

func (l *LSMTree) Put(key string, value string) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	operation := fmt.Sprintf("PUT|%s|%s", key, value)

	err := l.wal.Write(operation)
	if err != nil {
		return err
	}

	l.memTable.Put(key, value)

	if l.memTableSize() >= l.maxEntries {
		return l.flush()
	}

	return nil
}

func (l *LSMTree) Get(key string) (string, error) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	// First check the MemTable.
	value, ok := l.memTable.Get(key)
	if ok {
		return value, nil
	}

	// If not found in memory, check SSTables on disk.
	value, err := l.loadSSTable(key)
	if err != nil {
		return "", err
	}

	return value, nil
}
func (l *LSMTree) loadSSTable(key string) (string, error) {
	files, err := os.ReadDir(l.dataDir)
	if err != nil {
		return "", err
	}

	// Search SSTable files.
	for i := len(files) - 1; i >= 0; i-- {
		file := files[i]

		if !strings.HasPrefix(file.Name(), "sstable-") {
			continue
		}

		filename := filepath.Join(l.dataDir, file.Name())

		value, found, err := ReadSSTable(filename, key)
		if err != nil {
			return "", err
		}

		if found {
			return value, nil
		}
	}

	return "", ErrKeyNotFound
}

func (l *LSMTree) Delete(key string) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	operation := fmt.Sprintf("DELETE|%s", key)

	err := l.wal.Write(operation)
	if err != nil {
		return err
	}

	l.memTable.Delete(key)

	return nil
}

func (l *LSMTree) memTableSize() int {
	return len(l.memTable.Keys())
}

func (l *LSMTree) flush() error {
	if l.memTableSize() == 0 {
		return nil
	}

	l.sstableID++

	filename := filepath.Join(
		l.dataDir,
		"sstable-"+strconv.Itoa(l.sstableID)+".db",
	)

	err := WriteSSTable(filename, l.memTable)
	if err != nil {
		return err
	}

	l.memTable = NewMemTable()

	return nil
}

func (l *LSMTree) Flush() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	return l.flush()
}

func (l *LSMTree) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.memTableSize() > 0 {
		err := l.flush()
		if err != nil {
			return err
		}
	}

	return l.wal.Close()
}
func (l *LSMTree) recover() error {
	walPath := filepath.Join(l.dataDir, "wal.log")

	operations, err := l.wal.ReadAll(walPath)
	if err != nil {
		return err
	}

	for _, operation := range operations {
		parts := strings.SplitN(operation, "|", 3)

		if len(parts) == 0 {
			continue
		}

		switch parts[0] {
		case "PUT":
			if len(parts) != 3 {
				continue
			}

			l.memTable.Put(parts[1], parts[2])

		case "DELETE":
			if len(parts) != 2 {
				continue
			}

			l.memTable.Delete(parts[1])
		}
	}

	return nil
}
func findNextSSTableID(dataDir string) (int, error) {
	files, err := os.ReadDir(dataDir)
	if err != nil {
		return 0, err
	}

	maxID := 0

	for _, file := range files {
		name := file.Name()

		if !strings.HasPrefix(name, "sstable-") ||
			!strings.HasSuffix(name, ".db") {
			continue
		}

		number := strings.TrimPrefix(name, "sstable-")
		number = strings.TrimSuffix(number, ".db")

		id, err := strconv.Atoi(number)
		if err != nil {
			continue
		}

		if id > maxID {
			maxID = id
		}
	}

	return maxID + 1, nil
}
func (l *LSMTree) Compact() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	files, err := os.ReadDir(l.dataDir)
	if err != nil {
		return err
	}

	allData := make(map[string]string)
	deleted := make(map[string]bool)

	var sstableFiles []string

	// Read all SSTables.
	for _, file := range files {
		if !strings.HasPrefix(file.Name(), "sstable-") ||
			!strings.HasSuffix(file.Name(), ".db") {
			continue
		}

		sstableFiles = append(sstableFiles, file.Name())

		filename := filepath.Join(l.dataDir, file.Name())

		content, err := os.ReadFile(filename)
		if err != nil {
			return err
		}

		lines := strings.Split(string(content), "\n")

		for _, line := range lines {
			if line == "" {
				continue
			}

			parts := strings.SplitN(line, "=", 2)

			if len(parts) != 2 {
				continue
			}

			key := parts[0]
			value := parts[1]

			if value == "<deleted>" {
				delete(allData, key)
				deleted[key] = true
				continue
			}

			// A newer value replaces an older value.
			if !deleted[key] {
				allData[key] = value
			}
		}
	}

	// Nothing to compact.
	if len(sstableFiles) < 2 {
		return nil
	}

	// Create merged MemTable.
	merged := NewMemTable()

	for key, value := range allData {
		merged.Put(key, value)
	}

	// Create new SSTable.
	l.sstableID++

	filename := filepath.Join(
		l.dataDir,
		"sstable-"+strconv.Itoa(l.sstableID)+".db",
	)

	err = WriteSSTable(filename, merged)
	if err != nil {
		return err
	}

	// Remove old SSTables.
	for _, file := range sstableFiles {
		oldFile := filepath.Join(l.dataDir, file)

		if err := os.Remove(oldFile); err != nil {
			return err
		}
	}

	return nil
}

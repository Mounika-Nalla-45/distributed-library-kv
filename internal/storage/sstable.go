package storage

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func WriteSSTable(filename string, memtable *MemTable) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	for _, key := range memtable.Keys() {
		value, _ := memtable.Get(key)

		_, err := fmt.Fprintf(writer, "%s=%s\n", key, value)
		if err != nil {
			return err
		}
	}

	// Write tombstones.
	for _, key := range memtable.DeletedKeys() {
		_, err := fmt.Fprintf(writer, "%s=<deleted>\n", key)
		if err != nil {
			return err
		}
	}

	return writer.Flush()
}

func ReadSSTable(filename string, key string) (string, bool, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", false, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		parts := strings.SplitN(line, "=", 2)

		if len(parts) == 2 && parts[0] == key {
			if parts[1] == "<deleted>" {
				return "", false, nil
			}

			return parts[1], true, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return "", false, err
	}

	return "", false, nil
}

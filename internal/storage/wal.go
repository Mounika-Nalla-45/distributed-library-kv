package storage

import (
	"bufio"
	"os"
)

type WAL struct {
	file *os.File
}

func OpenWAL(filename string) (*WAL, error) {
	file, err := os.OpenFile(
		filename,
		os.O_CREATE|os.O_APPEND|os.O_WRONLY,
		0644,
	)

	if err != nil {
		return nil, err
	}

	return &WAL{file: file}, nil
}

func (w *WAL) Write(operation string) error {
	writer := bufio.NewWriter(w.file)

	_, err := writer.WriteString(operation + "\n")
	if err != nil {
		return err
	}

	return writer.Flush()
}
func (w *WAL) ReadAll(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var operations []string

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		operations = append(operations, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return operations, nil
}

func (w *WAL) Close() error {
	return w.file.Close()
}

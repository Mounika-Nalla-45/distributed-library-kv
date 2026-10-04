package storage

import (
	"fmt"
	"testing"
)

func BenchmarkLSMPut(b *testing.B) {
	dir := b.TempDir()

	lsm, err := NewLSMTree(dir, b.N+1)
	if err != nil {
		b.Fatal(err)
	}
	defer lsm.Close()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("book-%d", i)
		value := fmt.Sprintf("value-%d", i)

		if err := lsm.Put(key, value); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkLSMGet(b *testing.B) {
	dir := b.TempDir()

	lsm, err := NewLSMTree(dir, b.N+1)
	if err != nil {
		b.Fatal(err)
	}
	defer lsm.Close()

	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("book-%d", i)

		if err := lsm.Put(key, "library-book"); err != nil {
			b.Fatal(err)
		}
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("book-%d", i)

		_, err := lsm.Get(key)
		if err != nil {
			b.Fatal(err)
		}
	}
}

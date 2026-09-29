package main

import (
	"fmt"

	"distributed-library-kv/internal/kv"
	"distributed-library-kv/internal/library"
)

func main() {

	store := kv.NewStore()

	lib := library.NewLibrary(store)

	// Add books
	lib.AddBook("B101", "Computer Networks")
	lib.AddBook("B102", "Operating Systems")
	lib.AddBook("B103", "Data Structures")

	// Get a book
	book, exists := lib.GetBook("B101")

	if exists {
		fmt.Println("Book:", book)
	}

	// Delete a book
	lib.DeleteBook("B103")

	_, exists = lib.GetBook("B103")

	if !exists {
		fmt.Println("B103 deleted successfully")
	}
}

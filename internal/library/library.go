package library

import "distributed-library-kv/internal/kv"

type Library struct {
	store *kv.Store
}

func NewLibrary(store *kv.Store) *Library {
	return &Library{
		store: store,
	}
}

func (l *Library) AddBook(id string, title string) {
	key := "book:" + id
	l.store.Put(key, title)
}

func (l *Library) GetBook(id string) (string, bool) {
	key := "book:" + id
	return l.store.Get(key)
}

func (l *Library) DeleteBook(id string) {
	key := "book:" + id
	l.store.Delete(key)
}

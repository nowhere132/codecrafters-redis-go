package main

type Storage struct {
	data map[string]string
}

func NewStorage() *Storage {
	return &Storage{data: make(map[string]string)}
}

func (store *Storage) Get(k string) string {
	return store.data[k]
}

func (store *Storage) Set(k string, v string) {
	store.data[k] = v
}

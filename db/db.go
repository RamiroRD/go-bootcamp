package db

type Database interface {
	Insert(key, value string) bool
	Delete(key string) bool
	Update(key string, value string) bool
	Lookup(key string) (v string, ok bool)
	All() []string
}

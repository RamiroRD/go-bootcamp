package db

import (
	"sync"
	"time"
)

/*
	Database database. Part of Go Bootcamp
*/

type Database struct {
	contents map[int]*Student
	mutex    sync.RWMutex
}

type Student struct {
	Id       int
	Name     string
	LastName string
	Birthday time.Time
}

func NewDatabase() *Database {
	return &Database{contents: make(map[int]*Student)}
}

func (o *Database) Insert(s *Student) bool {
	if s == nil {
		return false
	}
	o.mutex.Lock()
	defer o.mutex.Unlock()

	_, ok := o.contents[s.Id]
	if ok {
		return false
	}
	o.contents[s.Id] = s
	return true

}

func (o *Database) Delete(s *Student) {
	if s == nil {
		return
	}
	o.mutex.Lock()
	delete(o.contents, s.Id)
	o.mutex.Unlock()
}

func (o *Database) Update(id int, s *Student) {
	if s == nil {
		return
	}
	o.mutex.Lock()
	t, ok := o.contents[id]
	if !ok {
		return
	}

	t.Name = s.Name
	t.LastName = s.LastName
	t.Birthday = s.Birthday
	o.mutex.Unlock()
}

func (o *Database) Lookup(id int) (s *Student, ok bool) {
	o.mutex.RLock()
	s, ok = o.contents[id]
	o.mutex.RUnlock()
	return
}

func (o *Database) All() []*Student {
	o.mutex.RLock()
	defer o.mutex.RUnlock()
	v := make([]*Student, 0, len(o.contents))
	for _, value := range o.contents {
		v = append(v, value)
	}
	return v
}

func Save(o *Database) {
}

func Restore(o *Database) {
}

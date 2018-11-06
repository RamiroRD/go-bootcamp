package db

import (
	"github.com/ramirord/go-bootcamp/util"
	"sync"
)

/*
	Simple key-value database.
*/
type FileDatabase struct {
	contents map[string]string
	mutex    sync.RWMutex
}

func NewFileDatabase() *FileDatabase {
	return &FileDatabase{contents: make(map[string]string)}
}

func (d *FileDatabase) Insert(key, value string) bool {
	if util.ContainsEmpty(key, value) {
		return false
	}

	d.mutex.Lock()
	defer d.mutex.Unlock()

	_, ok := d.contents[key]
	if ok {
		return false
	}
	d.contents[key] = value
	return true

}

func (d *FileDatabase) Delete(key string) bool {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	_, ok := d.contents[key]
	delete(d.contents, key)
	return ok
}

func (d *FileDatabase) Update(key, value string) bool {
	if util.ContainsEmpty(key, value) {
		return false
	}
	d.mutex.Lock()
	defer d.mutex.Unlock()

	_, ok := d.contents[key]
	if !ok {
		return false
	} else {
		delete(d.contents, key)
		d.contents[key] = value
		return true
	}
}

func (d *FileDatabase) Lookup(key string) (v string, ok bool) {
	d.mutex.RLock()
	v, ok = d.contents[key]
	d.mutex.RUnlock()
	return
}

func (d *FileDatabase) All() []string {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	v := make([]string, 0, len(d.contents))
	for _, value := range d.contents {
		v = append(v, value)
	}
	return v
}

func (d *FileDatabase) Save() {
	// TODO implementar
}

func (d *FileDatabase) Restore() {
	// TODO implementar
}

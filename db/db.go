package db

type Database interface {
	Insert(s *Student) bool
	Delete(s *Student)
	Update(id int, s *Student)
	Lookup(id int) (s *Student, ok bool)
	All() []*Student
}

package db

import "testing"
import "time"

func quickDate(year, month, day int) time.Time {
	loc, _ := time.LoadLocation("UTC")
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, loc)
}

func TestInsert(t *testing.T) {
	o := NewFileDatabase()
	elems := []Student{
		{1231232, "María", "Eufrasia", quickDate(1923, 12, 2)},
		{1636733, "José", "Ridículo", quickDate(1967, 3, 3)},
		{5691032, "Jemand", "Derexistiert", quickDate(1932, 7, 25)},
	}

	nelems := []int{23, 24, 26}
	t.Run("Insert", func(t *testing.T) {
		for _, p := range elems {
			o.Insert(&p)
			_, ok := o.Lookup(p.Id)
			if !ok {
				return
			}
			for _, id := range nelems {
				_, ok := o.Lookup(id)
				if ok {
					return
				}
			}
		}
		return
	})
}

func TestLookup(t *testing.T) {
	t.Run("Lookup", func(t *testing.T) {
	})
}

func TestDelete(t *testing.T) {
	t.Run("Delete", func(t *testing.T) {
	})
}

func TestUpdate(t *testing.T) {
	t.Run("Update", func(t *testing.T) {
	})
}

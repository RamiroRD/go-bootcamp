package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/ramirord/go-bootcamp/db"
	"net/http"
	"strconv"
	"time"
)

type Server struct {
	database db.Database
	router   *mux.Router
}

func quickDate(year, month, day int) time.Time {
	// TODO manejo de errores
	loc, _ := time.LoadLocation("UTC")
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, loc)
}

func containsEmpty(ss ...string) bool {
	for _, s := range ss {
		if s == "" {
			return true
		}
	}
	return false
}

func badRequest(w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("Bad request"))
}

type integerParser struct {
	err error
}

func (p *integerParser) parse(s string) int {
	if p.err != nil {
		return 0
	}
	n, err := strconv.Atoi(s)
	p.err = err
	return n
}

func (server *Server) studentsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		b, err := json.Marshal(server.database.All())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(b)
	case http.MethodPost:
		sid := r.FormValue("id")
		name := r.FormValue("name")
		lastname := r.FormValue("last_name")
		syear := r.FormValue("year")
		smonth := r.FormValue("month")
		sday := r.FormValue("day")

		if containsEmpty(name, lastname) {
			badRequest(w)
			return
		}

		var parser integerParser
		id := parser.parse(sid)
		year := parser.parse(syear)
		month := parser.parse(smonth)
		day := parser.parse(sday)

		if parser.err != nil {
			badRequest(w)
			return
		}

		s := db.Student{
			Id:       id,
			Name:     name,
			LastName: lastname,
			Birthday: quickDate(year, month, day),
		}

		ok := server.database.Insert(&s)
		if ok {
			w.Header().Add("Location", "/students/"+strconv.Itoa(id))
			w.WriteHeader(http.StatusCreated)
		} else {
			badRequest(w)
		}

	case http.MethodPut, http.MethodDelete:
		w.WriteHeader(http.StatusNotFound)
	default:
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte("Method not allowed"))
	}
}

func (server *Server) studentHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		values := mux.Vars(r)
		_, err := strconv.ParseUint(values["id"], 10, 64)
		if err == nil {
			w.WriteHeader(http.StatusBadRequest)
		}
		// TODO implementar
		w.WriteHeader(http.StatusOK)
	case http.MethodPost:
		w.WriteHeader(http.StatusMethodNotAllowed)
	case http.MethodPut, http.MethodDelete:
		// TODO implementar
		w.WriteHeader(http.StatusNotImplemented)
	default:
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte("Method not allowed"))
	}
}

func (server *Server) Serve(db db.Database) {
	server.router = mux.NewRouter()
	server.database = db

	server.router.HandleFunc("/students", func(w http.ResponseWriter, r *http.Request) {
		server.studentsHandler(w, r)
	})

	server.router.HandleFunc("/students/{id}", func(w http.ResponseWriter, r *http.Request) {
		server.studentHandler(w, r)
	})

	http.ListenAndServe(":8080", server.router)
}

package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/ramirord/go-bootcamp/db"
	"github.com/ramirord/go-bootcamp/util"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Server struct {
	database db.Database
	router   *mux.Router
}

type student struct {
	Id       int       `json:"id"`
	Name     string    `json:"name"`
	LastName string    `json:"last_name"`
	Birthday time.Time `json:"birthday"`
}

func quickDate(year, month, day int) time.Time {
	// TODO manejo de errores
	loc, _ := time.LoadLocation("UTC")
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, loc)
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

func (server *Server) listStudents(w http.ResponseWriter, r *http.Request) {
	b, err := json.Marshal(server.database.All())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func (server *Server) createStudent(w http.ResponseWriter, r *http.Request) {
	sid := r.FormValue("id")
	name := r.FormValue("name")
	lastname := r.FormValue("last_name")
	syear := r.FormValue("year")
	smonth := r.FormValue("month")
	sday := r.FormValue("day")

	if util.ContainsEmpty(name, lastname) {
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

	s := student{
		Id:       id,
		Name:     name,
		LastName: lastname,
		Birthday: quickDate(year, month, day),
	}

	str, err := json.Marshal(s)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ok := server.database.Insert(sid, string(str))

	if ok {
		w.Header().Add("Location", "/students/"+strconv.Itoa(id))
		w.WriteHeader(http.StatusCreated)
	} else {
		badRequest(w)
	}
}

func (server *Server) retrieveStudent(w http.ResponseWriter, r *http.Request) {
	values := mux.Vars(r)
	_, err := strconv.ParseUint(values["id"], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	v, ok := server.database.Lookup(values["id"])
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Write([]byte(v))
	w.WriteHeader(http.StatusOK)
}

func (server *Server) updateStudent(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (server *Server) deleteStudent(w http.ResponseWriter, r *http.Request) {
	values := mux.Vars(r)

	_, err := strconv.ParseUint(values["id"], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !server.database.Delete(values["id"]) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)

}

func logRequest(handler func(w http.ResponseWriter, r *http.Request)) func(http.ResponseWriter, *http.Request) {
	//log.Printf("%s from\n", r)
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		handler(w, r)
		//log.Printf("%s request from %s, response %s,  elapsed %s", r.Method, r.RemoteAddr, w.Header().Get("Status-Code"), time.Since(start))
		log.Printf("%s request from %s, elapsed %s", r.Method, r.RemoteAddr, time.Since(start))
	}
}

func (server *Server) Serve(db db.Database) {
	server.router = mux.NewRouter()
	server.database = db

	server.router.HandleFunc("/students", logRequest(server.listStudents)).Methods(http.MethodGet)
	server.router.HandleFunc("/students", logRequest(server.createStudent)).Methods(http.MethodPost)
	server.router.HandleFunc("/students/{id}", logRequest(server.retrieveStudent)).Methods(http.MethodGet)
	server.router.HandleFunc("/students/{id}", logRequest(server.updateStudent)).Methods(http.MethodPatch)
	server.router.HandleFunc("/students/{id}", logRequest(server.deleteStudent)).Methods(http.MethodDelete)

	http.ListenAndServe(":8080", server.router)
}

package main

import "github.com/ramirord/go-bootcamp/api"
import "github.com/ramirord/go-bootcamp/db"

var s api.Server

func main() {
	s.Serve(db.NewFileDatabase())
}

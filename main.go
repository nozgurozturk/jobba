package main

import (
	"github.com/nozgurozturk/jobba/app"
)

func main () {
	server := app.Server{}
	server.Run()
}
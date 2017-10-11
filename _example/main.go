package main

import (
	"flag"
	"log"

	"github.com/arteev/uriban"
)

var connectionString = flag.String("connection", "postgres://pqgotest:password@localhost/pqgotest?sslmode=verify-full", "connection string")

func main() {
	flag.Parse()
	bs := uriban.Replace(*connectionString, uriban.WithOption(uriban.Password, uriban.ModeValue("SECRET")))
	log.Printf("Connection string: %s", bs)
}

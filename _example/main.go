package main

import (
	"flag"
	"log"

	"github.com/arteev/uriban"
)

var connectionString = flag.String("connection", "postgres://pqgotest:password@myhost/pqgotest?sslmode=verify-full", "connection string")

func main() {
	flag.Parse()
	bs := uriban.Replace(*connectionString,
		uriban.WithOption(uriban.Password, uriban.ModeValue("SECRET")),
		uriban.WithOption(uriban.Host, uriban.ModeValue("localhost")))
	log.Printf("Connection string: %s", bs)
}

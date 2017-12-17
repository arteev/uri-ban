# uriban

[![Build Status](https://travis-ci.org/arteev/uriban.svg?branch=master)](https://travis-ci.org/arteev/uriban)
[![Coverage Status](https://coveralls.io/repos/arteev/uriban/badge.svg?branch=master&service=github)](https://coveralls.io/github/arteev/uriban?branch=master)
[![GoDoc](https://godoc.org/github.com/arteev/uriban?status.png)](https://godoc.org/github.com/arteev/uriban)

Description
-----------

Golang package for hiding information in a URI

Installation
------------

This package can be installed with the go get command:

    go get github.com/arteev/uriban

Documentation
-------------
Example:

```go
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
	log.Printf("Connection string:%s", bs)
}`
```

Output: 

```sh
2017/10/11 13:54:14 Connection string: postgres://pqgotest:SECRET@localhost/pqgotest?sslmode=verify-full
```

License
-------

  MIT


Author
------

Arteev Aleksey
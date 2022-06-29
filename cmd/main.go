package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/szucik/go-simple-rest-api/internal/server"
)

func main() {
	server.Run()
}

package main

import (
	"log"
	"os"
)

var logger = log.New(os.Stdout, "server ", log.Ldate|log.Ltime)

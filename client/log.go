package main

import (
	"log"
	"os"
)

var logger = log.New(os.Stdout, "client ", log.Ldate | log.Ltime)

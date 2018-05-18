package message

import (
	"log"
	"os"
)

var logger = log.New(os.Stdout, "message ", log.Ldate | log.Ltime)

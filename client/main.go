package main

import (
	"os"
	"net"
	"io"
	"fmt"
	"strconv"

	"github.com/satori/go.uuid"
	"github.com/frostoov/tunnel/message"
)



func serveTunnel(portb int16, addr string) {
	copier := func(aConn net.Conn, bConn net.Conn) {
		io.Copy(aConn, bConn)
		bConn.Close()
	}

	conna, err := net.Dial("tcp", addr)
	if err != nil {
		logger.Fatalf("%Failed dial s", err)
	}
	req := message.Request{
		Id: uuid.Must(uuid.NewV4()).String(),
		Name: "serveTunnel",
	}
	if err := req.Write(conna); err != nil {
		logger.Printf("Failed send serve tunnel %s", err)
		return
	}

	connb, err := net.Dial("tcp", fmt.Sprintf(":%d", portb))
	if err != nil {
		logger.Printf("Failed connect to portb %s", err)
		return
	}

	go copier(conna, connb)
	go copier(connb, conna)
}



func main() {
	var (
		req message.Request
		resp message.Response
	)

	portb, _ := strconv.Atoi(os.Args[1])
	addr := os.Args[2]
	port, err := strconv.Atoi(os.Args[3])
	if err != nil {
		logger.Fatalf("failed parse port %s", err)
	}
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		logger.Fatalf("Failed connect to server %s: %s", addr, err)
	}


	req = message.Request{
		Id: uuid.Must(uuid.NewV4()).String(),
		Name: "makeTunnel",
		Inputs: message.Args{
			"port": port,
		},
	}
	if err := req.Write(conn); err != nil {
		logger.Fatalf("Failed write request %s", err)
	}
	if err := resp.Read(conn); err != nil {
		logger.Fatalf("Failed read response %s", err)
	}
	for {
		if err := req.Read(conn); err != nil {
			logger.Fatalf("Failed read request %s", err)
		}
		if req.Name == "new_connection" {
			serveTunnel(int16(portb), addr)
		}
	}
}

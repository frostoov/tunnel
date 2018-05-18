package main

import (
	"fmt"
	"io"
	"net"

	"github.com/frostoov/tunnel/message"
	"github.com/satori/go.uuid"
)

type Status string

var aConns = make(chan net.Conn)
var bConns = make(chan net.Conn)

func tunnel(aConns <-chan net.Conn, bConns <-chan net.Conn) {
	copier := func(aConn net.Conn, bConn net.Conn) {
		io.Copy(aConn, bConn)
		bConn.Close()
	}

	for aConn := range aConns {
		logger.Printf("aConn %s", aConn.RemoteAddr())
		bConn := <-bConns
		logger.Printf("tunnel!! %s <=> %s", aConn.RemoteAddr(), bConn.RemoteAddr())
		go copier(aConn, bConn)
		go copier(bConn, aConn)
	}
}

func handleMaster(master net.Conn, ln net.Listener) {
	for {
		conn, err := ln.Accept()
		if err != nil {
			logger.Printf("Failed accept connection %s", err)
			return
		}
		logger.Printf("accepted %s from %s", conn.RemoteAddr(), ln.Addr())
		aConns <- conn
		req := message.Request{
			Id:   uuid.Must(uuid.NewV4()).String(),
			Name: "new_connection",
		}
		if err := req.Write(master); err != nil {
			logger.Printf("Failed write request %s", err)
			return
		}
	}
}

func makeTunnel(conn net.Conn, inputs message.Args) (message.Args, error) {
	logger.Printf("makeTunnel start!")
	port := int16(inputs["port"].(float64))

	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	go handleMaster(conn, ln)
	return message.Args{"port": port}, nil
}

var handlers = map[string]func(net.Conn, message.Args) (message.Args, error){
	"makeTunnel": makeTunnel,
}

func handleRequest(conn net.Conn, name string, inputs message.Args) (message.Args, error) {
	return handlers[name](conn, inputs)
}

func handleConn(conn net.Conn) {
	var (
		req  message.Request
	)
	for {
		if err := req.Read(conn); err != nil {
			logger.Printf("Failed read request %s", err)
			return
		}
		if req.Name == "serveTunnel" {
			bConns <- conn
			return
		}
		outputs, err := handleRequest(conn, req.Name, req.Inputs)
		if err != nil {
			logger.Printf("%s failure %s", conn.RemoteAddr(), err)
			return
		}
		resp := message.Response{
			Id:      req.Id,
			Status:  "success",
			Outputs: outputs,
		}
		if err := resp.Write(conn); err != nil {
			logger.Printf("Failed write response %s", err)
			return
		}
	}
}

func Listen(addr string) error {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	for {
		conn, err := ln.Accept()
		logger.Printf("New conn!")
		if err != nil {
			logger.Printf("Failed accept %s", err)
		}
		go handleConn(conn)
	}
}

func main() {
	go tunnel(aConns, bConns)
	Listen(":6089")
}

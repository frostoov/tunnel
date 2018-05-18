package message

import (
	"encoding/binary"
	"encoding/json"
	"io"
)

type Args = map[string]interface{}

type Request struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Inputs Args   `json:"inputs"`
}

type Response struct {
	Id      string `json:"id"`
	Status  string `json:"status"`
	Outputs Args   `json:"outputs"`
}

func read(r io.Reader, o interface{}) error {
	size := uint64(0)
	if err := binary.Read(r, binary.BigEndian, &size); err != nil {
		return err
	}
	data := make([]byte, size)
	if err := binary.Read(r, binary.BigEndian, data); err != nil {
		return err
	}
	if err := json.Unmarshal(data, o); err != nil {
		return err
	}
	return nil
}

func write(w io.Writer, o interface{}) error {
	data, err := json.Marshal(o)
	if err != nil {
		return err
	}
	if err := binary.Write(w, binary.BigEndian, uint64(len(data))); err != nil {
		return err
	}
	if err := binary.Write(w, binary.BigEndian, data); err != nil {
		return err
	}
	return nil
}

func (req *Request) Read(r io.Reader) error {
	err := read(r, req)
	if err == nil {
		logger.Printf("read request %s:%s\n", req.Name, req.Id)
	}
	return err
}

func (req *Request) Write(w io.Writer) error {
	err := write(w, req)
	if err == nil {
		logger.Printf("write request %s:%s\n", req.Name, req.Id)
	}
	return err
}

func (resp *Response) Read(r io.Reader) error {
	err := read(r, resp)
	if err == nil {
		logger.Printf("read response %s\n", resp.Id)
	}
	return err
}

func (resp *Response) Write(w io.Writer) error {
	err := write(w, resp)
	if err == nil {
		logger.Printf("write respose %s\n", resp.Id)
	}
	return err
}

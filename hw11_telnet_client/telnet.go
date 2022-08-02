package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type client struct {
	address string
	timeout time.Duration
	conn    net.Conn
	in      io.ReadCloser
	out     io.Writer
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &client{address: address, timeout: timeout, in: in, out: out}
}

func (cl *client) Connect() error {
	var err error
	cl.conn, err = net.DialTimeout("tcp", cl.address, cl.timeout)
	if err != nil {
		return fmt.Errorf("cannot connect to %s with the error %w", cl.address, err)
	}
	fmt.Fprintln(os.Stderr, "...Connected to", cl.address)
	return nil
}

func (cl *client) Close() error {
	return cl.conn.Close()
}

func (cl *client) Send() error {
	scanner := bufio.NewScanner(cl.in)
	for scanner.Scan() {
		text := scanner.Text()
		log.Printf("SENDED: %s", text)

		_, err := cl.conn.Write([]byte(fmt.Sprintln(text)))
		if err != nil {
			return fmt.Errorf("...Connection was closed by peer")
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	fmt.Fprintln(os.Stderr, "...EOF")

	return nil
}

func (cl *client) Receive() error {
	scanner := bufio.NewScanner(cl.conn)
	for scanner.Scan() {
		text := scanner.Text()
		log.Printf("RECEIVED: %s", text)

		fmt.Fprintln(cl.out, text)
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

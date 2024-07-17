package main

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

func main() {
	go func() {
		time.Sleep(2 * time.Second)
		sendFile(20000000)
	}()
	server := &FileServer{}
	server.start()
}

type FileServer struct{}

func (fs *FileServer) start() {
	ln, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go fs.readLoop(conn)
	}
}

func (fs *FileServer) readLoop(conn net.Conn) {
	// buf is a bytes Buffer struct
	buf := new(bytes.Buffer)
	var size int64
	// binary.Read reads in structured binary from conn into &size
	// which is a fixed-size int64
	// LittleEndian is a byte ordering system where the least-significant
	// byte is stored at the smallest address.
	binary.Read(conn, binary.LittleEndian, &size)
	for {
		// CopyN copies n (size) number of bytes from buf to conn
		// or until it encounters an error and essentially returns
		// the number of bytes that were copied
		n, err := io.CopyN(buf, conn, size)
		if err != nil {
			log.Fatal(err)
		}
		// buf.Bytes() returns the bytes in the buf struct
		fmt.Println(buf.Bytes())
		// n is the number of bytes that were copied by CopyN
		fmt.Printf("received %d bytes over the network/n", n)
	}
}

func sendFile(size int) error {
	file := make([]byte, size)
	_, err := io.ReadFull(rand.Reader, file)
	if err != nil {
		return err
	}
	conn, err := net.Dial("tcp", ":3000")
	if err != nil {
		return err
	}
	binary.Write(conn, binary.LittleEndian, int64(size))
	n, err := io.CopyN(conn, bytes.NewReader(file), int64(size))
	if err != nil {
		return err
	}

	fmt.Printf("wrote %d bytes over the network\n", n)
	return nil
}

//go:build unix || !windows
// +build unix !windows

package main

// This is simple demo code to show how to fd-pass over a unix domain socket in go. To aid
// readability, all function return value tests have been elided excepting error.
//
// This client connects to our companion server and sends a request. The response includes
// an fd which originated as the server's Stdout. We copy our Stdin to this passed fd; in
// short, connecting our Stdin to their Stdout. After you run this, anything you type into
// the client's Stdin should show up on the server's Stdout.

import (
	"fmt"
	"io"
	"net"
	"os"
	"syscall"
)

func main() {
	ra := &net.UnixAddr{socketPath, network}
	conn, err := net.DialUnix(network, nil, ra)
	if err != nil {
		fmt.Println("DialUnix", err)
		return
	}
	_, _, err = conn.WriteMsgUnix([]byte("Send me your Stdout"), nil, nil)
	if err != nil {
		fmt.Println("WriteMsgUnix", err)
		return
	}
	msg := make([]byte, 1024)
	oob := make([]byte, 1024)
	_, oobn, _, _, err := conn.ReadMsgUnix(msg, oob)
	if err != nil {
		fmt.Println("ReadMsgUnix", err)
		return
	}
	if oobn <= 0 {
		fmt.Println("Error: No Out-Of-Band Data provided by server")
		return
	}

	// Decode the raw OOB bytes into control messages. OOB can legitimately contain
	// multiple control messages of different types.
	scms, err := syscall.ParseSocketControlMessage(oob[:oobn])
	if err != nil {
		fmt.Println("ParseSocketControlMessage", err)
		return
	}

	// Find the first control message which contains a passed fd.
	for _, scm := range scms {
		fds, err := syscall.ParseUnixRights(&scm)
		if err != nil {
			fmt.Println("ParseUnixRights", err)
			continue
		}
		for _, fd := range fds {
			f := os.NewFile(uintptr(fd), "Stdout")
			fmt.Println("Sending my Stdin to Server's Stdout as fd", fd)
			io.Copy(f, os.Stdin)
			f.Close()
			return // All done
		}
	}
}

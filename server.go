// +build unix !windows

package main

// This is simple demo code to show how to fd-pass over a unix domain socket in go. To aid
// readability, all function return value tests have been elided excepting error.
//
// On receipt of a message this server responds with an SCM_RIGHTS message to fd-pass our
// Stdout back to the client. From that point on, it's up to the client what it does with
// that fd. For more on fd-passing, consult the recvmsg(2) man page.
import (
	"fmt"
	"net"
	"os"
	"syscall"
)

func main() {
	os.Remove(socketPath)
	listener, err := net.ListenUnix(network, &net.UnixAddr{socketPath, network})
	if err != nil {
		fmt.Println("ListenUnix", err)
		return
	}
	fmt.Println("Server ready on", listener)

	for {
		conn, err := listener.AcceptUnix()
		if err != nil {
			fmt.Println("AcceptUnix", err)
			return
		}
		fmt.Println("New Client:", conn)

		msg := make([]byte, 1024)
		_, _, _, _, err = conn.ReadMsgUnix(msg, nil)
		if err != nil {
			fmt.Println("ReadMsgUnix", err)
			conn.Close()
			continue
		}
		fmt.Println("Client Request:", string(msg))

		oob := syscall.UnixRights(int(stdout)) // Construct an SCM_RIGHTS CMSG with our fd
		_, _, err = conn.WriteMsgUnix([]byte("Here is my stdout"), oob, nil)
		if err != nil {
			fmt.Println("WriteMsgUnix", err)
		}

		conn.Close()
	}
}

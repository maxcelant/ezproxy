package listener

import (
	"fmt"
	"golang.org/x/sys/unix"
	"net"
	"os"
	"syscall"
)

// To allow multiple ip/ports to listen on the same address,
// we need to set specific configurations on the socket for
// the kernel to be aware. Otherwise it will return a SOCKET IN USE error
func Listen(network, address string) (net.Listener, error) {
	// Resolve address to get sockaddr
	tcpAddr, err := net.ResolveTCPAddr(network, address)
	if err != nil {
		return nil, err
	}

	// Create socket
	fd, err := unix.Socket(unix.AF_INET, unix.SOCK_STREAM, 0)
	if err != nil {
		return nil, err
	}

	// Configures address and port reusability to the socket
	if err := unix.SetsockoptInt(fd, unix.SOL_SOCKET, unix.SO_REUSEADDR, 1); err != nil {
		return nil, err
	}
	if err := unix.SetsockoptInt(fd, unix.SOL_SOCKET, unix.SO_REUSEPORT, 1); err != nil {
		return nil, err
	}

	// Binds the socket to the ip:port
	sa := &unix.SockaddrInet4{Port: tcpAddr.Port}
	copy(sa.Addr[:], tcpAddr.IP.To4())
	if err := unix.Bind(fd, sa); err != nil {
		return nil, err
	}

	if err := unix.Listen(fd, syscall.SOMAXCONN); err != nil {
		return nil, err
	}

	// Wrap in *os.File then net.Listener
	file := os.NewFile(uintptr(fd), fmt.Sprintf("fd %d", fd))
	defer file.Close()
	ln, err := net.FileListener(file)
	if err != nil {
		return nil, err
	}
	return ln, nil
}

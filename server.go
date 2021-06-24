package main

import (
	"fmt"
	"net"
	"syscall"
	"unsafe"

	"golang.org/x/sys/unix"
)

const SO_ORIGINAL_DST = 80

func panicOnErr(ctx string, err error) {
	if err != nil {
		panic(fmt.Sprintf("%s: %w", ctx, err))
	}
}

func main() {
	l, err := net.Listen("tcp4", "0.0.0.0:6666")
	panicOnErr("net.Listen", err)
	defer l.Close()

	for {
		conn, err := l.Accept()
		panicOnErr("l.Accept", err)
		go handle(conn)
	}
}

func handle(conn net.Conn) {
	defer conn.Close()

	tcpConn := conn.(*net.TCPConn)

	file, err := tcpConn.File()
	panicOnErr("tcpConn.File", err)
	defer file.Close()
	fd := file.Fd()

	addr, err := syscall.GetsockoptIPv6Mreq(int(fd), syscall.IPPROTO_IP, SO_ORIGINAL_DST)
	panicOnErr("getsockopt", err)
	ip := fmt.Sprintf("%d.%d.%d.%d", uint(addr.Multiaddr[4]),
		uint(addr.Multiaddr[5]),
		uint(addr.Multiaddr[6]),
		uint(addr.Multiaddr[7]))
	port := uint16(addr.Multiaddr[2])<<8 + uint16(addr.Multiaddr[3])
	fmt.Println("SO_ORIGNAL_DST:", ip, port)

	buf := make([]byte, 1024)
	_, err = conn.Read(buf)
	panicOnErr("conn.Read", err)
	conn.Write(buf)
}

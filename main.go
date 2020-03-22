// Based on https://github.com/vsergeev/tinytaptunnel

package main

import (
	"encoding/binary"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"syscall"
	"time"
	"unsafe"
)

const (
	TIMESTAMP_SIZE    = 8               // TIMESTAMP_SIZE is the size of the timestamp
	TIMESTAMP_TIMEOUT = time.Second * 3 // TIMESTAMP_TIMEOUT is the maximum duration after which a frame will be discarded

	TCP_MTU = 1472                          // The TCP MTU
	TAP_MTU = TCP_MTU - 14 - TIMESTAMP_SIZE // TAP_MTU is TCP_MTU - ethernet header (14) - TIMESTAMP_SIZE
)

// EncapsulateFrame encapsulates a frame
// A frame is composed of a nanosecond timestamp (8 bytes) and a plaintext frame (1-1432 bytes)
func EncapsulateFrame(frame []byte) ([]byte, error) {
	timeInByte := make([]byte, 8)
	binary.BigEndian.PutUint64(timeInByte, uint64(time.Now().UnixNano()))

	return append(timeInByte, frame...), nil
}

// DecapsulateFrame decapsulates a frame
func DecapsulateFrame(frame []byte) ([]byte, error) {
	if len(frame) < (TIMESTAMP_SIZE + 1) {
		return nil, errors.New("invalid encapsulated frame size")
	}

	timeInByte := frame[0:TIMESTAMP_SIZE]
	timeInUnixNano := int64(binary.BigEndian.Uint64(timeInByte))
	if (time.Now().UnixNano() - timeInUnixNano) > TIMESTAMP_TIMEOUT.Nanoseconds() {
		return nil, errors.New("timestamp out of acceptable range")
	}

	return frame[TIMESTAMP_SIZE:], nil
}

// TAPDevice is a TAP device
type TAPDevice struct {
	fd      int
	devName string
}

// Open opens the TAP device
func (t *TAPDevice) Open(mtu uint) error {
	fd, err := syscall.Open("/dev/net/tun", syscall.O_RDWR, syscall.S_IRUSR|syscall.S_IWUSR|syscall.S_IRGRP|syscall.S_IROTH)
	if err != nil {
		return fmt.Errorf("could not open /dev/net/tun: %v", err)
	}
	t.fd = fd

	ifreqFlags := uint16(syscall.IFF_TAP | syscall.IFF_NO_PI)
	ifreq := make([]byte, 32)
	ifreq[16] = byte(ifreqFlags)
	ifreq[17] = byte(ifreqFlags >> 8)
	currentFlags, _, err := syscall.Syscall(syscall.SYS_IOCTL, uintptr(t.fd), syscall.TUNSETIFF, uintptr(unsafe.Pointer(&ifreq[0])))
	if currentFlags != 0 {
		t.Close()

		return fmt.Errorf("could not set tun/tap type: %v", err)
	}

	t.devName = string(ifreq[0:16])

	rawSocketFd, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW, syscall.ETH_P_ALL)
	if err != nil {
		t.Close()

		return fmt.Errorf("could not create packet socket: %v", err)
	}
	defer syscall.Close(rawSocketFd)

	err = syscall.BindToDevice(rawSocketFd, t.devName)
	if err != nil {
		t.Close()

		return fmt.Errorf("could not bind packet socket to TAP device: %v", err)
	}

	ifreqMTU := mtu
	ifreq[16] = byte(ifreqMTU)
	ifreq[17] = byte(ifreqMTU >> 8)
	ifreq[18] = byte(ifreqMTU >> 16)
	ifreq[19] = byte(ifreqMTU >> 24)
	currentFlags, _, err = syscall.Syscall(syscall.SYS_IOCTL, uintptr(rawSocketFd), syscall.SIOCSIFMTU, uintptr(unsafe.Pointer(&ifreq[0])))
	if currentFlags != 0 {
		t.Close()

		return fmt.Errorf("could not set MTU on TAP device: %v", err)
	}

	currentFlags, _, err = syscall.Syscall(syscall.SYS_IOCTL, uintptr(rawSocketFd), syscall.SIOCGIFFLAGS, uintptr(unsafe.Pointer(&ifreq[0])))
	if currentFlags != 0 {
		t.Close()

		return fmt.Errorf("could not get TAP device flags: %v", err)
	}

	ifreqFlags = uint16(ifreq[16]) | (uint16(ifreq[17]) << 8)
	ifreqFlags |= syscall.IFF_UP | syscall.IFF_RUNNING
	ifreq[16] = byte(ifreqFlags)
	ifreq[17] = byte(ifreqFlags >> 8)
	currentFlags, _, err = syscall.Syscall(syscall.SYS_IOCTL, uintptr(rawSocketFd), syscall.SIOCSIFFLAGS, uintptr(unsafe.Pointer(&ifreq[0])))
	if currentFlags != 0 {
		t.Close()

		return fmt.Errorf("could not bring up TAP device: %v", err)
	}

	return nil
}

// Close closes the TAP device
func (t *TAPDevice) Close() {
	syscall.Close(t.fd)
}

// Read reads from the TAP device
func (t *TAPDevice) Read(b []byte) (int, error) {
	return syscall.Read(t.fd, b)
}

// Write writes to the TAP device
func (t *TAPDevice) Write(b []byte) (n int, err error) {
	return syscall.Write(t.fd, b)
}

// ForwardTCPtoTAP forwards TCP packets to a TAP device
func ForwardTCPtoTAP(tcpConn *net.TCPListener, tapDevice *TAPDevice, remoteAddr *net.TCPAddr) {
	log.Printf("forwarding TCP to TAP with remote %v:%v\n", remoteAddr.IP, remoteAddr.Port)

	for {
		packet := make([]byte, TCP_MTU)
		var frame []byte

		conn, err := tcpConn.AcceptTCP()
		if err != nil {
			log.Fatal(err)
		}

		n, err := conn.Read(packet)
		if err != nil {
			log.Fatalf("could not read from TCP socket: %v\n", err)
		}

		if err := conn.Close(); err != nil {
			log.Fatal(err)
		}

		frame, invalid := DecapsulateFrame(packet[0:n])
		if invalid != nil {
			continue
		}

		_, err = tapDevice.Write(frame)
		if err != nil {
			log.Fatalf("could not write to TAP device: %v\n", err)
		}
	}
}

// ForwardTAPtoTCP forwards frames from a TAP device to a TCP connection
func ForwardTAPtoTCP(tapDevice *TAPDevice, remoteAddr *net.TCPAddr) {
	log.Printf("forwarding TAP to TCP with remote %v:%v\n", remoteAddr.IP, remoteAddr.Port)

	for {
		frame := make([]byte, TAP_MTU+14)
		var encFrame []byte

		n, err := tapDevice.Read(frame)
		if err != nil {
			log.Fatalf("could not read from TAP device: %v\n", err)
		}

		encFrame, invalid := EncapsulateFrame(frame[0:n])
		if invalid != nil {
			continue
		}

		conn, err := net.Dial("tcp", remoteAddr.String())
		if err != nil {
			log.Printf("could not dial %v, retrying", remoteAddr)

			ForwardTAPtoTCP(tapDevice, remoteAddr)

			continue
		}

		_, err = conn.Write(encFrame)
		if err != nil {
			log.Fatalf("could not write to TCP socket: %v\n", err)
		}

		if err := conn.Close(); err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	localAddrFlag := flag.String("localAddr", "0.0.0.0:1234", "Local address")
	remoteAddrFlag := flag.String("remoteAddr", "0.0.0.0:12345", "Remote address")

	flag.Parse()

	localAddr, err := net.ResolveTCPAddr("tcp", *localAddrFlag)
	if err != nil {
		log.Fatalf("could not resolve local address: %v\n", err)
	}

	remoteAddr, err := net.ResolveTCPAddr("tcp", *remoteAddrFlag)
	if err != nil {
		log.Fatalf("could not resolve remote address: %v\n", err)
	}

	tcpListener, err := net.ListenTCP("tcp", localAddr)
	if err != nil {
		log.Fatalf("error creating a TCP socket: %v\n", err)
	}

	tapDev := new(TAPDevice)
	err = tapDev.Open(TAP_MTU)
	if err != nil {
		log.Fatalf("could not open a TAP device: %v\n", err)
	}

	log.Printf("started tunnel with TAP device %v", tapDev.devName)

	go ForwardTCPtoTAP(tcpListener, tapDev, remoteAddr)
	ForwardTAPtoTCP(tapDev, remoteAddr)
}

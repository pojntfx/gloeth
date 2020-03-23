package tap

import (
	"fmt"
	"syscall"
	"unsafe"
)

// Device is a TAP device
type Device struct {
	fd      int
	devName string
}

// Create a new device
func NewDevice() *Device {
	return &Device{}
}

// GetName returns the TAP device name
func (t *Device) GetName() string {
	return t.devName
}

// Open opens the TAP device
func (t *Device) Open(mtu uint) error {
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
func (t *Device) Close() {
	syscall.Close(t.fd)
}

// Read reads from the TAP device
func (t *Device) Read(b []byte) (int, error) {
	return syscall.Read(t.fd, b)
}

// Write writes to the TAP device
func (t *Device) Write(b []byte) (n int, err error) {
	return syscall.Write(t.fd, b)
}

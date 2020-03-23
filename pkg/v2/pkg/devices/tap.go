package devices

import (
	"syscall"
	"unsafe"

	"github.com/pojntfx/gloeth/pkg/v2/pkg/constants"
)

// TAPDevice is a TAP device
type TAPDevice struct {
	mtu        uint
	framesChan chan []byte
	fd         int
}

// NewTAPDevice creates a new TAP device
func NewTAPDevice(mtu uint, framesChan chan []byte) *TAPDevice {
	return &TAPDevice{
		mtu:        mtu,
		framesChan: framesChan,
	}
}

// Open opens the TAP device
func (t *TAPDevice) Open() error {
	fd, err := syscall.Open("/dev/net/tun", syscall.O_RDWR, syscall.S_IRUSR|syscall.S_IWUSR|syscall.S_IRGRP|syscall.S_IROTH)
	if err != nil {
		return err
	}
	t.fd = fd

	ifreqFlags := uint16(syscall.IFF_TAP | syscall.IFF_NO_PI)
	ifreq := make([]byte, 32)
	ifreq[16] = byte(ifreqFlags)
	ifreq[17] = byte(ifreqFlags >> 8)
	currentFlags, _, err := syscall.Syscall(syscall.SYS_IOCTL, uintptr(t.fd), syscall.TUNSETIFF, uintptr(unsafe.Pointer(&ifreq[0])))
	if currentFlags != 0 {
		t.Close()

		return err
	}

	devName := string(ifreq[0:16])

	rawSocketFd, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW, syscall.ETH_P_ALL)
	if err != nil {
		t.Close()

		return err
	}
	defer syscall.Close(rawSocketFd)

	err = syscall.BindToDevice(rawSocketFd, devName)
	if err != nil {
		t.Close()

		return err
	}

	ifreqMTU := t.mtu
	ifreq[16] = byte(ifreqMTU)
	ifreq[17] = byte(ifreqMTU >> 8)
	ifreq[18] = byte(ifreqMTU >> 16)
	ifreq[19] = byte(ifreqMTU >> 24)
	currentFlags, _, err = syscall.Syscall(syscall.SYS_IOCTL, uintptr(rawSocketFd), syscall.SIOCSIFMTU, uintptr(unsafe.Pointer(&ifreq[0])))
	if currentFlags != 0 {
		t.Close()

		return err
	}

	currentFlags, _, err = syscall.Syscall(syscall.SYS_IOCTL, uintptr(rawSocketFd), syscall.SIOCGIFFLAGS, uintptr(unsafe.Pointer(&ifreq[0])))
	if currentFlags != 0 {
		t.Close()

		return err
	}

	ifreqFlags = uint16(ifreq[16]) | (uint16(ifreq[17]) << 8)
	ifreqFlags |= syscall.IFF_UP | syscall.IFF_RUNNING
	ifreq[16] = byte(ifreqFlags)
	ifreq[17] = byte(ifreqFlags >> 8)
	currentFlags, _, err = syscall.Syscall(syscall.SYS_IOCTL, uintptr(rawSocketFd), syscall.SIOCSIFFLAGS, uintptr(unsafe.Pointer(&ifreq[0])))
	if currentFlags != 0 {
		t.Close()

		return err
	}

	return nil
}

// Close closes the TAP device
func (t *TAPDevice) Close() error {
	return syscall.Close(t.fd)
}

// Read reads from the TAP device
func (t *TAPDevice) Read() error {
	for {
		frame := make([]byte, constants.FRAME_SIZE)

		_, err := syscall.Read(t.fd, frame)
		if err != nil {
			return err
		}

		t.framesChan <- frame
	}
}

// Write writes to the TAP device
func (t *TAPDevice) Write(frame []byte) (int, error) {
	return syscall.Write(t.fd, frame)
}

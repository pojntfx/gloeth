package connections

import "net"

// GetConn returns a TCP connection to a remote TCP server
func GetConn(raddr *net.TCPAddr) (*net.TCPConn, error) {
	conn, err := net.DialTCP("tcp", nil, raddr)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

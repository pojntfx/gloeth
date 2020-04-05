package switchers

import (
	"fmt"
	"net"
	"reflect"
	"testing"
	"time"

	gm "github.com/cseeger-epages/mac-gen-go"
	cmap "github.com/orcaman/concurrent-map"
	"github.com/pojntfx/gloeth/pkg/wrappers"
)

func getMACAddress() (net.HardwareAddr, error) {
	prefix := gm.GenerateRandomLocalMacPrefix(false)
	suffix, err := gm.CalculateNICSufix(net.ParseIP("10.0.0.1"))
	if err != nil {
		return nil, err
	}

	rawDest := fmt.Sprintf("%v:%v", prefix, suffix)

	return net.ParseMAC(rawDest)
}

func getBroadcastMACAddress() (net.HardwareAddr, error) {
	return net.ParseMAC("ff:ff:ff:ff:ff:ff")
}

func getListener() (*net.TCPAddr, *net.TCPListener, error) {
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		return nil, nil, err
	}

	return l.Addr().(*net.TCPAddr), l.(*net.TCPListener), err
}

func getConn(raddr *net.TCPAddr) (*net.TCPConn, error) {
	conn, err := net.DialTCP("tcp", nil, raddr)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func getFrame() [wrappers.WrappedFrameSize]byte {
	return [wrappers.WrappedFrameSize]byte{1}
}

func writeFrame(conn *net.TCPConn, frame [wrappers.WrappedFrameSize]byte) error {
	_, err := conn.Write(frame[:])

	return err
}

func readFrame(conn *net.TCPConn) ([wrappers.WrappedFrameSize]byte, error) {
	frame := [wrappers.WrappedFrameSize]byte{}

	_, err := conn.Read(frame[:])
	if err != nil {
		return [wrappers.WrappedFrameSize]byte{}, err
	}

	return frame, nil
}

func getConns(conns map[string]*net.TCPConn) cmap.ConcurrentMap {
	iconns := cmap.New()

	for mac, conn := range conns {
		iconns.Set(mac, conn)
	}

	return iconns
}

func TestNewTCP(t *testing.T) {
	expectedReadChan := make(chan [wrappers.WrappedFrameSize]byte)
	expectedConnChan := make(chan *net.TCPConn)
	expectedLaddr, _, err := getListener()
	if err != nil {
		t.Error(err)
	}
	expectedConns := map[string]*net.TCPConn{}
	iexpectedConns := getConns(expectedConns)

	type args struct {
		readChan   chan [wrappers.WrappedFrameSize]byte
		connChan   chan *net.TCPConn
		listenAddr *net.TCPAddr
		conns      map[string]*net.TCPConn
	}
	tests := []struct {
		name string
		args args
		want *TCP
	}{
		{
			"New",
			args{
				expectedReadChan,
				expectedConnChan,
				expectedLaddr,
				expectedConns,
			},
			&TCP{
				expectedReadChan,
				expectedConnChan,
				expectedLaddr,
				nil,
				iexpectedConns,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTCP(tt.args.readChan, tt.args.connChan, tt.args.listenAddr, tt.args.conns); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTCP() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTCP_Open(t *testing.T) {
	readChan := make(chan [wrappers.WrappedFrameSize]byte)
	connChan := make(chan *net.TCPConn)
	laddr, listener, err := getListener()
	if err != nil {
		t.Error(err)
	}
	if err := listener.Close(); err != nil {
		t.Error(err)
	}

	type fields struct {
		readChan   chan [wrappers.WrappedFrameSize]byte
		connChan   chan *net.TCPConn
		listenAddr *net.TCPAddr
		listener   *net.TCPListener
		conns      cmap.ConcurrentMap
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			"Open",
			fields{
				readChan,
				connChan,
				laddr,
				nil,
				nil,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &TCP{
				readChan: tt.fields.readChan,
				connChan: tt.fields.connChan,
				laddr:    tt.fields.listenAddr,
				listener: tt.fields.listener,
				conns:    tt.fields.conns,
			}
			if err := s.Open(); (err != nil) != tt.wantErr {
				t.Errorf("TCP.Open() error = %v, wantErr %v", err, tt.wantErr)
			}

			timeoutChan := make(chan bool)
			timeout := time.Millisecond * 10
			go func() {
				_, err := getConn(laddr)
				if err != nil {
					t.Error(err)
				}

				timeoutChan <- false
			}()

			go func() {
				time.Sleep(timeout)

				timeoutChan <- true
			}()

			if <-timeoutChan {
				t.Errorf("TCP.Open() did not connect to TCP client within %v", timeout)
			}
		})
	}
}

func TestTCP_Close(t *testing.T) {
	readChan := make(chan [wrappers.WrappedFrameSize]byte)
	connChan := make(chan *net.TCPConn)
	laddr, listener, err := getListener()
	if err != nil {
		t.Error(err)
	}

	type fields struct {
		readChan chan [wrappers.WrappedFrameSize]byte
		connChan chan *net.TCPConn
		laddr    *net.TCPAddr
		listener *net.TCPListener
		conns    cmap.ConcurrentMap
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			"Close",
			fields{
				readChan,
				connChan,
				laddr,
				listener,
				nil,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &TCP{
				readChan: tt.fields.readChan,
				connChan: tt.fields.connChan,
				laddr:    tt.fields.laddr,
				listener: tt.fields.listener,
				conns:    tt.fields.conns,
			}
			if err := s.Close(); (err != nil) != tt.wantErr {
				t.Errorf("TCP.Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTCP_Read(t *testing.T) {
	readChan := make(chan [wrappers.WrappedFrameSize]byte)
	connChan := make(chan *net.TCPConn)
	laddr, listener, err := getListener()
	if err != nil {
		t.Error(err)
	}
	conn, err := getConn(laddr)
	if err != nil {
		t.Error(err)
	}
	expectedFrame := getFrame()

	type fields struct {
		readChan chan [wrappers.WrappedFrameSize]byte
		connChan chan *net.TCPConn
		laddr    *net.TCPAddr
		listener *net.TCPListener
		conns    cmap.ConcurrentMap
	}
	tests := []struct {
		name               string
		fields             fields
		frameToWrite, want [wrappers.WrappedFrameSize]byte
		framesToTransceive uint
		wantErr            bool
	}{
		{
			"Read",
			fields{
				readChan,
				connChan,
				laddr,
				listener,
				nil,
			},
			expectedFrame,
			expectedFrame,
			5,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &TCP{
				readChan: tt.fields.readChan,
				connChan: tt.fields.connChan,
				laddr:    tt.fields.laddr,
				listener: tt.fields.listener,
				conns:    tt.fields.conns,
			}

			go func() {
				if err := s.Read(); (err != nil) != tt.wantErr {
					t.Errorf("TCP.Read() error = %v, wantErr %v", err, tt.wantErr)
				}
			}()

			go func() {
				time.Sleep(time.Millisecond * 5)

				for i := 0; i < int(tt.framesToTransceive); i++ {
					if err := writeFrame(conn, expectedFrame); err != nil {
						t.Error(err)
					}
				}
			}()

			conn := <-connChan

			for i := 0; i < int(tt.framesToTransceive); i++ {
				got, err := readFrame(conn)
				if err != nil {
					t.Error(err)
				}

				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("TCP.Read() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestTCP_HandleFrame(t *testing.T) {
	readChan := make(chan [wrappers.WrappedFrameSize]byte)
	connChan := make(chan *net.TCPConn)
	expectedFrame := getFrame()

	type fields struct {
		readChan chan [wrappers.WrappedFrameSize]byte
		connChan chan *net.TCPConn
		laddr    *net.TCPAddr
		listener *net.TCPListener
		conns    cmap.ConcurrentMap
	}
	type args struct {
		frame [wrappers.WrappedFrameSize]byte
	}
	tests := []struct {
		name               string
		fields             fields
		args               args
		framesToTransceive uint
		want               [wrappers.WrappedFrameSize]byte
	}{
		{
			"HandleFrame",
			fields{
				readChan,
				connChan,
				nil,
				nil,
				nil,
			},
			args{
				expectedFrame,
			},
			5,
			expectedFrame,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &TCP{
				readChan: tt.fields.readChan,
				connChan: tt.fields.connChan,
				laddr:    tt.fields.laddr,
				listener: tt.fields.listener,
				conns:    tt.fields.conns,
			}

			go func() {
				for i := 0; i < int(tt.framesToTransceive); i++ {
					s.HandleFrame(tt.args.frame)
				}
			}()

			for i := 0; i < int(tt.framesToTransceive); i++ {
				got := <-readChan

				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("TCP.HandleFrame() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestTCP_Register(t *testing.T) {
	readChan := make(chan [wrappers.WrappedFrameSize]byte)
	connChan := make(chan *net.TCPConn)
	laddr, _, err := getListener()
	if err != nil {
		t.Error(err)
	}
	conn, err := getConn(laddr)
	if err != nil {
		t.Error(err)
	}
	mac, err := getMACAddress()
	if err != nil {
		t.Error(err)
	}

	type fields struct {
		readChan chan [wrappers.WrappedFrameSize]byte
		connChan chan *net.TCPConn
		laddr    *net.TCPAddr
		listener *net.TCPListener
		conns    cmap.ConcurrentMap
	}
	type args struct {
		mac  *net.HardwareAddr
		conn *net.TCPConn
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *net.TCPConn
	}{
		{
			"Register",
			fields{
				readChan,
				connChan,
				nil,
				nil,
				getConns(map[string]*net.TCPConn{}),
			},
			args{
				&mac,
				conn,
			},
			conn,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &TCP{
				readChan: tt.fields.readChan,
				connChan: tt.fields.connChan,
				laddr:    tt.fields.laddr,
				listener: tt.fields.listener,
				conns:    tt.fields.conns,
			}

			s.Register(tt.args.mac, tt.args.conn)

			got, ok := s.conns.Get(tt.args.mac.String())
			if !ok {
				t.Errorf("TCP.Register() ok = %v, want %v", ok, true)
			}

			if got.(*net.TCPConn) != tt.want {
				t.Errorf("TCP.Register() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTCP_GetConnectionForMAC(t *testing.T) {
	readChan := make(chan [wrappers.WrappedFrameSize]byte)
	connChan := make(chan *net.TCPConn)
	laddr, _, err := getListener()
	if err != nil {
		t.Error(err)
	}
	conn1, err := getConn(laddr)
	if err != nil {
		t.Error(err)
	}
	conn2, err := getConn(laddr)
	if err != nil {
		t.Error(err)
	}
	mac1, err := getMACAddress()
	if err != nil {
		t.Error(err)
	}
	mac2, err := getMACAddress()
	if err != nil {
		t.Error(err)
	}

	type fields struct {
		readChan chan [wrappers.WrappedFrameSize]byte
		connChan chan *net.TCPConn
		laddr    *net.TCPAddr
		listener *net.TCPListener
		conns    cmap.ConcurrentMap
	}
	type args struct {
		destMAC *net.HardwareAddr
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *net.TCPConn
		wantErr bool
	}{
		{
			"GetConnectionForMAC",
			fields{
				readChan,
				connChan,
				nil,
				nil,
				getConns(map[string]*net.TCPConn{
					mac1.String(): conn1,
					mac2.String(): conn2,
				}),
			},
			args{
				&mac1,
			},
			conn1,
			false,
		},
		{
			"GetConnectionForMAC (unknown connection)",
			fields{
				readChan,
				connChan,
				nil,
				nil,
				getConns(map[string]*net.TCPConn{
					mac1.String(): conn1,
				}),
			},
			args{
				&mac2,
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &TCP{
				readChan: tt.fields.readChan,
				connChan: tt.fields.connChan,
				laddr:    tt.fields.laddr,
				listener: tt.fields.listener,
				conns:    tt.fields.conns,
			}
			got, err := s.GetConnectionsForMAC(tt.args.destMAC)
			if (err != nil) != tt.wantErr {
				t.Errorf("TCP.GetConnectionsForMAC() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TCP.GetConnectionsForMAC() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTCP_Write(t *testing.T) {
	readChan := make(chan [wrappers.WrappedFrameSize]byte)
	connChan := make(chan *net.TCPConn)
	laddr, listener, err := getListener()
	if err != nil {
		t.Error(err)
	}
	conn, err := getConn(laddr)
	if err != nil {
		t.Error(err)
	}
	expectedFrame := getFrame()

	type fields struct {
		readChan chan [wrappers.WrappedFrameSize]byte
		connChan chan *net.TCPConn
		laddr    *net.TCPAddr
		listener *net.TCPListener
		conns    cmap.ConcurrentMap
	}
	type args struct {
		conn  *net.TCPConn
		frame [wrappers.WrappedFrameSize]byte
	}
	tests := []struct {
		name               string
		fields             fields
		args               args
		framesToTransceive uint
		want               [wrappers.WrappedFrameSize]byte
		wantErr            bool
	}{
		{
			"Write",
			fields{
				readChan,
				connChan,
				nil,
				nil,
				nil,
			},
			args{
				conn,
				expectedFrame,
			},
			5,
			expectedFrame,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &TCP{
				readChan: tt.fields.readChan,
				connChan: tt.fields.connChan,
				laddr:    tt.fields.laddr,
				listener: tt.fields.listener,
				conns:    tt.fields.conns,
			}

			doneChan := make(chan bool)

			go func() {
				time.Sleep(time.Millisecond * 5)

				conn, err := listener.AcceptTCP()
				if err != nil {
					t.Error(err)
				}

				for i := 0; i < int(tt.framesToTransceive); i++ {
					got, err := readFrame(conn)
					if err != nil {
						t.Error(err)
					}

					if !reflect.DeepEqual(got, tt.want) {
						t.Errorf("read(TCP.Write()) = %v, want %v", got, tt.want)
					}
				}

				doneChan <- true
			}()

			for i := 0; i < int(tt.framesToTransceive); i++ {
				if err := s.Write(tt.args.conn, tt.args.frame); (err != nil) != tt.wantErr {
					t.Errorf("TCP.Write() error = %v, wantErr %v", err, tt.wantErr)
				}
			}

			<-doneChan
		})
	}
}

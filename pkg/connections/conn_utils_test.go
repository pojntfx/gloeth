package connections

import (
	"net"
	"testing"
)

func TestGetConn(t *testing.T) {
	raddr, listener, err := getListener()
	if err != nil {
		t.Error(err)
	}

	type args struct {
		raddr *net.TCPAddr
	}
	tests := []struct {
		name                string
		args                args
		connectionsToCreate uint
		wantErr             bool
	}{
		{
			"GetConn",
			args{
				raddr,
			},
			5,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			go func() {
				for i := 0; i < int(tt.connectionsToCreate); i++ {
					got, err := GetConn(tt.args.raddr)
					if (err != nil) != tt.wantErr {
						t.Errorf("GetConn() error = %v, wantErr %v", err, tt.wantErr)
						return
					}
					if got == nil {
						t.Errorf("GetConn() = %v, want %v", got, nil)
					}
					if err := got.Close(); err != nil {
						t.Error(err)
					}
				}
			}()

			connectionsOpened := 0
			for i := 0; i < int(tt.connectionsToCreate); i++ {
				if _, err := listener.AcceptTCP(); err != nil {
					t.Error(err)
				}

				connectionsOpened = connectionsOpened + 1
			}

			if connectionsOpened != int(tt.connectionsToCreate) {
				t.Errorf("connectionsOpened = %v, want %v", connectionsOpened, tt.connectionsToCreate)
			}
		})
	}
}
